package psql

import (
	"time"

	"github.com/Phillezi/test-psql-conn/util"
	"go.uber.org/zap"
)

func (c *Client) Probe(freq ...time.Duration) error {
	zap.L().Info("started db prober")
	defer zap.L().Info("db prober exited")

	if c.db == nil {
		if err := c.Connect(); err != nil {
			zap.L().Error("failed to connect to db", zap.Error(err))
		} else {
			c.Query()
		}
	}

	var frequency time.Duration
	if len(freq) > 0 {
		frequency = freq[0]
	}

	ticker := time.NewTicker(util.Or(frequency, 10*time.Second))
	defer ticker.Stop()

	for {
		select {
		case <-c.ctx.Done():
			zap.L().Info("context cancelled, exitting prober")
			return nil
		case <-ticker.C:
			if c.db == nil {
				zap.L().Debug("db is nil, trying to reconnect")
				if err := c.Connect(); err != nil {
					zap.L().Error("failed to connect to db", zap.Error(err))
				} else {
					c.Query()
				}
				continue
			}
			c.dbMu.RLock()
			if err := c.db.Ping(); err != nil {
				c.dbMu.RUnlock()
				zap.L().Info("database connection lost", zap.Error(err))
				c.sendStatusAsync(false)
				func() {
					c.dbMu.Lock()
					defer c.dbMu.Unlock()
					c.db = nil
				}()
				continue
			}
			c.dbMu.RUnlock()
			c.sendStatusAsync(true)
			zap.L().Debug("succesfully pinged db")
		}
	}
}
