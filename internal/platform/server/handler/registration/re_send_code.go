package registration

import (
	"errors"
	"net/http"
	"rumm-api/internal/core/constants"
	"rumm-api/internal/core/service"
	"rumm-api/kit/security"
	"strconv"
	"time"
)

var ErrSendCode = errors.New("can't send more code")
var ErrCodeNoExist = errors.New("missing code status")

func ResendCode(s service.AccountService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		data, err := security.ExtractSnsTokenData(s.GetSnsSecret(), r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		value, err := s.Cache.Get(ctx, data.AccessID).Result()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		i, err := strconv.Atoi(value)
		if err != nil {
			http.Error(w, ErrCodeNoExist.Error(), http.StatusBadRequest)
			return
		}

		if i == constants.ConfirmationCodeThirdIntent || i == constants.ConfirmationFailure || i == constants.ConfirmationSuccess {
			http.Error(w, ErrSendCode.Error(), http.StatusBadRequest)
			return
		}

		if err := s.Cache.Set(ctx, data.AccessID, i+1, time.Hour).Err(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := s.RegisterCode(ctx, data.Cellphone); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusCreated)
		return
	}
}
