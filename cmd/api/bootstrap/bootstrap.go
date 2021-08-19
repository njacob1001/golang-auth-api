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
	err := envconfig.Process("AUTH_API", &cfg)
	if err != nil {
		return err
	}

	postgresURI := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=disable TimeZone=-5", cfg.DbHost, cfg.DbUser, cfg.DbPass, cfg.DbName, cfg.DbPort)

	db, err := gorm.Open(postgres.Open(postgresURI), &gorm.Config{})
	if err != nil {
		return err
	}

	redisURI := fmt.Sprintf("%v:%v", cfg.CacheHost, cfg.CachePort)
	rdb := redis.NewClient(&redis.Options{
		Addr: redisURI,
		DB:   cfg.CacheIndex,
	})

	accountRepository := postgresdb.NewAccountRepository(db, cfg.DbTimeout, cfg.JwtSecret, rdb)

	accountService := service.NewAccountService(accountRepository)

	isDevelopMode := !(cfg.Mode == "release")

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
	Host            string        `default:"0.0.0.0" split_words:"true"`
	Port            uint          `default:"80" split_words:"true"`
	ShutdownTimeout time.Duration `default:"10s" split_words:"true"`
	Mode            string        `default:"DEVELOP" split_words:"true"`
	DbUser          string        `default:"admin" split_words:"true"`
	DbPass          string        `default:"admin123" split_words:"true"`
	DbHost          string        `default:"0.0.0.0" split_words:"true"`
	DbPort          uint          `default:"5432" split_words:"true"`
	DbName          string        `default:"develop" split_words:"true"`
	DbSchema        string        `default:"public" split_words:"true"`
	DbTimeout       time.Duration `default:"5s" split_words:"true"`
	JwtSecret       string        `default:"example_secret" split_words:"true"`
	CacheIndex      int           `default:"1" split_words:"true"`
	CacheHost       string        `default:"0.0.0.0" split_words:"true"`
	CachePort       uint          `default:"6379" split_words:"true"`
}
