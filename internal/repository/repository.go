package repository

import (
	"database/sql"
	"fmt"
	//"os"
	_ "github.com/lib/pq"
)

type DBConfig interface {
	GetHost() string
	GetPort() int
	GetUser() string
	GetPassword() string
	GetName() string
}

func NewDB(cfg DBConfig) (*sql.DB, error) {
	//sqlInfo := os.Getenv("DATABASE_URL")

	sqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.GetHost(), cfg.GetPort(), cfg.GetUser(), cfg.GetPassword(), cfg.GetName())

	db, err := sql.Open("postgres", sqlInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	return db, nil
}
