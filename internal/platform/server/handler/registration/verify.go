package registration

import (
	"encoding/json"
	"errors"
	"net/http"
	"rumm-api/internal/core/constants"
	"rumm-api/internal/core/service"
	"rumm-api/kit/security"
	"time"
)

type requestPayload struct {
	Code string `json:"code" validate:"required"`
}

var incorrectCode = errors.New("incorrect code")
var ErrStatusOutOfRange = errors.New("status is out of range")

func Verify(s service.AccountService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var rp requestPayload
		ctx := r.Context()

		err := security.IsTokenValid(s.GetSnsSecret(), r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		if err := json.NewDecoder(r.Body).Decode(&rp); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := s.Validate.Struct(rp); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		data, err := security.ExtractSnsTokenData(s.GetSnsSecret(), r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		code, err := s.Cache.Get(ctx, data.Cellphone).Result()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		status, err := s.Cache.Get(ctx, data.AccessID).Int()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if status <= constants.ConfirmationCodeThirdIntent {
			if code != rp.Code {
				http.Error(w, incorrectCode.Error(), http.StatusBadRequest)
				return
			}

			if err := s.Cache.Set(ctx, data.AccessID, constants.ConfirmationSuccess, time.Minute*20).Err(); err != nil {
				http.Error(w, incorrectCode.Error(), http.StatusBadRequest)
				return
			}

			w.WriteHeader(http.StatusCreated)
			return
		}

		http.Error(w, ErrStatusOutOfRange.Error(), http.StatusBadRequest)
		return

	}
}
