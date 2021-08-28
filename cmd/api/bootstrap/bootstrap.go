package bootstrap

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
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

	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("us-east-2"),
	}))

	svc := sns.New(sess)

	accountRepository := postgresdb.NewAccountRepository(db, cfg.DbTimeout, cfg.JwtSecret, rdb)

	isDevelopMode := cfg.Mode == "DEVELOP"

	validate := validator.New()

	accountService := service.NewAccountService(accountRepository, svc, validate, rdb, cfg.SnsTimeout, cfg.JwtSecret, cfg.SmsJwtSecret)

	ctx, srv, err := server.New(
		context.Background(),
		server.WithTimeout(cfg.ShutdownTimeout),
		server.WithAddress(cfg.Host, cfg.Port),
		server.WithAccountService(accountService),
		server.WithDevelopEnv(isDevelopMode),
		server.WithJwtSecret(cfg.JwtSecret),
		server.WithRedis(rdb))

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
	SmsJwtSecret    string        `default:"example_secret" split_words:"true"`
	CacheIndex      int           `default:"1" split_words:"true"`
	CacheHost       string        `default:"0.0.0.0" split_words:"true"`
	CachePort       uint          `default:"6379" split_words:"true"`
	SnsTimeout      time.Duration `default:"5s" split_words:"true"`
}
