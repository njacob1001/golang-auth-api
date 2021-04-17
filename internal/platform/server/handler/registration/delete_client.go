package registration

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"rumm-api/internal/core/service"
	"rumm-api/kit/identifier"
)


func DeleteClient(clientService service.AccountService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		requestID := chi.URLParam(r, "id")
		if requestID=="" {
			http.Error(w, "Client ID is required", http.StatusBadRequest)
		}

		id, err := identifier.ValidateIdentifier(requestID)
		if err != nil {
			http.Error(w, err.Error(),http.StatusNotAcceptable)
			return
		}

		err = clientService.DeleteClientByID(ctx, id.String)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	}
}