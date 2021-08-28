package registration

import (
	"net/http"
	"rumm-api/internal/core/service"
	"rumm-api/kit/security"
)

func ResendCode(s service.AccountService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		data, err := security.ExtractSnsTokenData(s.GetSnsSecret(), r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := s.RegisterCode(ctx, data.Cellphone, ""); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusCreated)
		return
	}
}
