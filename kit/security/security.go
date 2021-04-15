package security

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"rumm-api/kit/identifier"
	"time"
)

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AccessUuid   string
	RefreshUuid  string
	AtExpires    int64
	RtExpires    int64
}

var ErrInvalidClientUUID = errors.New("hashing generator error")
var ErrTokenCreator = errors.New("can't create token")

func GetHash(password string) (string, error) {
	str := []byte(password)

	hashStr, err := bcrypt.GenerateFromPassword(str, bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("%w", ErrInvalidClientUUID)
	}

	return string(hashStr), nil
}

func ValidatePassword(hash, password string) (bool, error) {
	bch := []byte(hash)
	bcp := []byte(password)
	err := bcrypt.CompareHashAndPassword(bch, bcp)

	if err != nil {
		return false, err
	}

	return true, nil
}

func CreateToken(secret, uuid string) (TokenDetails, error) {
	td := TokenDetails{
		AtExpires: time.Now().Add(time.Minute*15).Unix(),
		AccessUuid: identifier.CreateUUID(),
		RtExpires: time.Now().Add(time.Hour*24*7).Unix(),
		RefreshUuid: identifier.CreateUUID(),
	}

	var err error
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["access_uuid"] = td.AccessUuid
	atClaims["id"] = uuid
	atClaims["exp"] = td.AtExpires

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte(secret))

	if err != nil {
		return TokenDetails{}, ErrTokenCreator
	}

	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = td.RefreshUuid
	rtClaims["id"] = uuid
	rtClaims["exp"] = td.RtExpires

	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte(secret))
	if err != nil {
		return TokenDetails{}, ErrTokenCreator
	}

	return td, nil
}
