package server

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v8"
	"log"
	"net/http"
	"os"
	"os/signal"
	"rumm-api/internal/core/service"
	"rumm-api/internal/platform/server/apimiddleware"
	"rumm-api/internal/platform/server/handler/registration"
	"rumm-api/internal/platform/server/routeset"
	"time"
)

type Option func(*Server) error

type Server struct {
	httpAddress     string
	router          *chi.Mux
	shutdownTimeout time.Duration
	developMode     bool
	jwtSecret       string
	rdb             *redis.Client
	validator       *validator.Validate
	//deps
	accountService service.AccountService
}

func New(ctx context.Context, options ...Option) (context.Context, Server, error) {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)

	server := Server{
		router: r,
	}
	for _, option := range options {
		err := option(&server)
		if err != nil {
			return nil, server, err
		}
	}
	server.router.Use(middleware.Recoverer)

	if server.developMode {
		server.router.Use(middleware.Logger)
	}

	server.registerRoutes()
	return serverContext(ctx), server, nil
}

func (s *Server) Run(ctx context.Context) error {
	log.Println("Server running on", s.httpAddress)

	srv := &http.Server{
		Addr:    s.httpAddress,
		Handler: s.router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Server shut down", err)
		}
	}()

	<-ctx.Done()

	ctxShutDown, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	return srv.Shutdown(ctxShutDown)
}

func (s *Server) registerRoutes() {

	// protected endpoints
	s.router.Group(func(r chi.Router) {
		r.Use(apimiddleware.JwtAuth(s.jwtSecret, s.rdb))

		r.Route("/clients", routeset.Client(s.accountService, s.validator))
	})

	s.router.Group(func(r chi.Router) {
		s.router.Post("/client", registration.CreateTemporalClient(s.accountService, s.validator))
		s.router.Post("/logout", registration.Logout(s.accountService, s.jwtSecret))
		s.router.Post("/account", registration.CreateAccount(s.accountService, s.validator))
		s.router.Post("/login", registration.ValidateAccount(s.accountService))
		s.router.Post("/refresh", registration.RefreshToken(s.accountService))
	})
}


func serverContext(ctx context.Context) context.Context {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	ctx, cancel := context.WithCancel(ctx)
	go func() {
		<-c
		cancel()
	}()

	return ctx
}

func WithJwtSecret(secret string) Option {
	return func(server *Server) error {
		server.jwtSecret = secret
		return nil
	}
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
func WithClientService(clientService service.AccountService) Option {
	return func(server *Server) error {
		server.accountService = clientService
		return nil
	}
}
func WithDevelopEnv(isDevelopMode bool) Option {
	return func(server *Server) error {
		server.developMode = isDevelopMode
		return nil
	}
}
func WithRedis(rdb *redis.Client) Option {
	return func(server *Server) error {
		server.rdb = rdb
		return nil
	}
}
func WithValidator(v *validator.Validate) Option {
	return func(server *Server) error {
		server.validator = v
		return nil
	}
}
