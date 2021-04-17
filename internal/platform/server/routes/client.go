package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	service "rumm-api/internal/core/services/clients"
	"rumm-api/internal/platform/server/handler/accounthandler"
)

func Client(accountService service.AccountService, validate *validator.Validate) func(r chi.Router) {
	return func (r chi.Router) {
		r.Post("/", accounthandler.CreateClient(accountService, validate))
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", accounthandler.FindClient(accountService))
			r.Delete("/", accounthandler.DeleteClient(accountService))
			r.Put("/", accounthandler.UpdateClient(accountService))
		})
	}
}