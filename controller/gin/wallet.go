package gin

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"

	"github.com/croixxant/go-sandbox/usecase"
	"github.com/croixxant/go-sandbox/usecase/repo"
	"github.com/croixxant/go-sandbox/util/token"
)

type createWalletRequest struct {
	Currency string `json:"currency" binding:"required,currency"`
}

func (c *Controller) CreateWallet(ctx *gin.Context) {
	var req createWalletRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	arg := repo.CreateWalletParams{
		UserID:   authPayload.UserID,
		Currency: req.Currency,
		Balance:  0,
	}
	w, err := usecase.NewUsecase(c.repo).CreateWallet(ctx, arg)
	if err != nil {
		if dbErr, ok := err.(*mysql.MySQLError); ok {
			switch dbErr.Number {
			case 1452, 1062: // 1452: foreign key constraint fails, Duplicate entry
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, w)
}

type getWalletRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (c *Controller) GetWallet(ctx *gin.Context) {
	var req getWalletRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	args := usecase.GetWalletParams{
		RequestUserID: authPayload.UserID,
		ID:            req.ID,
	}
	w, err := usecase.NewUsecase(c.repo).GetWallet(ctx, args)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, w)
}

type listWalletRequest struct {
	Offset int32 `form:"offset"`
	Limit  int32 `form:"limit" binding:"required,min=5,max=10"`
}

func (c *Controller) ListWallets(ctx *gin.Context) {
	var req listWalletRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	arg := repo.ListwalletsParams{
		UserID: authPayload.UserID,
		Limit:  req.Limit,
		Offset: req.Offset,
	}

	ws, err := usecase.NewUsecase(c.repo).ListWallets(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, ws)
}
