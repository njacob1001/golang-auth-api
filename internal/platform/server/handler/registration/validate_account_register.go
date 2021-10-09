package registration

import (
	"encoding/json"
	"net/http"
	"rumm-api/internal/core/constants"
	"rumm-api/internal/core/service"
	"time"
)

type ValidateAccountResponse struct {
	Token string `json:"token"`
}

func ValidateAccountRegister(s service.AccountService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var req createAccountRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		// validate person schema
		if err := s.Validate.Struct(req.Person); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// check if person already exist
		if err := s.VerifyAccountRegister(ctx, req.Person); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// send code and save in cache
		if err := s.RegisterCode(ctx, req.Person.Cellphone); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Create token
		td, err := s.RegisterSnsToken(req.Person.Cellphone)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Init status of confirmation
		if err := s.Cache.Set(ctx, td.AccessID, constants.ConfirmationCodeInit, time.Hour).Err(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		response := ValidateAccountResponse{
			Token: td.SnsToken,
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

func ValidateAccountExists(s service.AccountService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var req createAccountRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		// validate person schema
		if err := s.Validate.Struct(req.Account); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// check if person already exist
		if err := s.ValidateAccountExists(ctx, req.Account); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("OK"))

		return

	}
}