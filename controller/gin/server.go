package gin

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"

	"github.com/croixxant/go-sandbox/config"
	"github.com/croixxant/go-sandbox/usecase/repo"
	"github.com/croixxant/go-sandbox/util/token"
)

func NewServer(cfg config.Config, r repo.Repository) (*gin.Engine, error) {
	tm, err := token.NewPasetoMaker(cfg.TokenSymmetricKey)
	if err != nil {
		return nil, err
	}

	c := NewController(cfg, r, tm)
	s := gin.Default()
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("currency", validCurrency)
	}

	s.POST("/users", c.CreateUser)
	s.POST("/users/login", c.LoginUser)
	s.POST("/tokens/renew_access", c.RenewAccessToken)

	authRoutes := s.Group("/").Use(authMiddleware(c.tokenmaker))
	authRoutes.POST("/wallets", c.CreateWallet)
	authRoutes.GET("/wallets/:id", c.GetWallet)
	authRoutes.GET("/wallets", c.ListWallets)
	authRoutes.POST("/transfers", c.CreateTransfer)

	return s, nil
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

// Controller 外界とusecaseのIFを変換するレイヤーです。
// メソッドの数が多くなってきたら役割ごとにControllerを分割します。
type Controller struct {
	config     config.Config
	repo       repo.Repository
	tokenmaker token.Maker
}

func NewController(cfg config.Config, r repo.Repository, tm token.Maker) *Controller {
	return &Controller{
		config:     cfg,
		repo:       r,
		tokenmaker: tm,
	}
}
