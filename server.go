package dewu

import (
	"context"
	"net/http"
	"time"
)

type Server struct {
	httpServer *http.Server
}

func (s *Server) Run(port string, hadler http.Handler) error {
	s.httpServer = &http.Server{
		Addr:           ":" + port,
		MaxHeaderBytes: 1 << 20,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		Handler:        hadler,
	}
	return s.httpServer.ListenAndServe()
}
func (s *Server) Close(ctx context.Context) {
	s.httpServer.Shutdown(ctx)
}
