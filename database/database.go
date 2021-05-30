package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Podcast struct {
	title         string
	emailId       string
	pathToRawFile string
	pathToMP3File string
	isConverted   bool
	links         []string
}

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "password"
	dbname   = "postgres"
)

func OpenConnection() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	return db
}
