package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"rumm-api/internal/core/constants"
	"rumm-api/internal/core/domain"
	"rumm-api/kit/security"
	"time"
)

type AccountRepository struct {
	db        *gorm.DB
	dbTimeout time.Duration
	jwtSecret string
	rdb       *redis.Client
}

func NewAccountRepository(db *gorm.DB, dbTimeout time.Duration, jwtSecret string, rdb *redis.Client) *AccountRepository {
	return &AccountRepository{
		db:        db,
		dbTimeout: dbTimeout,
		jwtSecret: jwtSecret,
		rdb:       rdb,
	}
}

func (r *AccountRepository) Create(_ context.Context, account domain.Account, profile domain.Profile, person domain.Person) (*security.TokenDetails, error) {

	if account.TypeID == constants.ClientAccount {
		if err := r.db.Omit("company_id").Create(person).Error; err != nil {
			return nil, fmt.Errorf("error trying to persist account on database: %v", err)
		}
	} else {
		if err := r.db.Create(&person).Error; err != nil {

			return nil, fmt.Errorf("error trying to persist account on database: %v", err)
		}
	}

	if err := r.db.Omit("last_login").Create(&account).Error; err != nil {
		return nil, fmt.Errorf("error trying to persist account on database: %v", err)
	}

	if err := r.db.Create(&profile).Error; err != nil {
		return nil, fmt.Errorf("error trying to persist account on database: %v", err)
	}

	td, err := security.CreateToken(r.jwtSecret, account.ID)
	if err != nil {

		return nil, err
	}

	return td, nil
}

func (r *AccountRepository) Authenticate(ctx context.Context, accIdentifier, password, filterByType string) (*security.TokenDetails, error) {

	var acc domain.Account

	if filterByType != "" {
		if err := r.db.Where("identifier = ? AND type_id = ?", accIdentifier, filterByType).First(&acc).Error; err != nil {
			return nil, fmt.Errorf("error trying to find account on database, account doesn't exist: %v", err)
		}
	} else {
		if err := r.db.Where("identifier = ?", accIdentifier).First(&acc).Error; err != nil {
			return nil, fmt.Errorf("error trying to find account on database, account doesn't exist: %v", err)
		}
	}

	isValid, err := acc.ValidatePassword(password)

	if err != nil {
		return nil, err
	}
	if isValid {
		td, err := security.CreateToken(r.jwtSecret, acc.ID)
		if err != nil {
			return nil, err
		}

		if err := security.CreateAuth(ctx, acc.ID, td, r.rdb); err != nil {
			return nil, err
		}

		return td, nil
	}

	return nil, domain.ErrAccountValidation
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

var ErrAccountRegistered = errors.New("account already registered")
var ErrPersonRegistered = errors.New("person already registered")
var ErrProfileRegistered = errors.New("profile already registered")

func (r *AccountRepository) ValidateRegister(ctx context.Context, person domain.Person) error {
	ctxTimeout, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	pr := r.db.WithContext(ctxTimeout).Where("email = ?", person.Email).Or("cellphone = ?", person.Photo).Or("id_number = ?", person.IDNumber).Find(&person)
	if pr.RowsAffected > 0 {
		return ErrPersonRegistered
	}

	return nil
}

// ValidateAccount verify if account already exists
func (r AccountRepository) ValidateAccount(ctx context.Context, acc domain.Account) error {
	ctxTimeout, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	ar := r.db.WithContext(ctxTimeout).Where("identifier = ?", acc.Identifier).Or("id = ?", acc.ID).Find(&acc)
	if ar.RowsAffected > 0 {
		return ErrAccountRegistered
	}

	return nil

}
func (r AccountRepository) IdentifyUser(ctx context.Context, accountID string) (domain.User, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, r.dbTimeout)
	defer cancel()

	var account domain.Account
	var profile domain.Profile
	var person domain.Person


	if err := r.db.WithContext(ctxTimeout).Omit("password").Where("id = ?", accountID).Find(&account).Error; err != nil {
		return domain.User{}, err
	}


	if err := r.db.WithContext(ctxTimeout).Where("account_id = ?", accountID).Find(&profile).Error; err != nil {
		return domain.User{}, err
	}

	if err := r.db.WithContext(ctxTimeout).Where("id = ?",  account.PersonID).Find(&person).Error; err != nil {
		return domain.User{}, err
	}

	return domain.User{
		Account: account,
		Profile: profile,
		Person: person,
	}, nil

}
