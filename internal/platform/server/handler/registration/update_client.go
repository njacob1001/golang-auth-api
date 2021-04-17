package registration

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"net/http"
	"rumm-api/internal/core/service"
	"rumm-api/kit/identifier"
)


type updateRequest struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	LastName  string `json:"lastName"`
	BirthDay  string `json:"birthday"`
	Email     string `json:"email"`
	City      string `json:"city"`
	Address   string `json:"address"`
	Cellphone string `json:"cellphone"`
}

func UpdateClient(accountService service.AccountService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var req updateRequest


		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		clientID := chi.URLParam(r, "id")

		err := accountService.UpdateClientByID(ctx, clientID, req.Name, req.LastName, req.BirthDay, req.Email, req.City, req.Address, req.Cellphone)

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
