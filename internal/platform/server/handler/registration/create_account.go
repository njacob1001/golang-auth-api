package registration

import (
	"encoding/json"
	"errors"
	"net/http"
	"rumm-api/internal/core/domain"
	"rumm-api/internal/core/service"
	"rumm-api/kit/identifier"
)

type createAccountRequest struct {
	Person  domain.Person  `json:"person"`
	Account domain.Account `json:"account"`
	Profile domain.Profile `json:"profile"`
}

type createResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func CreateAccount(s service.AccountService) http.HandlerFunc {
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

		td, err := s.CreateAccount(ctx, req.Person, req.Account, req.Profile)

		if err != nil {
			switch {
			case errors.Is(err, identifier.ErrInvalidClientUUID):
				http.Error(w, err.Error(), http.StatusBadRequest)

				return
			default:
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		response := createResponse{
			AccessToken:  td.AccessToken,
			RefreshToken: td.RefreshToken,
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
