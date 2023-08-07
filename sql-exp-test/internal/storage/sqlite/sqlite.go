package sqlite

import (
	"context"
	"database/sql"
	"sql-exp-test/internal/lib/e"
	"sql-exp-test/internal/storage"

	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

// New creates new SQLite storage.
func New(path string) (*Storage, error) {
	const operation = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, e.Wrap(operation, err)
	}

	if err := db.Ping(); err != nil {
		return nil, e.Wrap(operation, err)
	}

	return &Storage{db: db}, nil
}

// Create table if not exists
func (s *Storage) Init(ctx context.Context) error {
	const operation = "storage.sqlite.Init"

	qry, err := s.db.Prepare(`
	CREATE TABLE IF NOT EXISTS pages(
		id INTEGER PRIMARY KEY,
		url TEXT NOT NULL,
		username TEXT NOT NULL);
	`)

	if err != nil {
		return e.Wrap(operation, err)
	}

	if _, err = qry.Exec(); err != nil {
		return e.Wrap(operation, err)
	}

	return nil
}

// Save saves page to storage.
func (s *Storage) Save(ctx context.Context, p *storage.Page) error {
	const operation = "storage.sqlite.Save"

	q := `INSERT INTO pages (url, username) VALUES (?, ?)`
	_, err := s.db.ExecContext(ctx, q, p.URL, p.UserName)
	if err != nil {
		return e.Wrap(operation, err)
	}

	return nil
}

// Remove page from storage.
func (s *Storage) Remove(ctx context.Context, p *storage.Page) error {
	const operation = "storage.sqlite.Remove"

	q := `DELETE FROM pages WHERE url = ? AND username = ?`
	_, err := s.db.ExecContext(ctx, q, p.URL, p.UserName)
	if err != nil {
		return e.Wrap(operation, err)
	}
	return nil
}

// IsExists checks if page exists in storage.
func (s *Storage) IsExists(ctx context.Context, p *storage.Page) (bool, error) {
	const operation = "storage.sqlite.IsExists"

	var count int

	q := `SELECT COUNT(*) FROM pages WHERE url = ? AND username = ?`
	err := s.db.QueryRowContext(ctx, q, p.URL, p.UserName).Scan(&count)
	if err != nil {
		return false, e.Wrap(operation, err)
	}

	return count > 0, nil
}
