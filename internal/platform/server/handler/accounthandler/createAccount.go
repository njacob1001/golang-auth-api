package accounthandler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	service "rumm-api/internal/core/services/clients"
	"rumm-api/kit/identifier"
)

type createAccountRequest struct {
	ID          string `json:"id" binding:"required"`
	Password    string `json:"password" binding:"required"`
	Identifier  string `json:"identifier" binding:"required"`
	AccountType string `json:"accountType" binding:"required"`
}

func CreateAccountHandler(accountService service.AccountService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req createAccountRequest
		if err := ctx.BindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, err.Error())
			return
		}
		err := accountService.CreateAccount(ctx, req.ID, req.Identifier, req.Password, req.AccountType)

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
