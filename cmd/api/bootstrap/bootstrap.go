package bootstrap

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/huandu/go-sqlbuilder"
	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq"
	"rumm-api/internal/core/services/clients"

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

	clientRepository := postgres.NewClientRepository(db, cfg.DbTimeout)
	clientService := clients.NewClientService(clientRepository)


	ctx, srv, err := server.NewServer(
		context.Background(),
		server.WithTimeout(cfg.ShutdownTimeout),
		server.WithAddress(cfg.Host, cfg.Port),
		server.WithClientService(clientService),
	)

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
	DbUser    string        `required:"true"`
	DbPass    string        `required:"true"`
	DbHost    string        `required:"true"`
	DbPort    uint          `required:"true"`
	DbName    string        `required:"true"`
	DbTimeout time.Duration `default:"5s"`
}
