package web

import "github.com/Phillezi/test-psql-conn/pkg/model/table"

type TemplateData struct {
	ConnState bool
	Tables    []table.Table
}
