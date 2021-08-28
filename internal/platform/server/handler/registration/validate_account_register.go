package registration

import (
	"encoding/json"
	"net/http"
	"rumm-api/internal/core/service"
)

type ValidateAccountResponse struct {
	SnsToken string `json:"sns_token"`
	FinishID string `json:"finish_id"`
}

func ValidateAccountRegister(s service.AccountService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var req createAccountRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := s.Validate.Struct(req.Account); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := s.Validate.Struct(req.Profile); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := s.Validate.Struct(req.Person); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		snsToken, err := s.VerifyAccountRegister(ctx, req.Person, req.Account, req.Profile)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		response := ValidateAccountResponse{
			SnsToken: snsToken,
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
