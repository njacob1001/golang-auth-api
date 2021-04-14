package accounthandler

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	service "rumm-api/internal/core/services/clients"
	"rumm-api/kit/identifier"
)

func UpdateHandler(accountService service.AccountService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req createRequest
		if err := ctx.BindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, err.Error())
			return
		}
		var paramReq clientRequest
		if err := ctx.ShouldBindUri(&paramReq); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": fmt.Sprintf("%v", err)})
			return
		}
		err := accountService.UpdateClientByID(ctx, paramReq.ID, req.Name, req.LastName, req.BirthDay.Format("2006-01-02"), req.Email, req.City, req.Address, req.Cellphone)

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
