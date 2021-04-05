package bootstrap

import (
	"database/sql"
	"fmt"
	"github.com/huandu/go-sqlbuilder"
	_ "github.com/lib/pq"
	"rumm-api/internal/creating"
	"rumm-api/internal/platform/server"
	"rumm-api/internal/platform/storage/postgres"
)


const (
	host = "0.0.0.0"
	port = 8080
	dbUser = "rummmain"
	dbPass = "2020rummworkerdevs"
	dbHost = "localhost"
	dbPort = 5432
	dbName = "rumm"
)

func Run() error {
	sqlbuilder.DefaultFlavor = sqlbuilder.PostgreSQL
	postgresURI := fmt.Sprintf("postgresql://%s@%s:%d/%s?password=%s", dbUser, dbHost, dbPort, dbName, dbPass)
	db, err := sql.Open("postgres", postgresURI)
	if err != nil {
		return err
	}
	clientRepository := postgres.NewClientRepository(db)

	creatingClientService := creating.NewClientService(clientRepository)

	srv := server.New(host, port, creatingClientService)
	return srv.Run()
}
