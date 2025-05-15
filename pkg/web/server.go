package web

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/Phillezi/test-psql-conn/pkg/model/table"
	"github.com/Phillezi/test-psql-conn/util"
)

type Server struct {
	ctx context.Context

	server *http.Server
	mux    *http.ServeMux

	dataMu sync.RWMutex
	data   TemplateData

	connectionChannel *chan bool
	tablesChannel     *chan []table.Table
}

type ServerOpts struct {
	Port              *int
	ConnectionChannel *chan bool
	TablesChannel     *chan []table.Table
}

func New(opts ...ServerOpts) *Server {
	var (
		port              int
		connectionChannel *chan bool
		tablesChannel     *chan []table.Table
	)
	if len(opts) > 0 {
		util.IfNotNilDo(opts[0].Port, func() { port = *opts[0].Port })
		util.IfNotNilDo(opts[0].ConnectionChannel, func() { connectionChannel = opts[0].ConnectionChannel })
		util.IfNotNilDo(opts[0].TablesChannel, func() { tablesChannel = opts[0].TablesChannel })
	}

	mux := http.NewServeMux()
	return &Server{
		ctx: context.Background(),

		server: &http.Server{
			Addr:    fmt.Sprintf(":%d", util.Or(port, DefaultHttpPort)),
			Handler: mux,
		},
		mux: mux,

		connectionChannel: connectionChannel,
		tablesChannel:     tablesChannel,
	}
}

func (s *Server) WithContext(ctx context.Context) *Server {
	s.ctx = ctx
	return s
}
