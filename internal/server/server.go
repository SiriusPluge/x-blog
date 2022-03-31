package server

import (
	"context"
	"net/http"
	"time"

	"github.com/spf13/viper"
)

type Server struct {
	httpServer *http.Server
}

const defaultPort = ":8080"

func (s *Server) RunServer(handler http.Handler) error {
	port := viper.GetString("port")
	if port == "" {
		port = defaultPort
	}

	s.httpServer = &http.Server{
		Addr:           ":" + port,
		Handler:        handler,
		MaxHeaderBytes: 1 << 20, // 1 MB
		ReadTimeout:    time.Duration(viper.GetInt("server.ReadTimeout")) * time.Second,
		WriteTimeout:   time.Duration(viper.GetInt("server.WriteTimeout")) * time.Second,
	}

	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx *context.Context) error {
	return s.httpServer.Shutdown(*ctx)
}
