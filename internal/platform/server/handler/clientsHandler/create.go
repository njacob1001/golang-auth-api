package clientsHandler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"rumm-api/internal/core/services/clients"
	"rumm-api/kit/identifier"
)

type createRequest struct {
	ID        string `json:"id" binding:"required"`
	Name      string `json:"name" binding:"required"`
	LastName  string `json:"lastName" binding:"required"`
	BirthDay  string `json:"birthday" binding:"required"`
	Email     string `json:"email" binding:"required"`
	Password  string `json:"password" binding:"required"`
	City      string `json:"city" binding:"required"`
	Address   string `json:"address" binding:"required"`
	Cellphone string `json:"cellphone" binding:"required"`
}

func CreateHandler(creatingClientService clients.ClientService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req createRequest
		if err := ctx.BindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, err.Error())
			return
		}

		err := creatingClientService.CreateClient(ctx, req.ID, req.Name, req.LastName, req.BirthDay, req.Email, req.City, req.Address, req.Cellphone, req.Password)

		if err != nil {
			switch {
			case errors.Is(err, identifier.ErrCreatingClientUUID):
				ctx.JSON(http.StatusBadRequest, err.Error())
				return
			default:
				ctx.JSON(http.StatusInternalServerError, err.Error())
				return
			}
		}

		ctx.Status(http.StatusCreated)
	}
}
