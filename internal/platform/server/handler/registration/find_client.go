package registration

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
	"rumm-api/internal/core/service"
	"rumm-api/kit/identifier"
)

type clientResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	LastName  string `json:"last_name"`
	Birthday  string `json:"birthday"`
	Email     string `json:"email"`
	City      string `json:"city"`
	Address   string `json:"address"`
	Cellphone string `json:"cellphone"`
}

func FindClient(accountService service.AccountService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		requestID := chi.URLParam(r, "id")
		if requestID=="" {
			http.Error(w, "Client ID is required", http.StatusBadRequest)
			return
		}



		id, err := identifier.ValidateIdentifier(requestID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotAcceptable)
			return
		}

		client, err := accountService.FindClientByID(ctx, id.String)

		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		response := clientResponse{
			ID: client.ID(),
			Name: client.Name(),
			LastName: client.LastName(),
			Birthday: client.BirthDay(),
			Email: client.Email(),
			City: client.City(),
			Address: client.Address(),
			Cellphone: client.Cellphone(),
		}

		j, err := json.Marshal(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}



		if _, err := w.Write(j); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusCreated)


		return

	}
}
