package server

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"rumm-api/internal/core/services/clients"
	"rumm-api/internal/platform/server/handler/clientsHandler"
	"time"
)

type Option func(*Server) error

type Server struct {
	httpAddress string
	engine      *gin.Engine
	shutdownTimeout time.Duration


	//deps
	creatingClientService clients.ClientService
}

func NewServer(ctx context.Context, options ...Option) (context.Context, Server, error) {
	server := Server{
		engine:                gin.New(),
	}
	for _, option := range options {
		err := option(&server)
		if err != nil {
			return nil, server, err
		}
	}
	server.engine.Use(gin.Recovery())
	// debug mode
	//server.engine.Use(gin.Logger())
	server.registerRoutes()
	return serverContext(ctx), server, nil
}

func (server *Server) Run(ctx context.Context) error {
	log.Println("Server running on", server.httpAddress)

	srv := &http.Server{
		Addr:    server.httpAddress,
		Handler: server.engine,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Server shut down", err)
		}
	}()

	<- ctx.Done()

	ctxShutDown, cancel := context.WithTimeout(context.Background(), server.shutdownTimeout)
	defer cancel()

	return srv.Shutdown(ctxShutDown)

}

func (server *Server) registerRoutes() {
	server.engine.POST("/clients", clientsHandler.CreateHandler(server.creatingClientService))
}

func serverContext(ctx context.Context) context.Context {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	ctx, cancel := context.WithCancel(ctx)
	go func(){
		<-c
		cancel()
	}()

	return ctx
}

func WithAddress(host string, port uint) Option {
	return func(server *Server) error {
		server.httpAddress = fmt.Sprintf("%s:%d", host, port)
		return nil
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(server *Server) error {
		server.shutdownTimeout = timeout
		return nil
	}
}

func WithClientService(clientService clients.ClientService) Option {
	return func(server *Server) error {
		server.creatingClientService = clientService
		return nil
	}
}