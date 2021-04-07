package server

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"rumm-api/internal/creating"
	"rumm-api/internal/platform/server/handler/clients"
	"time"
)

type Server struct {
	httpAddress string
	engine      *gin.Engine
	shutdownTimeout time.Duration


	//deps
	creatingClientService creating.ClientService
}

func New(ctx context.Context, host string, port uint, shutdownTimeout time.Duration,creatingCourseService creating.ClientService) (context.Context, Server) {
	server := Server{
		engine:                gin.New(),
		httpAddress:           fmt.Sprintf("%s:%d", host, port),

		shutdownTimeout: shutdownTimeout,
		creatingClientService: creatingCourseService,

	}
	server.registerRoutes()
	return serverContext(ctx), server
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
	server.engine.POST("/clients", clients.CreateHandler(server.creatingClientService))
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
