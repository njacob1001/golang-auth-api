package registration

import (
	"encoding/json"
	"errors"
	"github.com/go-playground/validator/v10"
	"net/http"
	"rumm-api/internal/core/service"
	"rumm-api/kit/identifier"
)


type createAccountRequest struct {
	ID          string `json:"id" validate:"required,uuid4"`
	Password    string `json:"password" validate:"required"`
	Identifier  string `json:"identifier" validate:"required"`
	AccountType string `json:"accountType" validate:"required,uuid4"`
	ClientID string `json:"clientID" validate:"required,uuid4"`
}


type createResponse struct {
	AccessToken string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}


func CreateAccount(accountService service.AccountService, validate *validator.Validate) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var req createAccountRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := validate.Struct(req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		td, err := accountService.CreateAccount(ctx, req.ID, req.Identifier, req.Password, req.AccountType, req.ClientID)

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
			AccessToken: td.AccessToken,
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
