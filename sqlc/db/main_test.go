package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

const (
	dbDriver = "mysql"
	dbDSN    = "root:example@tcp(127.0.0.1:3306)/sqlc?parseTime=true"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	var err error
	testDB, err = sql.Open(dbDriver, dbDSN)
	if err != nil {
		log.Fatal("connot connect to db: %w", err)
	}
	testQueries = New(testDB)
	os.Exit(m.Run())
}
