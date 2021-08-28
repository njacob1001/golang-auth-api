package registration

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"rumm-api/internal/core/service"
	"rumm-api/kit/identifier"
	"rumm-api/kit/security"
)

type requestPayload struct {
	Code string `json:"phone" validate:"required"`
}

var incorrectCode = errors.New("account already registered")

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
		fmt.Println(code)
		fmt.Println(rp.Code)
		if code != rp.Code {
			http.Error(w, incorrectCode.Error(), http.StatusBadRequest)
			return
		}

		tokenData, err := security.ExtractSnsTokenData(s.GetSnsSecret(), r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		finishID := identifier.CreateUUID()
		td, err := security.CreateSnsToken(s.GetSnsSecret(), tokenData.Cellphone, finishID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := ValidateAccountResponse{
			SnsToken: td.SnsToken,
			FinishID: td.FinishID,
		}

		j, err := json.Marshal(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if _, err := w.Write(j); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusCreated)
		return

	}
}
