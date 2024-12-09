package main

import (
	"context"
	"time"

	"github.com/Phillezi/test-psql-conn/config"
	"github.com/Phillezi/test-psql-conn/internal/client"
	"github.com/Phillezi/test-psql-conn/internal/models"
	"github.com/Phillezi/test-psql-conn/internal/server"
	"github.com/Phillezi/test-psql-conn/util"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

func main() {
	cfg := config.Load()
	defer func() {
		if cfg.SleepWhenDone {
			sleepInf()
		}
	}()

	if cfg.DBUser == "" || cfg.DBPass == "" || cfg.DBName == "" {
		logrus.Errorln("DB_USER, DB_PASS and DB all have to be set")

		util.LogEnvTable()
		return
	}

	cfg.LogConfig()

	connStatus := make(chan bool)
	tablesChan := make(chan []models.Table)

	srv := server.New(context.Background(), 8080, connStatus, tablesChan)
	db := client.New(cfg, connStatus, tablesChan)

	if cfg.ServeHTTP {
		go srv.Start()
	}
	db.Start()
}

func sleepInf() {
	logrus.Infoln("Sleeping forever since program exited")
	for {
		time.Sleep(100000 * time.Hour)
	}
}
