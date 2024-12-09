package client

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/Phillezi/test-psql-conn/config"
	"github.com/Phillezi/test-psql-conn/internal/models"
	"github.com/Phillezi/test-psql-conn/util"
	"github.com/sirupsen/logrus"
)

type Client struct {
	dsn        string
	connStatus chan bool
	tablesChan chan []models.Table
}

func New(cfg *config.Config, connStatus chan bool, tablesChan chan []models.Table) *Client {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPass, cfg.DBName)
	return &Client{
		dsn:        dsn,
		connStatus: connStatus,
		tablesChan: tablesChan,
	}
}

func (c *Client) Start() {
	maxOpenConns := util.GetEnvAsInt("DB_MAX_OPEN_CONNS", 10)
	maxIdleConns := util.GetEnvAsInt("DB_MAX_IDLE_CONNS", 5)
	connMaxLifetime := util.GetEnvAsDuration("DB_CONN_MAX_LIFETIME", 30*time.Minute)

	db, err := sql.Open("postgres", c.dsn)
	if err != nil {
		logrus.Errorln("Did not connect")
		logrus.Errorln("failed to open database: ", err)
		return
	}
	defer db.Close()

	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxLifetime(connMaxLifetime)

	err = db.Ping()
	if err != nil {
		c.connStatus <- false
		logrus.Errorln("Did not connect")
	} else {
		c.connStatus <- true
		logrus.Infoln("Connected")
		c.fetchTablesAndCounts(db)
	}
}

func (c *Client) fetchTablesAndCounts(db *sql.DB) {
	query := `
		SELECT table_name 
		FROM information_schema.tables 
		WHERE table_schema = 'public'
		ORDER BY table_name;
	`
	rows, err := db.Query(query)
	if err != nil {
		logrus.Errorln("Failed to fetch tables:", err)
		c.tablesChan <- nil
		return
	}
	defer rows.Close()

	var tableQueries []string
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			logrus.Errorln("Error scanning table name:", err)
			continue
		}
		tableQueries = append(tableQueries, fmt.Sprintf("SELECT '%s' AS table_name, COUNT(*) AS row_count FROM %s", tableName, tableName))
	}

	if err := rows.Err(); err != nil {
		logrus.Errorln("Error iterating through table names:", err)

		return
	}

	fullQuery := strings.Join(tableQueries, " UNION ALL ")

	countRows, err := db.Query(fullQuery)
	if err != nil {
		logrus.Errorln("Failed to execute count query:", err)
		return
	}
	defer countRows.Close()

	var tables []models.Table
	for countRows.Next() {
		var table_name string
		var row_count int
		if err := countRows.Scan(&table_name, &row_count); err != nil {
			logrus.Errorln("Error scanning row count:", err)
			continue
		}
		tables = append(tables, models.Table{Name: table_name, Count: row_count})
		fmt.Println("name: ", table_name, " count: ", row_count)
	}

	select {
	case c.tablesChan <- tables:
		logrus.Infoln("Tables sent successfully")
	default:
		logrus.Warnln("Channel is full, tables not sent")
	}
}
