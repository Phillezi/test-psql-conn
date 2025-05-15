package web

import (
	_ "embed"
	"net/http"
	"slices"
	"text/template"

	"github.com/Phillezi/test-psql-conn/pkg/model/table"
	"go.uber.org/zap"
)

//go:embed static/tables.html.template
var tablesTmpl string

func (s *Server) Start() error {
	s.mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	s.mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		zap.L().Debug("got request on /")
		tmpl, err := template.New("status").Parse(tablesTmpl)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		s.dataMu.RLock()
		defer s.dataMu.RUnlock()
		if err := tmpl.Execute(w, s.data); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	})

	go func() {
		defer zap.L().Info("exiting channel listener")
		var failedValidation bool
		if s.connectionChannel == nil {
			zap.L().Error("connectionChannel is nil, this is not allowed")
			failedValidation = true
		}
		if s.tablesChannel == nil {
			zap.L().Error("tablesChannel is nil, this is not allowed")
			failedValidation = true
		}
		if failedValidation {
			zap.L().Error("validation failed due to previous errors")
			return
		}

		for {
			select {
			case status := <-*s.connectionChannel:
				go func() {
					s.dataMu.Lock()
					defer s.dataMu.Unlock()
					s.data.ConnState = status
				}()
				zap.L().Debug("received status", zap.Bool("connected", status))
			case tables := <-*s.tablesChannel:
				go func() {
					s.dataMu.Lock()
					defer s.dataMu.Unlock()
					s.data.Tables = slices.Clone(tables)
				}()
				zap.L().Debug("received tables", zap.Array("tables", table.MarshalTables(tables)))
			case <-s.ctx.Done():
				s.Stop()
				return
			}
		}
	}()

	//logrus.Printf("Starting server on %s\n", s.server.Addr)
	zap.L().Info("starting server on", zap.String("address", s.server.Addr))
	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		zap.L().Error("server failed to start", zap.Error(err))
		return err
	}
	return nil
}
