package accounthandler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	service "rumm-api/internal/core/services/clients"
	"rumm-api/kit/identifier"
)

type validateAccountRequest struct {
	Password    string `json:"password" binding:"required"`
	Identifier  string `json:"identifier" binding:"required"`
}

func ValidateAccountHandler(accountService service.AccountService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req validateAccountRequest
		if err := ctx.BindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, err.Error())
			return
		}
		_, err := accountService.Authenticate(ctx,req.Identifier, req.Password)

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

		ctx.Status(http.StatusOK)

	}
}
