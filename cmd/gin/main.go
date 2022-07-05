package main

import (
	"database/sql"
	"log"

	"github.com/croixxant/go-sandbox/config"
	"github.com/croixxant/go-sandbox/controller/gin"
	"github.com/croixxant/go-sandbox/repo/sqlc"
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

	s, err := gin.NewServer(cfg, r)
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(s.Run(cfg.HTTPServerAddress))
}
