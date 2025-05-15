package psql

import (
	"context"
	"database/sql"
	"fmt"
	"slices"
	"strings"

	"github.com/Phillezi/test-psql-conn/pkg/model/table"
	"github.com/Phillezi/test-psql-conn/util"
	"go.uber.org/zap"
)

func (c *Client) Query() (tables []table.Table, err error) {
	c.dbMu.RLock()
	if c.db == nil {
		c.dbMu.RUnlock()
		return nil, fmt.Errorf("db is nil")
	}

	queryTablesCtx, cancelQueryTablesCtx := context.WithTimeout(c.ctx, util.Or(DefaultDBQueryTimeout))
	defer cancelQueryTablesCtx()
	rows, err := c.db.QueryContext(queryTablesCtx, QueryAllTables)
	c.dbMu.RUnlock()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	countQuery, err := constructCountQuery(rows)
	if err != nil {
		return nil, err
	}

	queryCountCtx, cancelQueryCountCtx := context.WithTimeout(c.ctx, util.Or(DefaultDBQueryTimeout))
	defer cancelQueryCountCtx()
	c.dbMu.RLock()
	countRows, err := c.db.QueryContext(queryCountCtx, countQuery)
	c.dbMu.RUnlock()
	if err != nil {
		return nil, err
	}
	defer countRows.Close()

	for countRows.Next() {
		var table_name string
		var row_count int
		if err := countRows.Scan(&table_name, &row_count); err != nil {
			zap.L().Error("error scanning row count", zap.Error(err))
			continue
		}
		tables = append(tables, table.Table{Name: table_name, Count: row_count})
	}

	c.sendTablesAsync(slices.Clone(tables))

	return tables, nil
}

func constructCountQuery(rows *sql.Rows) (string, error) {
	if rows == nil {
		return "", fmt.Errorf("rows is nil")
	}
	var tableQueries []string
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			zap.L().Error("error scanning table name", zap.Error(err))
			continue
		}
		tableQueries = append(
			tableQueries,
			fmt.Sprintf(
				"SELECT '%s' AS table_name, COUNT(*) AS row_count FROM %s",
				tableName,
				tableName,
			),
		)
	}

	if err := rows.Err(); err != nil {
		zap.L().Error("error iterating through table names", zap.Error(err))
		return "", err
	}

	return strings.Join(tableQueries, " UNION ALL "), nil
}
