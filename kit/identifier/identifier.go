package identifier

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
)

var ErrInvalidClientUUID = errors.New("invalid Client UUID")

type Identifier struct {
	String string
}

func ValidateIdentifier(value string) (Identifier, error) {
	v, err := uuid.Parse(value)
	if err != nil {
		return Identifier{}, fmt.Errorf("%w: %s", ErrInvalidClientUUID, value)
	}

	return Identifier{
		String: v.String(),
	}, nil
}

func CreateUUID() string {
	v := uuid.New()
	return v.String()
}
