package main

import (
	"database/sql"
	"log"
	"net"
	"net/http"

	"github.com/croixxant/go-sandbox/config"
	ctrl "github.com/croixxant/go-sandbox/controller/grpc"
	"github.com/croixxant/go-sandbox/repo/sqlc"
	"github.com/croixxant/go-sandbox/usecase/repo"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}
	conn, err := sql.Open("mysql", cfg.DBSource)
	if err != nil {
		log.Fatal(err)
	}
	r := sqlc.NewRepository(conn)

	go runGateway(cfg.GRPCServerAddress, cfg.HTTPServerAddress)
	runServer(cfg, r)
}

func runServer(cfg config.Config, r repo.Repository) {
	s, err := ctrl.NewServer(cfg, r)
	if err != nil {
		log.Fatal(err)
	}

	l, err := net.Listen("tcp", cfg.GRPCServerAddress)
	if err != nil {
		log.Fatal("cannot create listener:", err)
	}
	log.Printf("start gRPC server at %s", cfg.GRPCServerAddress)
	log.Fatal(s.Serve(l))
}

func runGateway(grpcServerAddress, httpServerAddress string) {
	gw, err := ctrl.NewGatewayServer(grpcServerAddress)
	if err != nil {
		log.Fatal(err)
	}

	l, err := net.Listen("tcp", httpServerAddress)
	if err != nil {
		log.Fatal("cannot create listener:", err)
	}
	log.Printf("start HTTP gateway server at %s", httpServerAddress)
	log.Fatal(http.Serve(l, gw))
}
