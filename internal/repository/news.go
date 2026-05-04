package repository

import (
	"database/sql"
	"fmt"
	"notification-bot/internal/models"

	_ "github.com/lib/pq"
)

type NewsRepository struct {
	db *sql.DB
}

func NewNewsRepository(db *sql.DB) (*NewsRepository, error) {
	return &NewsRepository{db: db}, nil
}

func (r *NewsRepository) Add(news *models.News) error {
	_, err := r.db.Exec("INSERT INTO news (source, title, text, link, date)"+
		"VALUES ($1, $2, $3, $4, $5)",
		news.Source,
		news.Title,
		news.Text,
		news.Link,
		news.Date,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *NewsRepository) IsNewsSent(link string) (bool, error) {
	var news models.News
	query := "SELECT link FROM news WHERE link = $1;"
	row := r.db.QueryRow(query, link)

	err := row.Scan(&news.Link)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("err querying data from DB: %v", err)
	}

	return true, nil
}