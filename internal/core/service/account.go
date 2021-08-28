package service

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v8"
	"rumm-api/internal/core/domain"
	"rumm-api/internal/core/port"
	"rumm-api/kit/security"
	"rumm-api/kit/sender"
	"time"
)

type AccountService struct {
	accountRepository port.AccountRepository
	smsJwtSecret      string
	authJwtSecret     string
	sns               *sns.SNS
	snsTimeout        time.Duration
	Validate          *validator.Validate
	Cache             *redis.Client
}

func NewAccountService(accountRepository port.AccountRepository, sns *sns.SNS, validate *validator.Validate, cache *redis.Client, timeout time.Duration, authJwtSecret, smsJwtSecret string) AccountService {
	return AccountService{
		accountRepository: accountRepository,
		smsJwtSecret:      smsJwtSecret,
		authJwtSecret:     authJwtSecret,
		sns:               sns,
		Validate:          validate,
		Cache:             cache,
		snsTimeout:        timeout,
	}
}

func (s AccountService) CreateAccount(ctx context.Context, person domain.Person, account domain.Account, profile domain.Profile) (*security.TokenDetails, error) {

	hash, err := security.GetHash(account.Password)

	if err != nil {
		return nil, err
	}

	newAcc := domain.Account{
		ID:         account.ID,
		Identifier: account.Identifier,
		TypeID:     account.TypeID,
		Password:   string(hash),
		PersonID:   account.PersonID,
	}

	if err != nil {
		return nil, err
	}

	return s.accountRepository.Create(ctx, newAcc, profile, person)
}

func (s AccountService) Authenticate(ctx context.Context, accIdentifier, password string) (*security.TokenDetails, error) {
	return s.accountRepository.Authenticate(ctx, accIdentifier, password)
}

func (s AccountService) Logout(ctx context.Context, accessUUID string) error {
	return s.accountRepository.Logout(ctx, accessUUID)
}

func (s AccountService) Refresh(ctx context.Context, refreshToken string) (*security.TokenDetails, error) {
	return s.accountRepository.Refresh(ctx, refreshToken)
}

func (s AccountService) SendVerificationCode(ctx context.Context, phone, messageContent string) error {
	ctxTimeout, cancel := context.WithTimeout(ctx, s.snsTimeout)
	defer cancel()

	message := sender.NewMessage(phone, messageContent)
	_, err := s.sns.PublishWithContext(ctxTimeout, message)

	return err
}

func (s AccountService) RegisterCode(ctx context.Context, cellphone, finishID string) error {
	code := security.EncodeToString(4)
	expiration := time.Unix(time.Now().Add(time.Minute).Unix(), 0)

	if err := s.Cache.Set(ctx, cellphone, code, expiration.Sub(time.Now())).Err(); err != nil {
		return err
	}

	if finishID != "" {
		if err := s.Cache.Set(ctx, finishID, "success", expiration.Sub(time.Now())).Err(); err != nil {
			return err
		}
	}

	message := fmt.Sprintf("Nadie en RUMM te va a solicitar este dato. No lo compartas. Tu codigo de seguridad es %v", code)

	if err := s.SendVerificationCode(ctx, cellphone, message); err != nil {
		return err
	}

	return nil
}

func (s AccountService) VerifyAccountRegister(ctx context.Context, person domain.Person, account domain.Account, profile domain.Profile) (string, error) {

	if err := s.accountRepository.ValidateRegister(ctx, account, profile, person); err != nil {
		return "", err
	}

	if err := s.RegisterCode(ctx, person.Cellphone, ""); err != nil {
		return "", err
	}

	td, err := security.CreateSnsToken(s.smsJwtSecret, person.Cellphone, "")
	if err != nil {
		return "", err
	}

	return td.SnsToken, nil
}

func (s AccountService) GetSnsSecret() string {
	return s.smsJwtSecret
}
