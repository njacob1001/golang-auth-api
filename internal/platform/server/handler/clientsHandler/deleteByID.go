package clientsHandler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	service "rumm-api/internal/core/services/clients"
	"rumm-api/kit/identifier"
)


func DeleteByIDHandler(clientService service.ClientService) gin.HandlerFunc {
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

		err = clientService.DeleteClientByID(ctx, id.String)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": fmt.Sprintf("%v", err)})
			return
		}
		ctx.Status(http.StatusOK)
		return
	}
}