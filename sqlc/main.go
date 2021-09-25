package main

import (
	"context"
	"database/sql"
	"log"
	"reflect"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"

	"github.com/croixxant/golang-examples/sqlc/db"
)

func run() error {
	ctx := context.Background()
	conn, err := sql.Open("mysql", "root:example@tcp(127.0.0.1:3306)/sqlc")
	if err != nil {
		return err
	}

	queries := db.New(conn)

	// list all users
	users, err := queries.ListUsers(ctx)
	if err != nil {
		return err
	}
	log.Println(users)

	// create an user
	hashed, err := bcrypt.GenerateFromPassword([]byte("password"), 12)
	if err != nil {
		return err
	}
	result, err := queries.CreateUser(ctx, db.CreateUserParams{
		Email:          "example@example.com",
		HashedPassword: string(hashed),
	})
	if err != nil {
		return err
	}

	insertedUserID, err := result.LastInsertId()
	if err != nil {
		return err
	}
	log.Println(insertedUserID)

	// get the user we just inserted
	fetchedUser, err := queries.GetUser(ctx, insertedUserID)
	if err != nil {
		return err
	}

	// prints true
	log.Println(reflect.DeepEqual(insertedUserID, fetchedUser.ID))
	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
