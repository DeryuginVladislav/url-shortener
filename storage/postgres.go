package storage

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

var db *pgxpool.Pool

// Инициализация подключения к БД
func InitDB(connString string) {
	var err error
	db, err = pgxpool.New(context.Background(), connString)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	log.Println("Connected to PostgreSQL!")
}

// Сохранение URL в БД
func SaveURL(shortID, longURL string) error {
	_, err := db.Exec(context.Background(),
		"INSERT INTO urls (short_id, long_url) VALUES ($1, $2)", shortID, longURL)
	return err
}

// Получение оригинального URL из БД
func GetURL(shortID string) (string, bool) {
	var longURL string
	err := db.QueryRow(context.Background(),
		"SELECT long_url FROM urls WHERE short_id = $1", shortID).Scan(&longURL)
	if err != nil {
		return "", false
	}
	return longURL, true
}
