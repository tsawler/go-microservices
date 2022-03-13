package main

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"log"
	"net/http"
	"os"
	"time"
)

var conn *sql.DB
var counts int64

func main() {
	// connect to postgres
	dsn := "host=postgres port=5432 user=postgres password=password dbname=users sslmode=disable timezone=UTC connect_timeout=5"
	//dsn := "host=postgres user=postgres password=password dbname=users"

	for {
		connection, err := openDB(dsn)
		if err != nil {
			fmt.Println("Postgres not ready...")
			counts++
		} else {
			conn = connection
			break
		}

		if counts > 5 {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println("Backing off for one second...")
		time.Sleep(1 * time.Second)
		continue
	}

	http.HandleFunc("/authenticate", func(w http.ResponseWriter, r *http.Request) {
		var requestPayload struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		err := readJSON(w, r, &requestPayload)
		if err != nil {
			_ = errorJSON(w, err, http.StatusBadRequest)
		}

		// TODO validate against database
		payload := jsonResponse{
			Error:   false,
			Message: fmt.Sprintf("Authenticated user %s", requestPayload.Email),
		}

		_ = writeJSON(w, http.StatusAccepted, payload)
	})

	fmt.Println("Starting authentication end service on port 80")
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		log.Panic(err)
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
