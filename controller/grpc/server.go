package grpc

import (
	"context"
	"fmt"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/croixxant/go-sandbox/config"
	"github.com/croixxant/go-sandbox/controller/grpc/internal"
	"github.com/croixxant/go-sandbox/usecase/repo"
	"github.com/croixxant/go-sandbox/util/token"
)

// NewServer creates a new gRPC server.
func NewServer(cfg config.Config, store repo.Repository) (*grpc.Server, error) {
	tm, err := token.NewPasetoMaker(cfg.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	c := NewController(cfg, store, tm)
	s := grpc.NewServer()
	internal.RegisterSandboxServer(s, c)
	reflection.Register(s)

	return s, nil
}

// NewServer creates a new gRPC gateway server.
func NewGatewayServer(grpcServerAddress string) (*runtime.ServeMux, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	jsonOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})

	mux := runtime.NewServeMux(jsonOption)
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err := internal.RegisterSandboxHandlerFromEndpoint(ctx, mux, grpcServerAddress, opts)
	if err != nil {
		return nil, fmt.Errorf("cannot register handler server: %w", err)
	}

	return mux, nil
}

type Controller struct {
	internal.UnimplementedSandboxServer
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
