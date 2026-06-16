package config

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func ConnectDB() {
	connStr := GetDatabaseURL()

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("❌ DB connection error: ", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("❌ DB not reachable: ", err)
	}

	// Set connection pool settings for production
	if IsProduction() {
		db.SetMaxOpenConns(25)
		db.SetMaxIdleConns(5)
	}

	DB = db
	log.Println("✅ Successfully connected to DB !!")
}