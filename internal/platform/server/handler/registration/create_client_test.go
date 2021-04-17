package registration

import (
	"bytes"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/huandu/go-assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"rumm-api/internal/core/service"
	"rumm-api/mocks/mockups"
	"testing"
)

func TestCreateHandler(t *testing.T) {
	clientRepository := new(storagemocks.ClientRepository)
	clientRepository.On("Create", mock.Anything, mock.Anything).Return(nil)
	accountRepository := new(storagemocks.AccountRepository)

	createClientService := service.NewAccountService(accountRepository, clientRepository)
	v := validator.New()
	r := chi.NewRouter()
	r.Use(middleware.GetHead)
	r.Post("/clients", CreateClient(createClientService, v))


	t.Run("given and invalid request it return 400", func(t *testing.T) {
		createClientReq := createRequest{
			ID:  "invalid-uuid",
			Email: "some",
		}

		b, err := json.Marshal(createClientReq)
		// stop execution test if the condition fails
		require.NoError(t, err) // comprobar que no exista error al marchalear el json

		req, err := http.NewRequest(http.MethodPost, "/clients", bytes.NewBuffer(b))
		require.NoError(t, err) //  comprobar que no exista error al llamar el endpoint

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	})

	t.Run("Given a valid request it returns 201", func(t *testing.T) {
		createClientReq := createRequest{
			ID:     "01c8488f-a19d-492d-8a00-9f22278db342",
			Email:    "jacob@testing.com",
			Address:  "test",
			City:     "test",
			BirthDay: "2021-02-28",
			LastName: "test",
			Name:     "test",
			Cellphone: "testing",
		}
		b, err := json.Marshal(createClientReq)
		// stop execution test if the condition fails
		require.NoError(t, err) // comprobar que no exista error al marchalear el json

		req, err := http.NewRequest(http.MethodPost, "/clients", bytes.NewBuffer(b))
		require.NoError(t, err) //  comprobar que no exista error al llamar el endpoint

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusCreated, res.StatusCode)
	})
}
