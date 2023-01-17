package tests

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jackc/pgx/v5"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = "5432"
	user     = "gen"
	password = "1111"
	dbname   = "gen"
)

func TestPostgresConnect(t *testing.T) {
	psqlConnect := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", host, port, dbname, user, password)
	db, err := sql.Open("postgres", psqlConnect)
	if err != nil {
		log.Fatal("Error", err)
	}

	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	innerSTMT := `select id, first_name from person`

	rows, e := db.Query(innerSTMT)
	if e != nil {
		log.Fatal(e)
	}

	for rows.Next() {
		var id int
		var name string

		err := rows.Scan(&id, &name)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(id, name)
	}
	log.Println("Connection set")
}

func TestConPqx(t *testing.T) {
	psqlConnect := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", host, port, dbname, user, password)

	conn, err := pgx.Connect(context.Background(), psqlConnect)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	defer conn.Close(context.Background())

	rows, _ := conn.Query(context.Background(), "select id, first_name, last_name from person")
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	type row struct {
		id        int
		firstName string
		lastName  string
	}

	rowArr := row{}

	for rows.Next() {
		err := rows.Scan(&rowArr.id, &rowArr.firstName, &rowArr.lastName)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(rowArr)
	}
}
