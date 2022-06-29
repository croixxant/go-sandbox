package gapi

import (
	"fmt"

	"github.com/croixxant/go-sandbox/simplebank/db"
	"github.com/croixxant/go-sandbox/simplebank/pb"
	"github.com/croixxant/go-sandbox/simplebank/token"
	"github.com/croixxant/go-sandbox/simplebank/util"
)

// Server serves gRPC requests for our banking service.
type Server struct {
	pb.UnimplementedSimpleBankServer
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
}

// NewServer creates a new gRPC server.
func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	return server, nil
}
