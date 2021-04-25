package bootstrap

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v8"
	"github.com/huandu/go-sqlbuilder"
	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq"
	"rumm-api/internal/core/service"
	"rumm-api/internal/platform/server"
	"rumm-api/internal/platform/storage/postgres"
	"time"
)

func Run() error {
	var cfg config
	err := envconfig.Process("RUMM", &cfg)
	if err != nil {
		return nil
	}

	sqlbuilder.DefaultFlavor = sqlbuilder.PostgreSQL

	postgresURI := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable", cfg.DbUser, cfg.DbPass, cfg.DbHost, cfg.DbPort, cfg.DbName)
	db, err := sql.Open("postgres", postgresURI)
	if err != nil {
		return err
	}

	redisURI := fmt.Sprintf("%v:%v", cfg.RdbHost, cfg.RdbPort)
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisURI,
		Password: cfg.RdbPassword,
		DB:       cfg.RdbIndex,
	})

	clientRepository := postgres.NewClientRepository(db, cfg.DbTimeout, rdb)
	accountRepository := postgres.NewAccountRepository(db, cfg.DbTimeout, cfg.JwtSecret, rdb)

	accountService := service.NewAccountService(accountRepository, clientRepository)

	isDevelopMode := !(cfg.ServerMode == "release")

	validate := validator.New()

	ctx, srv, err := server.New(
		context.Background(),
		server.WithTimeout(cfg.ShutdownTimeout),
		server.WithAddress(cfg.Host, cfg.Port),
		server.WithClientService(accountService),
		server.WithDevelopEnv(isDevelopMode),
		server.WithJwtSecret(cfg.JwtSecret),
		server.WithRedis(rdb),
		server.WithValidator(validate))

	if err != nil {
		return err
	}
	return srv.Run(ctx)
}

type config struct {
	// Server configuration
	Host            string        `default:"0.0.0.0"`
	Port            uint          `default:"8080"`
	ShutdownTimeout time.Duration `default:"10s"`

	// Database configuration
	DbUser     string        `required:"true"`
	DbPass     string        `required:"true"`
	DbHost     string        `required:"true"`
	DbPort     uint          `required:"true"`
	DbName     string        `required:"true"`
	DbTimeout  time.Duration `default:"5s"`
	ServerMode string        `default:"develop"`

	// authentication
	JwtSecret string `required:"true"`

	// Redis database
	RdbIndex    int    `default:"0"`
	RdbPassword string `default:""`
	RdbHost     string `default:"0.0.0.0"`
	RdbPort     uint   `default:"6379"`
}
