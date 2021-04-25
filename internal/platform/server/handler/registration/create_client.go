package registration

import (
	"encoding/json"
	"errors"
	"github.com/go-playground/validator/v10"
	"net/http"
	"rumm-api/internal/core/service"
	"rumm-api/kit/identifier"
)

type createRequest struct {
	ID        string `json:"id" validate:"required"`
	Name      string `json:"name" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
	BirthDay  string `json:"birthday" validate:"required,datetime=2006-01-02"`
	Email     string `json:"email" validate:"required,email"`
	City      string `json:"city" validate:"required"`
	Address   string `json:"address" validate:"required"`
	Cellphone string `json:"cellphone" validate:"required"`
}

func CreateClient(accountService service.AccountService, validate *validator.Validate) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()
		var req createRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := validate.Struct(req); err != nil {
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

func CreateTemporalClient(accountService service.AccountService, validate *validator.Validate) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var req createRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := validate.Struct(req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err := accountService.CreateTemporalClient(ctx, req.ID, req.Name, req.LastName, req.BirthDay, req.Email, req.City, req.Address, req.Cellphone)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.WriteHeader(http.StatusCreated)
		return
	}
}