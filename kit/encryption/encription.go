package encryption

import (
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	)

var ErrInvalidClientUUID = errors.New("hashing generator error")

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