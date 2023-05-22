package service

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

// PostgresDBHandler contains a pointer to the SQL database.
type PostgresDBHandler struct {
	*sql.DB
}

var Db PostgresDBHandler

const (
	host        = "localhost"
	port        = "5432"
	dbName      = "labora-proyect-1"
	rolName     = "postgres"
	rolPassword = "1234"
)

// PingOrDie it is a function that sends a ping to the server and returns an error if it does not receive the information.
func (db *PostgresDBHandler) PingOrDie() {
	if err := db.Ping(); err != nil {
		log.Fatalf("Cannot reach database, error: %v", err)
	}
}

// DeleteItem it is a function that connects to the database
func Connect_DB() {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, rolName, rolPassword, dbName)
	dbConn, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Successful connection to the database:", dbConn)
	Db = PostgresDBHandler{dbConn}
	Db.PingOrDie()
	if err != nil {
		log.Fatalf("Error while database connection: %v", err)
	}
}
