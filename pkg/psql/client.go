package psql

import (
	"context"
	"database/sql"
	"sync"

	"github.com/Phillezi/test-psql-conn/pkg/model/table"
	"github.com/Phillezi/test-psql-conn/util"
)

type Client struct {
	ctx context.Context

	dbHost       string
	dbPort       int
	dbUser       string
	dbPass       string
	dbName       string
	dbSSLEnabled bool

	dbMu sync.RWMutex
	db   *sql.DB

	connStatus chan bool
	tables     chan []table.Table
}

func New() *Client {
	return &Client{
		ctx: context.Background(),

		connStatus: make(chan bool),
		tables:     make(chan []table.Table),
	}
}

func (c *Client) WithContext(ctx context.Context) *Client {
	c.ctx = ctx
	return c
}

type ClientOpts struct {
	DBHost       *string
	DBPort       *int
	DBUser       *string
	DBPass       *string
	DBName       *string
	DBSSLEnabled *bool
}

func (c *Client) WithOptions(opts ClientOpts) *Client {
	util.IfNotNilDo(opts.DBHost, func() { c.dbHost = *opts.DBHost })
	util.IfNotNilDo(opts.DBPort, func() { c.dbPort = *opts.DBPort })
	util.IfNotNilDo(opts.DBUser, func() { c.dbUser = *opts.DBUser })
	util.IfNotNilDo(opts.DBPass, func() { c.dbPass = *opts.DBPass })
	util.IfNotNilDo(opts.DBName, func() { c.dbName = *opts.DBName })
	util.IfNotNilDo(opts.DBSSLEnabled, func() { c.dbSSLEnabled = *opts.DBSSLEnabled })
	return c
}
