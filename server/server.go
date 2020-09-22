package server

import (
	"context"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"net"
	"net/http"
	"time"
)

type Server struct {
	Router  chi.Router
	mServer http.Server
}

func New(host, port string) *Server {
	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		//AllowCredentials: true,
	}))
	r.Use(middleware.Logger)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.AllowContentType("application/json"))
	r.Use(middleware.Timeout(60 * time.Second))

	return &Server{
		Router: r,
		mServer: http.Server{
			Addr:    net.JoinHostPort(host, port),
			Handler: r,
		},
	}
}

func (server *Server) ListenAndServe() error {
	return server.mServer.ListenAndServe()
}

func (server *Server) Stop() error {
	return server.mServer.Shutdown(context.Background())
}
