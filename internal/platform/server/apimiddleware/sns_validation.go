package apimiddleware

import (
	"net/http"
	"rumm-api/internal/core/service"
	"rumm-api/kit/security"
)

func SnsValidation(s service.AccountService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			data, err := security.ExtractSnsTokenData(s.GetSnsSecret(), r)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			value, err := s.Cache.Get(ctx, data.Cellphone).Result()
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
			if value == "" {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
