package psql

import (
	"database/sql"
	"fmt"

	"github.com/Phillezi/test-psql-conn/util"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

func (c *Client) Connect() (err error) {
	if c.db != nil {
		return fmt.Errorf("aldready connected")
	}

	c.dbMu.Lock()
	c.db, err = sql.Open("postgres", c.buildDSN())
	c.dbMu.Unlock()
	if err != nil {
		c.sendStatusAsync(false)
		return err
	}

	c.dbMu.RLock()
	if err = c.db.Ping(); err != nil {
		c.dbMu.RUnlock()
		c.sendStatusAsync(false)
		func() {
			c.dbMu.Lock()
			defer c.dbMu.Unlock()
			c.db.Close()
			c.db = nil
		}()
		zap.L().Error("database ping failed", zap.Error(err))
		return err
	}
	c.dbMu.RUnlock()

	c.sendStatusAsync(true)
	zap.L().Info("succesfully connected to db")

	go func() {
		<-c.ctx.Done()
		c.dbMu.RLock()
		defer c.dbMu.RUnlock()
		util.IfNotNilDo(c.db, func() {
			if err := c.db.Close(); err != nil {
				zap.L().Error("failed to close db", zap.Error(err))
			}
			c.db = nil
		})
	}()

	return nil
}

func (c *Client) buildDSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		util.Or(c.dbHost, DefaultDBHost),
		util.Or(c.dbPort, DefaultDBPort),
		util.Or(c.dbUser, DefaultDBUser),
		util.Or(c.dbPass, DefaultDBPass),
		util.Or(c.dbName, DefaultDBName),
		func() string {
			if c.dbSSLEnabled {
				return "enable"
			}
			return "disable"
		}(),
	)
}
