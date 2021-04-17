package registration

import (
	"net/http"
	"rumm-api/internal/core/services/clients"
	"rumm-api/kit/security"
)

func Logout(accountService service.AccountService, secret string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		td, err := security.ExtractTokenMetadata(secret, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := accountService.Logout(ctx, td.AccessUuid); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
