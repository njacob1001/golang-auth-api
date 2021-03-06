package service

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/sns"
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
	sns               *sns.Client
	snsTimeout        time.Duration
	Validate          *validator.Validate
	Cache             *redis.Client
}

func NewAccountService(accountRepository port.AccountRepository, sns *sns.Client, validate *validator.Validate, cache *redis.Client, timeout time.Duration, authJwtSecret, smsJwtSecret string) AccountService {
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

func (s AccountService) Authenticate(ctx context.Context, accIdentifier, password, filterByType string) (*security.TokenDetails, error) {
	return s.accountRepository.Authenticate(ctx, accIdentifier, password, filterByType)
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

	_, err := s.sns.Publish(ctxTimeout, message)

	return err
}

func (s AccountService) RegisterCode(ctx context.Context, cellphone string) error {
	code := security.EncodeToString(4)

	if err := s.Cache.Set(ctx, cellphone, code, time.Minute*5).Err(); err != nil {
		return err
	}

	message := fmt.Sprintf("Nadie en RUMM te va a solicitar este dato. No lo compartas. Tu codigo de seguridad es %v", code)

	if err := s.SendVerificationCode(ctx, cellphone, message); err != nil {
		return err
	}

	return nil
}

func (s AccountService) VerifyAccountRegister(ctx context.Context, person domain.Person) error {
	return s.accountRepository.ValidateRegister(ctx, person)
}

func (s AccountService) ValidateAccountExists(ctx context.Context, acc domain.Account) error {
	return s.accountRepository.ValidateAccount(ctx, acc)
}


func (s AccountService) RegisterSnsToken(cellphone string) (security.SnsTokenDetails, error) {
	return security.CreateSnsToken(s.smsJwtSecret, cellphone)
}

func (s AccountService) GetSnsSecret() string {
	return s.smsJwtSecret
}

func (s AccountService) IdentifyUser(ctx context.Context, accountID string) (domain.User, error) {
	return s.accountRepository.IdentifyUser(ctx, accountID)
}