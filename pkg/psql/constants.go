package psql

import "time"

const (
	// configuration
	DefaultDBHost         = "localhost"
	DefaultDBPort         = 5432
	DefaultDBUser         = "postgres"
	DefaultDBPass         = "password"
	DefaultDBName         = "postgres"
	DefaultDBQueryTimeout = 10 * time.Second

	// queries
	QueryAllTables string = `
		SELECT table_name 
		FROM information_schema.tables 
		WHERE table_schema = 'public'
		ORDER BY table_name;
	`
)
