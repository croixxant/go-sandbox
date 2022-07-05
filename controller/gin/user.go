package gin

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"

	"github.com/croixxant/go-sandbox/usecase"
)

func (c *Controller) CreateUser(ctx *gin.Context) {
	var args usecase.CreateUserParams
	if err := ctx.ShouldBindJSON(&args); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	u, err := usecase.NewUserUsecase(c.repo, c.tokenmaker, c.config).CreateUser(ctx, args)
	if err != nil {
		if dbErr, ok := err.(*mysql.MySQLError); ok {
			switch dbErr.Number {
			case 1062:
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, u)
}

func (c *Controller) LoginUser(ctx *gin.Context) {
	var args usecase.LoginUserParams
	if err := ctx.ShouldBindJSON(&args); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	args.UserAgent = ctx.Request.UserAgent()
	args.ClientIP = ctx.ClientIP()

	resp, err := usecase.NewUserUsecase(c.repo, c.tokenmaker, c.config).LoginUser(ctx, args)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

func (c *Controller) RenewAccessToken(ctx *gin.Context) {
	var req usecase.RenewAccessTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	resp, err := usecase.NewUserUsecase(c.repo, c.tokenmaker, c.config).RenewAccessToken(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, resp)
}
