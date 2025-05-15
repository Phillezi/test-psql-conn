package psql

import "github.com/Phillezi/test-psql-conn/pkg/model/table"

func (c *Client) ConnectChannel() *chan bool {
	return &c.connStatus
}

func (c *Client) TableChannel() *chan []table.Table {
	return &c.tables
}
