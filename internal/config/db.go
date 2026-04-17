package config

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func ConnectDB() {
	connStr := "user=apple dbname=goapp sslmode=disable"

	db, err:= sql.Open("postgres", connStr)
	if err != nil{
		log.Fatal("DB connection error: ", err)
	}

	err = db.Ping()
	if err != nil{
		log.Fatal("DB not reachable: ", err)
	}

	DB = db
	log.Println("Successfully connected to DB !!")
}