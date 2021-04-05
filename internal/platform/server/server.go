package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"rumm-api/internal/creating"
	"rumm-api/internal/platform/server/handler/clients"
)

type Server struct {
	httpAddress string
	engine *gin.Engine

	//deps
	creatingClientService creating.ClientService
}

func New(host string, port uint, creatingCourseService creating.ClientService) Server {
	server := Server{
		engine: gin.New(),
		httpAddress: fmt.Sprintf("%s:%d", host, port),
		creatingClientService: creatingCourseService,
	}
	server.registerRoutes()
	return server
}

func (server *Server) Run() error {
	log.Println("Server running on", server.httpAddress)
	return server.engine.Run(server.httpAddress)
}

func (server *Server) registerRoutes() {
	server.engine.POST("/clients", clients.CreateHandler(server.creatingClientService))
}