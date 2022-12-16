package server

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(handler http.Handler) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:    "localhost:8080",
			Handler: handler,
		},
	}
}

func (server *Server) Start() {
	err := server.httpServer.ListenAndServe()
	if err != nil {
		log.Fatal("[DEBUG11] Server starting failed")
	}
}

func (server *Server) Stop(ctx *gin.Context) {
	err := server.httpServer.Shutdown(ctx)
	if err != nil {
		log.Fatal("[DEBUG12] Server stop failed")
	}
}
