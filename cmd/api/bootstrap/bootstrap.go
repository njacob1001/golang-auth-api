package bootstrap

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/huandu/go-sqlbuilder"
	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq"
	"rumm-api/internal/creating"
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

	creatingClientService := creating.NewClientService(clientRepository)

	ctx, srv := server.New(context.Background(), cfg.Host, cfg.Port, cfg.ShutdownTimeout, creatingClientService)
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
