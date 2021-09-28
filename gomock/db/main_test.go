package db_test

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"

	"github.com/croixxant/golang-examples/gomock/db"
)

const (
	dbDriver = "mysql"
)

var (
	dbSource    = "root:example@tcp(127.0.0.1:3306)/gomock"
	testQueries *db.Queries
	testDB      *sql.DB
)

func TestMain(m *testing.M) {
	var err error
	if env := os.Getenv("DB_SOURCE"); env != "" {
		dbSource = env
	}

	testDB, err = sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal(err)
	}
	testQueries = db.New(testDB)
	os.Exit(m.Run())
}
