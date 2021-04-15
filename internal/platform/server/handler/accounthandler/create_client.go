package accounthandler

import (
	"encoding/json"
	"errors"
	"net/http"
	"rumm-api/internal/core/services/clients"
	"rumm-api/kit/identifier"
)

type createRequest struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	LastName  string `json:"lastName"`
	BirthDay  string `json:"birthday"`
	Email     string `json:"email"`
	City      string `json:"city"`
	Address   string `json:"address"`
	Cellphone string `json:"cellphone"`
}

func CreateHandler(accountService accountservice.AccountService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var req createRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err := accountService.CreateClient(ctx, req.ID, req.Name, req.LastName, req.BirthDay, req.Email, req.City, req.Address, req.Cellphone)

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
		return

	}
}
