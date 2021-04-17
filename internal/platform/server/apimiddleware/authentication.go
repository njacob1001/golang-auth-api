package apimiddleware

import (
	"github.com/go-redis/redis/v8"
	"net/http"
	"rumm-api/kit/security"
)

func JwtAuth(secret string, rdb *redis.Client) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			tokenAuth, err := security.ExtractTokenMetadata(secret, r)

			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			userID, err := security.FetchAuth(ctx, tokenAuth, rdb)

			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			if userID != tokenAuth.UserID {
				http.Error(w, "Corrupt", http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)

		})
	}
}
