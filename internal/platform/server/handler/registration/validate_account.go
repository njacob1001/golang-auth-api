package registration

import (
	"encoding/json"
	"errors"
	"net/http"
	"rumm-api/internal/core/service"
	"rumm-api/kit/identifier"
)

type validateAccountRequest struct {
	Password   string `json:"password" binding:"required"`
	Identifier string `json:"identifier" binding:"required"`
}

type authResponse struct {
	AccessToken string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func ValidateAccount(accountService service.AccountService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var req validateAccountRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		_, td,err := accountService.Authenticate(ctx, req.Identifier, req.Password)

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

		response := authResponse{
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

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		return

	}
}
