package db_test

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"

	"github.com/croixxant/go-sandbox/sqlc/db"
)

const (
	dbDriver = "mysql"
	dbSource = "root:example@tcp(127.0.0.1:3306)/sqlc"
)

var (
	testQueries *db.Queries
	testDB      *sql.DB
)

func TestMain(m *testing.M) {
	var err error

	testDB, err = sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal(err)
	}
	testQueries = db.New(testDB)
	os.Exit(m.Run())
}
