package bootstrap

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go-v2/service/sns"
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

type databaseCredential struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

type seed struct {
	Auth string `json:"auth"`
	Sns  string `json:"sns"`
}

func Run() error {
	secretName := "development/app-db"
	seedsSecretName := "seeds"
	var cfg configEnv
	err := envconfig.Process("AUTH_API", &cfg)
	if err != nil {
		return err
	}

	awsCfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(cfg.Region),
	)
	if err != nil {
		return err
	}

	svc := sns.NewFromConfig(awsCfg)

	secretClient := secretsmanager.New(secretsmanager.Options{
		Region:      cfg.Region,
		Credentials: awsCfg.Credentials,
	})

	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secretName),
		VersionStage: aws.String("AWSCURRENT"),
	}

	seedsInput := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(seedsSecretName),
		VersionStage: aws.String("AWSCURRENT"),
	}

	secretResult, err := secretClient.GetSecretValue(context.Background(), input)
	if err != nil {
		return err
	}

	seedsResult, err := secretClient.GetSecretValue(context.Background(), seedsInput)
	if err != nil {
		return err
	}

	var secretString string

	var dbc databaseCredential

	if secretResult.SecretString != nil {
		secretString = *secretResult.SecretString
		if err := json.Unmarshal([]byte(secretString), &dbc); err != nil {
			return err
		}
	}

	var seedsSecretString string
	var sr seed

	if seedsResult.SecretString != nil {
		seedsSecretString = *seedsResult.SecretString
		if err := json.Unmarshal([]byte(seedsSecretString), &sr); err != nil {
			return err
		}
	}

	postgresURI := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=disable TimeZone=-5", cfg.DbHost, dbc.Username, dbc.Password, cfg.DbName, cfg.DbPort)

	db, err := gorm.Open(postgres.Open(postgresURI), &gorm.Config{})
	if err != nil {
		return err
	}

	redisURI := fmt.Sprintf("%v:%v", cfg.CacheHost, cfg.CachePort)
	rdb := redis.NewClient(&redis.Options{
		Addr: redisURI,
		DB:   cfg.CacheIndex,
	})

	accountRepository := postgresdb.NewAccountRepository(db, cfg.DbTimeout, sr.Auth, rdb)

	isDevelopMode := cfg.Mode == "DEVELOP"

	validate := validator.New()

	accountService := service.NewAccountService(accountRepository, svc, validate, rdb, cfg.SnsTimeout, sr.Auth, sr.Sns)

	ctx, srv, err := server.New(
		context.Background(),
		server.WithTimeout(cfg.ShutdownTimeout),
		server.WithAddress(cfg.Host, cfg.Port),
		server.WithAccountService(accountService),
		server.WithDevelopEnv(isDevelopMode),
		server.WithJwtSecret(sr.Auth),
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
	Region          string        `default:"us-east-2" split_words:"true"`
	DbUser          string        `default:"admin" split_words:"true"`
	DbHost          string        `default:"0.0.0.0" split_words:"true"`
	DbPort          uint          `default:"5432" split_words:"true"`
	DbName          string        `default:"develop" split_words:"true"`
	DbSchema        string        `default:"public" split_words:"true"`
	DbTimeout       time.Duration `default:"5s" split_words:"true"`
	CacheIndex      int           `default:"1" split_words:"true"`
	CacheHost       string        `default:"0.0.0.0" split_words:"true"`
	CachePort       uint          `default:"6379" split_words:"true"`
	SnsTimeout      time.Duration `default:"5s" split_words:"true"`
}
