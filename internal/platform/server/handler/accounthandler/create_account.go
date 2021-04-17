package accounthandler

import (
	"encoding/json"
	"errors"
	"net/http"
	service "rumm-api/internal/core/services/clients"
	"rumm-api/kit/identifier"
)

type createAccountRequest struct {
	ID          string `json:"id" binding:"required"`
	Password    string `json:"password" binding:"required"`
	Identifier  string `json:"identifier" binding:"required"`
	AccountType string `json:"accountType" binding:"required"`
}

func CreateAccountHandler(accountService service.AccountService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var req createAccountRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err := accountService.CreateAccount(ctx, req.ID, req.Identifier, req.Password, req.AccountType)

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
		w.WriteHeader(http.StatusCreated)

	}
}
