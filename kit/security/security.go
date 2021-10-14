package security

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis/v8"
	"golang.org/x/crypto/bcrypt"
	"io"
	"net/http"
	"rumm-api/kit/identifier"
	"strings"
	"time"
)

type TokenDetails struct {
	ID           string
	AccessToken  string
	RefreshToken string
	AccessUuid   string
	RefreshUuid  string
	AtExpires    int64
	RtExpires    int64
}

var ErrInvalidClientUUID = errors.New("hashing generator error")
var ErrTokenCreator = errors.New("can't create token")

func GetHash(password string) ([]byte, error) {
	str := []byte(password)

	hashStr, err := bcrypt.GenerateFromPassword(str, bcrypt.DefaultCost)
	if err != nil {
		return []byte{}, fmt.Errorf("%w", ErrInvalidClientUUID)
	}

	return hashStr, nil
}

func ValidatePassword(hash []byte, password string) (bool, error) {
	bch := hash
	bcp := []byte(password)

	err := bcrypt.CompareHashAndPassword(bch, bcp)

	if err == nil {
		return true, nil
	}

	return false, err
}

type SnsTokenDetails struct {
	SnsToken  string
	Cellphone string
	AccessID  string
	AtExpires int64
}

func CreateSnsToken(secret, phone string) (SnsTokenDetails, error) {
	td := SnsTokenDetails{}
	td.AtExpires = time.Now().Add(time.Minute * 30).Unix()
	td.Cellphone = phone

	var err error
	id := identifier.CreateUUID()
	td.AccessID = id
	atClaims := jwt.MapClaims{}
	atClaims["cellphone"] = phone
	atClaims["id"] = id
	atClaims["exp"] = td.AtExpires

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.SnsToken, err = at.SignedString([]byte(secret))
	if err != nil {
		return SnsTokenDetails{}, ErrTokenCreator
	}
	return td, nil
}

func CreateToken(secret, uuid string) (*TokenDetails, error) {
	td := new(TokenDetails)
	td.ID = uuid
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

func IsTokenValid(secret string, r *http.Request) error {
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
			UserID:     userId,
		}, nil
	}
	return nil, err
}

func ExtractSnsTokenData(secret string, r *http.Request) (*SnsTokenDetails, error) {
	token, err := VerifyToken(secret, r)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		cellphone, ok := claims["cellphone"].(string)
		if !ok {
			return nil, err
		}

		accessID, ok := claims["id"].(string)
		if !ok {
			return nil, err
		}

		return &SnsTokenDetails{
			AccessID:  accessID,
			Cellphone: cellphone,
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

var table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}

func EncodeToString(max int) string {
	b := make([]byte, max)
	n, err := io.ReadAtLeast(rand.Reader, b, max)
	if n != max {
		panic(err)
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	return string(b)
}
