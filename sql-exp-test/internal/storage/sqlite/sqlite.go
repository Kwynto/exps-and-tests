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

	query, err := s.db.Prepare(`
	CREATE TABLE IF NOT EXISTS entities(
		id INTEGER PRIMARY KEY,
		name TEXT NOT NULL,
		value DOUBLE
		description TEXT,
		flag BOOL);
	`)

	if err != nil {
		return e.Wrap(operation, err)
	}

	if _, err = query.Exec(); err != nil {
		return e.Wrap(operation, err)
	}

	return nil
}

// Create entity to storage.
func (s *Storage) Create(ctx context.Context, entity *storage.Entities) (int64, error) {
	const operation = "storage.sqlite.Create"

	query := `INSERT INTO entities (name, value, description, flag) VALUES (?, ?, ?, ?);`
	result, err := s.db.ExecContext(ctx, query, entity.Name, entity.Value, entity.Description, entity.Flag)
	if err != nil {
		return 0, e.Wrap(operation, err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, e.Wrap(operation, err)
	}
	return id, nil
}

// Read entity from storage
func (s *Storage) Read(ctx context.Context, id int64) (*storage.Entities, error) {
	const operation = "storage.sqlite.Read"

	var entity storage.Entities

	query := `SELECT * FROM entities WHERE id = ? LIMIT 1;`
	err := s.db.QueryRowContext(ctx, query, id).Scan(&entity)
	if err != nil {
		return nil, e.Wrap(operation, err)
	}

	return &entity, nil
}

// Update entity to storage
func (s *Storage) Update(ctx context.Context, entity *storage.Entities) error {
	const operation = "storage.sqlite.Update"

	query := `UPDATE entities SET (name, value, description, flag) VALUES (?, ?, ?, ?) WHERE id = ?;`
	_, err := s.db.ExecContext(ctx, query, entity.Name, entity.Value, entity.Description, entity.Flag, entity.Id)
	if err != nil {
		return e.Wrap(operation, err)
	}

	return nil
}

// Remove entity from storage by entity.
func (s *Storage) Delete(ctx context.Context, entity *storage.Entities) error {
	const operation = "storage.sqlite.Delete"

	query := `DELETE FROM entities WHERE name = ?;`
	_, err := s.db.ExecContext(ctx, query, entity.Name)
	if err != nil {
		return e.Wrap(operation, err)
	}
	return nil
}

// Remove entity from storage by Id
func (s *Storage) DeleteId(ctx context.Context, id int64) error {
	const operation = "storage.sqlite.DeleteId"

	query := `DELETE FROM entities WHERE id = ?;`
	_, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return e.Wrap(operation, err)
	}
	return nil
}

// IsExists checks if entity exists in storage.
func (s *Storage) IsExists(ctx context.Context, entity *storage.Entities) (bool, error) {
	const operation = "storage.sqlite.IsExists"

	var count int

	query := `SELECT COUNT(*) FROM entities WHERE name = ?;`
	err := s.db.QueryRowContext(ctx, query, entity.Name).Scan(&count)
	if err != nil {
		return false, e.Wrap(operation, err)
	}

	return count > 0, nil
}
