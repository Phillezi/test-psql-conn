package psql

import (
	"github.com/Phillezi/test-psql-conn/pkg/model/table"
	"go.uber.org/zap"
)

func (c *Client) sendStatusAsync(status bool) {
	go func() {
		select {
		case c.connStatus <- status:
		default:
			zap.L().Warn("failed to update status, channel is full")
		}
	}()
}

func (c *Client) sendTablesAsync(tables []table.Table) {
	go func() {
		select {
		case c.tables <- tables:
		default:
			zap.L().Warn("failed to update tables, channel is full")
		}
	}()
}
