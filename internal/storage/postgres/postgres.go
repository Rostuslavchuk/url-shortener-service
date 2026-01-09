package postgres

import (
	"database/sql"
	"errors"
	"fmt"

	"url_shortener/internal/storage"

	"github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.postgres.New"

	db, err := sql.Open("postgres", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s %w", op, err)
	}
	return &Storage{db: db}, nil
}

func (s *Storage) SaveURL(urlToSave, alias string) (int64, error) {
	const op = "storage.postgresql.SaveURL"

	var id int64
	sqlStatment := `INSERT INTO url (url, alias) VALUES ($1, $2) RETURNING id`
	err := s.db.QueryRow(sqlStatment, urlToSave, alias).Scan(&id)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" {
				return 0, fmt.Errorf("%w", storage.ErrURLExists)
			}
		}
		return 0, fmt.Errorf("%s %w", op, err)
	}

	return id, nil
}

func (s *Storage) GetURL(alias string) (string, error) {
	const op = "storage.postgresql.GetURL"
	var resURL string

	sqlStatment := "SELECT url FROM url WHERE alias = $1"
	err := s.db.QueryRow(sqlStatment, alias).Scan(&resURL)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("%s %w", op, storage.ErrNotFound)
		}
		return "", fmt.Errorf("%s, %w", op, err)
	}

	return resURL, nil
}

func (s *Storage) DeleteURL(alias string) (int64, error) {
	const op = "storage.postgresql.DeleteURL"
	sqlStatement := "DELETE FROM url WHERE alias = $1 RETURNING id"
	var id int64

	err := s.db.QueryRow(sqlStatement, alias).Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, fmt.Errorf("%s %w", op, storage.ErrNotFound)
		}
		return 0, fmt.Errorf("%s %w", op, err)
	}
	return id, nil
}
