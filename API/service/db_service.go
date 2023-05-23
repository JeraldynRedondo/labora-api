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

//var Db PostgresDBHandler

const (
	host        = "localhost"
	port        = "5432"
	dbName      = "labora-proyect-1"
	rolName     = "postgres"
	rolPassword = "1234"
)

// Connect_DB it is a function that connects to the database
func Connect_DB() (*PostgresDBHandler, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, rolName, rolPassword, dbName)
	dbConn, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	fmt.Println("Successful connection to the database:", dbConn)
	DbHandler := &PostgresDBHandler{dbConn}
	return DbHandler, nil
}
