package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis/v8"
	"github.com/huandu/go-sqlbuilder"
	"rumm-api/internal/core/domain"
	"rumm-api/kit/security"
	"time"
)

type AccountRepository struct {
	db        *sql.DB
	dbTimeout time.Duration
	jwtSecret string
	rdb       *redis.Client
}

var accountSQLStruck = sqlbuilder.NewStruct(new(sqlAccount)).For(sqlbuilder.PostgreSQL)
var accountInfoSQLStruck = sqlbuilder.NewStruct(new(accountInfo)).For(sqlbuilder.PostgreSQL)

func NewAccountRepository(db *sql.DB, dbTimeout time.Duration, jwtSecret string, rdb *redis.Client) *AccountRepository {
	return &AccountRepository{
		db:        db,
		dbTimeout: dbTimeout,
		jwtSecret: jwtSecret,
		rdb:       rdb,
	}
}

func (r AccountRepository) GetTemporalClient(ctx context.Context, storeKey string) (domain.Client, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()
	bc, err := r.rdb.Get(ctxTimeout, storeKey).Result()
	if err != nil {
		return domain.Client{}, fmt.Errorf("error trying to recover client data: %v", err)
	}
	var c domain.Client

	if err := json.Unmarshal([]byte(bc), &c); err != nil {
		return c, fmt.Errorf("error trying to recover client data: %v", err)
	}

	return c, nil
}

func (r *AccountRepository) Create(ctx context.Context, account domain.Account, client domain.Client) (*security.TokenDetails, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	c := sqlClient{
		ID: client.ID(),
		Cellphone: client.Cellphone(),
		Address: client.Address(),
		City: client.City(),
		Email: client.Email(),
		LastName: client.LastName(),
		Name: client.Name(),
		Birthday: client.BirthDay(),
	}


	createClientQuery := "INSERT INTO clients (id, name, last_name, birth_day, email, city, address, cellphone) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)"
	createAccountQuery := "INSERT INTO accounts (id, identifier, password, type_id, client_id) VALUES ($9, $10, $11, $12, $13)"

	query := fmt.Sprintf("WITH new_client as (%s) %s", createClientQuery, createAccountQuery)

	_, err := r.db.ExecContext(ctxTimeout, query, c.ID, c.Name, c.LastName, c.Birthday, c.Email, c.City, c.Address, c.Cellphone, account.ID(), account.Identifier(), account.Password(), account.AccountType(), c.ID)
	if err != nil {
		return nil, fmt.Errorf("error trying to persist account on database: %v", err)
	}

	td, err := security.CreateToken(r.jwtSecret, account.ID())
	if err != nil {
		return nil, err
	}

	return td, nil
}

func (r *AccountRepository) Authenticate(ctx context.Context, accIdentifier, password string) (domain.Account, *security.TokenDetails, error) {

	sb := accountInfoSQLStruck.SelectFrom(sqlAccountTable)
	query, args := sb.Where(sb.Equal("identifier", accIdentifier)).Build()

	ctxTimeout, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	row := r.db.QueryRowContext(ctxTimeout, query, args...)

	var acc accountInfo
	if err := row.Scan(accountInfoSQLStruck.Addr(&acc)...); err != nil {
		return domain.Account{}, nil, fmt.Errorf("error trying to find account on database, account doesn't exist: %v", err)
	}

	account, err := domain.NewAccount(
		domain.WithAccountID(acc.ID),
		domain.WithAccountIdentifier(acc.Identifier),
		domain.WithAccountHashedPass(acc.Password),
		domain.WithAccountType(acc.AccountID))
	if err != nil {
		return domain.Account{}, nil, err
	}

	isValid, err := account.ValidatePassword(password)

	if err != nil {
		return domain.Account{}, nil, err
	}
	if isValid {
		td, err := security.CreateToken(r.jwtSecret, account.ID())
		if err != nil {
			return domain.Account{}, nil, err
		}

		if err := security.CreateAuth(ctxTimeout, account.ID(), td, r.rdb); err != nil {
			return domain.Account{}, nil, err
		}

		return account, td, nil
	}

	return domain.Account{}, nil, domain.ErrAccountValidation
}

func (r *AccountRepository) Logout(ctx context.Context, accessUUID string) error {
	ctxTimeout, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	_, err := security.DeleteAuth(ctxTimeout, r.rdb, accessUUID)

	if err != nil {
		return err
	}
	return nil
}

func (r *AccountRepository) Refresh(ctx context.Context, refreshToken string) (*security.TokenDetails, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(r.jwtSecret), nil
	})
	if err != nil {
		//http.Error(w, "Refresh token expired", http.StatusUnauthorized)
		return nil, err
	}

	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		//http.Error(w, err.Error(), http.StatusUnauthorized)
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		refreshUuid, ok := claims["refresh_uuid"].(string)

		if !ok {
			//http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return nil, err
		}
		userID, ok := claims["id"].(string)
		if !ok {
			//http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return nil, err
		}

		deleted, err := security.DeleteAuth(ctxTimeout, r.rdb, refreshUuid)

		if err != nil || deleted == 0 {
			return nil, err //No autorizado
		}

		ts, err := security.CreateToken(r.jwtSecret, userID)

		if err != nil {
			return nil, err // forbidden
		}
		if err := security.CreateAuth(ctxTimeout, userID, ts, r.rdb); err != nil {
			return nil, err // forbidden
		}
		return ts, nil

	}
	return nil, security.ErrTokenCreator

}
