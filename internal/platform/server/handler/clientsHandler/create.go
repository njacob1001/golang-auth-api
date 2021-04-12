package clientsHandler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"rumm-api/internal/core/services/clients"
	"rumm-api/kit/identifier"
	"time"
)

type createRequest struct {
	ID        string `json:"id" binding:"required"`
	Name      string `json:"name" binding:"required"`
	LastName  string `json:"lastName" binding:"required"`
	BirthDay  time.Time `json:"birthday" binding:"required" time_format:"2006-01-02"`
	Email     string `json:"email" binding:"required"`
	City      string `json:"city" binding:"required"`
	Address   string `json:"address" binding:"required"`
	Cellphone string `json:"cellphone" binding:"required"`
}

func CreateHandler(clientService clients.ClientService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req createRequest
		if err := ctx.BindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, err.Error())
			return
		}

		err := clientService.CreateClient(ctx, req.ID, req.Name, req.LastName, req.BirthDay.Format("2006-01-02"), req.Email, req.City, req.Address, req.Cellphone)

		if err != nil {
			switch {
			case errors.Is(err, identifier.ErrInvalidClientUUID):
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
