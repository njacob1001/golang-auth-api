package routeset

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"rumm-api/internal/core/service"
	"rumm-api/internal/platform/server/handler/registration"
)

func Client(accountService service.AccountService, validate *validator.Validate) func(r chi.Router) {
	return func (r chi.Router) {
		r.Post("/", registration.CreateClient(accountService, validate))
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", registration.FindClient(accountService))
			r.Delete("/", registration.DeleteClient(accountService))
			r.Put("/", registration.UpdateClient(accountService))
		})
	}
}