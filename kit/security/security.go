package security

import (
	"context"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis/v8"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"rumm-api/kit/identifier"
	"strings"
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

func CreateToken(secret, uuid string) (*TokenDetails, error) {
	td := new(TokenDetails)
	td.AtExpires = time.Now().Add(time.Minute * 15).Unix()
	td.AccessUuid = identifier.CreateUUID()
	td.RtExpires = time.Now().Add(time.Hour * 24 * 7).Unix()
	td.RefreshUuid = identifier.CreateUUID()

	var err error
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["access_uuid"] = td.AccessUuid
	atClaims["id"] = uuid
	atClaims["exp"] = td.AtExpires

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte(secret))

	if err != nil {
		return new(TokenDetails), ErrTokenCreator
	}

	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = td.RefreshUuid
	rtClaims["id"] = uuid
	rtClaims["exp"] = td.RtExpires

	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte(secret))
	if err != nil {
		return new(TokenDetails), ErrTokenCreator
	}

	return td, nil
}

func CreateAuth(ctx context.Context, accountID string, td *TokenDetails, rdb *redis.Client) error {

	at := time.Unix(td.AtExpires, 0)
	rt := time.Unix(td.RtExpires, 0)
	now := time.Now()

	errAccess := rdb.Set(ctx, td.AccessUuid, accountID, at.Sub(now)).Err()
	if errAccess != nil {
		return errAccess
	}

	errRefresh := rdb.Set(ctx, td.RefreshUuid, accountID, rt.Sub(now)).Err()
	if errRefresh != nil {
		return errRefresh
	}
	return nil
}

func ExtractToken(r *http.Request) string {
	bearToken := r.Header.Get("Authorization")
	//normally Authorization the_token_xxx
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

func VerifyToken(secret string, r *http.Request) (*jwt.Token, error) {
	tokenString := ExtractToken(r)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func TokenValid(secret string, r *http.Request) error {
	token, err := VerifyToken(secret, r)
	if err != nil {
		return err
	}
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return err
	}
	return nil
}

type AccessDetails struct {
	AccessUuid string
	UserID     string
}

func ExtractTokenMetadata(secret string, r *http.Request) (*AccessDetails, error) {
	token, err := VerifyToken(secret, r)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		accessUuid, ok := claims["access_uuid"].(string)
		if !ok {
			return nil, err
		}
		userId := fmt.Sprintf("%s", claims["id"])

		return &AccessDetails{
			AccessUuid: accessUuid,
			UserID:   userId,
		}, nil
	}
	return nil, err
}

func FetchAuth(ctx context.Context, authD *AccessDetails, rdb *redis.Client) (string, error) {
	userid, err := rdb.Get(ctx, authD.AccessUuid).Result()
	if err != nil {
		return "", err
	}
	return userid, nil
}

func DeleteAuth(ctx context.Context, rdb *redis.Client, givenUUID string) (int64, error) {
	deleted, err := rdb.Del(ctx, givenUUID).Result()
	if err != nil {
		return 0, err
	}

	return deleted, nil
}