package server

import (
	"context"
	"net"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Server struct {
	Router  gin.IRouter
	mServer http.Server
}

func New(host, port string) *Server {
	r := gin.Default()
	r.Use(
		cors.New(cors.Config{
			AllowOrigins: []string{"*"},
			AllowHeaders: []string{"Accept", "Authorization", "Content-Type"},
			AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			// AllowCredentials: true,
		}),
	)

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
