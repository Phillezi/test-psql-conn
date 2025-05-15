package web

import "go.uber.org/zap"

func (s *Server) Stop() error {
	if s.server != nil {
		zap.L().Info("stopping server")
		return s.server.Close()
	}
	return nil
}
