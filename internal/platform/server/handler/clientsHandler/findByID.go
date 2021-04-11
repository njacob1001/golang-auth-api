package clientsHandler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"rumm-api/internal/core/services/clients"
	"rumm-api/kit/identifier"
)

type clientRequest struct {
	ID string `uri:"id" binding:"required"`
}

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

func FindByIDHandler(clientService clients.ClientService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request clientRequest
		if err := ctx.ShouldBindUri(&request); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": fmt.Sprintf("%v", err)})
			return
		}

		id, err := identifier.ValidateIdentifier(request.ID)
		if err != nil {
			ctx.JSON(http.StatusNotAcceptable, gin.H{"message": fmt.Sprintf("%v", err)})
			return
		}

		client, err := clientService.FindClientByID(ctx, id.String)

		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("%v", err)})
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

		ctx.JSON(http.StatusOK, response)
		return

	}
}
