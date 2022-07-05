package gin

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/croixxant/go-sandbox/usecase"
	"github.com/croixxant/go-sandbox/util/token"
)

type transferRequest struct {
	FromWalletID int64  `json:"from_wallet_id" binding:"required,min=1"`
	ToWalletID   int64  `json:"to_wallet_id" binding:"required,min=1"`
	Amount       int64  `json:"amount" binding:"required,gt=0"`
	Currency     string `json:"currency" binding:"required,currency"`
}

func (c *Controller) CreateTransfer(ctx *gin.Context) {
	var req transferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	arg := usecase.CreateTransferParams{
		RequestUserID: authPayload.UserID,
		FromWalletID:  req.FromWalletID,
		ToWalletID:    req.ToWalletID,
		Amount:        req.Amount,
		Currency:      req.Currency,
	}
	t, err := usecase.NewUsecase(c.repo).CreateTransfer(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, t)
}
