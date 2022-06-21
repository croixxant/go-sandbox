package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"

	"github.com/croixxant/go-sandbox/gomock/api"
	"github.com/croixxant/go-sandbox/gomock/db"
)

var (
	serverAddress = "0.0.0.0:8080"
	dbSource      = "root:example@tcp(127.0.0.1:3306)/gomock"
)

func run() error {
	if env := os.Getenv("DB_SOURCE"); env != "" {
		dbSource = env
	}
	conn, err := sql.Open("mysql", dbSource)
	if err != nil {
		return err
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(serverAddress)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
