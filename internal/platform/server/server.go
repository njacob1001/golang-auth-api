package server

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-redis/redis/v8"
	"log"
	"net/http"
	"os"
	"os/signal"
	"rumm-api/internal/core/service"
	"rumm-api/internal/platform/server/apimiddleware"
	"rumm-api/internal/platform/server/handler"
	"rumm-api/internal/platform/server/handler/registration"
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

	s.router.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "Origin"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	// protected endpoints
	s.router.Group(func(r chi.Router) {
		r.Use(apimiddleware.JwtAuth(s.jwtSecret, s.rdb))
		r.Post("/logout", registration.Logout(s.accountService, s.jwtSecret))
	})

	s.router.Group(func(r chi.Router) {
		r.Use(apimiddleware.SnsValidation(s.accountService))
		r.Post("/account-send-code", registration.ResendCode(s.accountService))
		r.Post("/account-verify-code", registration.Verify(s.accountService))
		r.Post("/account-register", registration.CreateAccount(s.accountService))
	})

	s.router.Group(func(r chi.Router) {
		s.router.Post("/account-init-register", registration.ValidateAccountRegister(s.accountService))
		s.router.Post("/login", registration.ValidateAccount(s.accountService))
		s.router.Post("/refresh", registration.RefreshToken(s.accountService))
		s.router.Get("/health-check", handler.HealthCheck())
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
func WithAccountService(accountService service.AccountService) Option {
	return func(server *Server) error {
		server.accountService = accountService
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
