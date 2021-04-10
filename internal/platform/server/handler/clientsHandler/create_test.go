package clientsHandler

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/huandu/go-assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"rumm-api/internal/core/services/clients"
	"rumm-api/mocks/mockups"
	"testing"
)

func TestCreateHandler(t *testing.T) {
	clientRepository := new(storagemocks.ClientRepository)
	clientRepository.On("Save", mock.Anything, mock.Anything).Return(nil)

	createClientService := clients.NewClientService(clientRepository)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/clientsHandler", CreateHandler(createClientService))

	t.Run("given and invalid request it return 400", func(t *testing.T) {
		createClientReq := createRequest{
			ID:  "01c8488f-a19d-492d-8a00-9f22278db342",
			Email: "some",
		}

		b, err := json.Marshal(createClientReq)
		// stop execution test if the condition fails
		require.NoError(t, err) // comprobar que no exista error al marchalear el json

		req, err := http.NewRequest(http.MethodPost, "/clientsHandler", bytes.NewBuffer(b))
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
			Email:    "some",
			Address:  "test",
			City:     "test",
			BirthDay: "2020-12-12",
			LastName: "test",
			Name:     "test",
			Password: "test",
			Cellphone: "testing",
		}
		b, err := json.Marshal(createClientReq)
		// stop execution test if the condition fails
		require.NoError(t, err) // comprobar que no exista error al marchalear el json

		req, err := http.NewRequest(http.MethodPost, "/clientsHandler", bytes.NewBuffer(b))
		require.NoError(t, err) //  comprobar que no exista error al llamar el endpoint

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusCreated, res.StatusCode)
	})
}
