package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"github.com/croixxant/golang-examples/gomock/db"
)

type createUserRequest struct {
	Email    string `binding:"required"`
	Password string `binding:"required"`
}

func (server *Server) createUser(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	arg := db.CreateUserParams{
		Email:          req.Email,
		HashedPassword: string(hashedPassword),
	}
	if _, err = server.store.CreateUser(ctx, arg); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, struct{}{})
}

type getUserRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getUser(ctx *gin.Context) {
	var req getUserRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUser(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, user)
}
