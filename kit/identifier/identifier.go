package identifier

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
)

var ErrInvalidClientUUID = errors.New("invalid Client UUID")
var ErrCreatingClientUUID = errors.New("can't create UUID")

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

func NewIdentifier() Identifier {
	value := uuid.New()

	return Identifier{
		String: value.String(),
	}

}
