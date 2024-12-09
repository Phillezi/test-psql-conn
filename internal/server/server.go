package server

import (
	"context"
	_ "embed"
	"fmt"
	"html/template"
	"net/http"

	"github.com/Phillezi/test-psql-conn/internal/models"
	"github.com/sirupsen/logrus"
)

//go:embed static/tables.html.tmpl
var tablesTmpl string

type Server struct {
	server     *http.Server
	mux        *http.ServeMux
	connStatus chan bool
	tablesChan chan []models.Table
	connState  bool
	ctx        context.Context
	tableState []models.Table
}

type TemplateData struct {
	ConnState bool
	Tables    []models.Table
}

func New(ctx context.Context, port int, connStatus chan bool, tablesChan chan []models.Table) *Server {
	mux := http.NewServeMux()
	return &Server{
		server: &http.Server{
			Addr:    fmt.Sprintf(":%d", port),
			Handler: mux,
		},
		mux:        mux,
		connStatus: connStatus,
		tablesChan: tablesChan,
		connState:  false,
		ctx:        ctx,
	}
}

func (s *Server) Start() {
	s.mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	s.mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.New("status").Parse(tablesTmpl)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		data := TemplateData{
			ConnState: s.connState,
			Tables:    s.tableState,
		}

		if err := tmpl.Execute(w, data); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	})

	go func(ctx context.Context) {
		defer fmt.Println("exiting chan listner")
		for {
			select {
			case status := <-s.connStatus:
				s.connState = status
			case tables := <-s.tablesChan:
				s.tableState = tables
			case <-ctx.Done():
				s.Stop()
				return
			}
		}
	}(s.ctx)

	logrus.Printf("Starting server on %s\n", s.server.Addr)
	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logrus.Fatalf("Server failed to start: %v\n", err)
	}
}

func (s *Server) Stop() error {
	if s.server != nil {
		logrus.Println("Stopping server...")
		return s.server.Close()
	}
	return nil
}
