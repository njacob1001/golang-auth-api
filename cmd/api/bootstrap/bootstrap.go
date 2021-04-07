package bootstrap

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/huandu/go-sqlbuilder"
	_ "github.com/lib/pq"
	"rumm-api/internal/creating"
	"rumm-api/internal/platform/server"
	"rumm-api/internal/platform/storage/postgres"
	"time"
)

const (
	host            = "0.0.0.0"
	port            = 8080
	dbUser          = "rummmain"
	dbPass          = "2020rummworkerdevs"
	dbHost          = "localhost"
	dbPort          = 5432
	dbName          = "rumm"
	shutdownTimeout = 10 * time.Second
)

func Run() error {
	sqlbuilder.DefaultFlavor = sqlbuilder.PostgreSQL
	postgresURI := fmt.Sprintf("postgresql://%s@%s:%d/%s?password=%s", dbUser, dbHost, dbPort, dbName, dbPass)
	db, err := sql.Open("postgres", postgresURI)
	if err != nil {
		return err
	}
	clientRepository := postgres.NewClientRepository(db, shutdownTimeout)

	creatingClientService := creating.NewClientService(clientRepository)

	ctx, srv := server.New(context.Background(), host, port, shutdownTimeout, creatingClientService)
	return srv.Run(ctx)
}
