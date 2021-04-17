package registration

import (
	"encoding/json"
	"net/http"
	"rumm-api/internal/core/service"
	"rumm-api/kit/security"
)

type refreshResponse struct {
	AccessToken string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func RefreshToken(accountService service.AccountService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		refreshToken := security.ExtractToken(r)
		ctx := r.Context()

		td, err := accountService.Refresh(ctx, refreshToken)

		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		response := refreshResponse{
			AccessToken: td.AccessToken,
			RefreshToken: td.RefreshToken,
		}

		j, err := json.Marshal(response)

		if err !=nil {
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