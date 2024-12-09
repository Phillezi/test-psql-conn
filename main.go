package main

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/Phillezi/test-psql-conn/config"
	"github.com/Phillezi/test-psql-conn/util"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

func main() {
	cfg := config.Load()
	defer func() {
		if cfg.SleepWhenDone {
			sleepInfinityly()
		}
	}()

	if cfg.DBUser == "" || cfg.DBPass == "" || cfg.DBName == "" {
		logrus.Errorln("DB_USER, DB_PASS and DB all have to be set")

		util.LogEnvTable()
		return
	}

	cfg.LogConfig()

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPass, cfg.DBName)

	maxOpenConns := util.GetEnvAsInt("DB_MAX_OPEN_CONNS", 10)
	maxIdleConns := util.GetEnvAsInt("DB_MAX_IDLE_CONNS", 5)
	connMaxLifetime := util.GetEnvAsDuration("DB_CONN_MAX_LIFETIME", 30*time.Minute)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		logrus.Errorln("Did not connect")
		logrus.Errorln("failed to open database: ", err)
		if cfg.SleepWhenDone {
			sleepInfinityly()
		}
		return
	}
	defer db.Close()

	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxLifetime(connMaxLifetime)

	err = db.Ping()
	if err != nil {
		logrus.Errorln("Did not connect")
	} else {
		logrus.Infoln("Connected")
	}
}

func sleepInfinityly() {
	logrus.Infoln("Sleeping forever since program exited")
	for {
		time.Sleep(100000 * time.Hour)
	}
}
