// скрипт для создания таблиц(ы), в прод такое нельзя, просто для удобства пока
package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	sqlInfo := "host=localhost port=5432 user=postgres password=goose dbname=News sslmode=disable"

	db, err := sql.Open("postgres", sqlInfo)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	query := `
	CREATE EXTENSION IF NOT EXISTS "pgcrypto";

	CREATE TABLE IF NOT EXISTS news (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		source TEXT,
		title TEXT,
		text TEXT,
		link TEXT UNIQUE,
		date TEXT,
		created_at TIMESTAMP DEFAULT NOW()
	);
	`

	_, err = db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("news table ensured")
}