package bootstrap

import (
	"context"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v8"
	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"rumm-api/internal/core/service"
	"rumm-api/internal/platform/server"
	postgresdb "rumm-api/internal/platform/storage/postgres"
	"time"
)

func Run() error {
	var cfg configEnv
	err := envconfig.Process("RUMM", &cfg)
	if err != nil {
		return err
	}

	postgresURI := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=disable TimeZone=-5", cfg.DbHost, cfg.DbUser, cfg.DbPass, cfg.DbName, cfg.DbPort)

	//postgresURI := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable", cfg.DbUser, authenticationToken, cfg.DbHost, cfg.DbPort, cfg.DbName)
	//sqlDB, err := sql.Open("postgres", postgresURI)

	db, err := gorm.Open(postgres.Open(postgresURI), &gorm.Config{})
	if err != nil {
		return err
	}


	//if err := sqlDB.Ping(); err != nil {
	//	panic("Ping error: "+err.Error())
	//}


	redisURI := fmt.Sprintf("%v:%v", cfg.RdbHost, cfg.RdbPort)
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisURI,
		Password: cfg.RdbPassword,
		DB:       cfg.RdbIndex,
	})


	accountRepository := postgresdb.NewAccountRepository(db, cfg.DbTimeout, cfg.JwtSecret, rdb)

	accountService := service.NewAccountService(accountRepository)

	isDevelopMode := !(cfg.ServerMode == "release")

	validate := validator.New()

	ctx, srv, err := server.New(
		context.Background(),
		server.WithTimeout(cfg.ShutdownTimeout),
		server.WithAddress(cfg.Host, cfg.Port),
		server.WithAccountService(accountService),
		server.WithDevelopEnv(isDevelopMode),
		server.WithJwtSecret(cfg.JwtSecret),
		server.WithRedis(rdb),
		server.WithValidator(validate))

	if err != nil {
		return err
	}
	return srv.Run(ctx)
}

type configEnv struct {
	// Server configuration
	Host            string        `default:"0.0.0.0" split_words:"true"`
	Port            uint          `default:"8080" split_words:"true"`
	ShutdownTimeout time.Duration `default:"10s" split_words:"true"`

	// Database configuration
	DbUser     string        `required:"true" split_words:"true"`
	DbPass     string        `required:"true" split_words:"true"`
	DbHost     string        `required:"true" split_words:"true"`
	DbPort     uint          `required:"true" split_words:"true"`
	DbName     string        `required:"true" split_words:"true"`
	DbTimeout  time.Duration `default:"5s" split_words:"true"`
	ServerMode string        `default:"develop" split_words:"true"`

	// authentication
	JwtSecret string `required:"true" split_words:"true"`

	// Redis database
	RdbIndex    int    `default:"0" split_words:"true"`
	RdbPassword string `default:"" split_words:"true"`
	RdbHost     string `default:"0.0.0.0" split_words:"true"`
	RdbPort     uint   `default:"6379" split_words:"true"`
}
