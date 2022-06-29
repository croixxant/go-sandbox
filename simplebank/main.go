package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"net/http"

	"github.com/croixxant/go-sandbox/simplebank/api"
	"github.com/croixxant/go-sandbox/simplebank/db"
	"github.com/croixxant/go-sandbox/simplebank/gapi"
	"github.com/croixxant/go-sandbox/simplebank/pb"
	"github.com/croixxant/go-sandbox/simplebank/util"
	_ "github.com/go-sql-driver/mysql"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	// runDBMigration(config.MigrationURL, config.DBSource)

	store := db.NewStore(conn)
	// runGinServer(config, store)
	go runGatewayServer(config, store)
	runGrpcServer(config, store)
}

// func runDBMigration(migrationURL string, dbSource string) {
// 	migration, err := migrate.New(migrationURL, dbSource)
// 	if err != nil {
// 		log.Fatal("cannot create new migrate instance:", err)
// 	}

// 	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
// 		log.Fatal("failed to run migrate up:", err)
// 	}

// 	log.Println("db migrated successfully")
// }

func runGrpcServer(config util.Config, store db.Store) {
	// simplebankの実装サーバー
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	// gRPCサーバー
	grpcServer := grpc.NewServer()
	pb.RegisterSimpleBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal("cannot create listener:", err)
	}

	log.Printf("start gRPC server at %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("cannot start gRPC server:", err)
	}
}

func runGatewayServer(config util.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	jsonOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})

	grpcMux := runtime.NewServeMux(jsonOption)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = pb.RegisterSimpleBankHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatal("cannot register handler server:", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	// statikFS, err := fs.New()
	// if err != nil {
	// 	log.Fatal("cannot create statik fs:", err)
	// }

	// swaggerHandler := http.StripPrefix("/swagger/", http.FileServer(statikFS))
	swaggerHandler := http.StripPrefix("/swagger/", http.FileServer(http.Dir("./doc/swagger")))
	mux.Handle("/swagger/", swaggerHandler)

	listener, err := net.Listen("tcp", config.HTTPServerAddress)
	if err != nil {
		log.Fatal("cannot create listener:", err)
	}

	log.Printf("start HTTP gateway server at %s", listener.Addr().String())
	err = http.Serve(listener, mux)
	if err != nil {
		log.Fatal("cannot start HTTP gateway server:", err)
	}
}

func runGinServer(config util.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
