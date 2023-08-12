package mysql

import (
	"context"
	"database/sql"
	"sql-exp-test/internal/lib/e"
	"sql-exp-test/internal/storage"

	_ "github.com/go-sql-driver/mysql"
)

type Storage struct {
	db *sql.DB
}

// New creates new SQLite storage.
func New(path string) (storage.Storage, error) {
	const operation = "storage.mysql.New"

	db, err := sql.Open("mysql", path)
	if err != nil {
		return nil, e.Wrap(operation, err)
	}

	if err := db.Ping(); err != nil {
		return nil, e.Wrap(operation, err)
	}

	return &Storage{db: db}, nil
}

// New creates new SQLite storage.
func (s *Storage) Close() error {
	const operation = "storage.mysql.Close"

	return e.WrapIfErr(operation, s.db.Close())
}

// Create table if not exists
func (s *Storage) Init(ctx context.Context) error {
	const operation = "storage.mysql.Init"

	query, err := s.db.Prepare(`
	CREATE TABLE IF NOT EXISTS entities(
		id INTEGER PRIMARY KEY AUTO_INCREMENT,
		name TEXT NOT NULL,
		value DOUBLE,
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
	const operation = "storage.mysql.Create"

	query := "INSERT INTO entities (name, value, description, flag) VALUES (?, ?, ?, ?);"
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
	const operation = "storage.mysql.Read"

	var entity storage.Entities

	query := `SELECT * FROM entities WHERE id = ? LIMIT 1;`
	err := s.db.QueryRowContext(ctx, query, id).Scan(&entity.Id, &entity.Name, &entity.Value, &entity.Description, &entity.Flag)
	if err != nil {
		return nil, e.Wrap(operation, err)
	}

	return &entity, nil
}

// Update entity to storage
func (s *Storage) Update(ctx context.Context, entity *storage.Entities) error {
	const operation = "storage.mysql.Update"

	// query := `UPDATE entities SET (name, value, description, flag) VALUES (?, ?, ?, ?) WHERE id = ?;`
	query := `UPDATE entities SET name = ?, value = ?, description = ?, flag = ? WHERE id = ?;`
	_, err := s.db.ExecContext(ctx, query, entity.Name, entity.Value, entity.Description, entity.Flag, entity.Id)
	if err != nil {
		return e.Wrap(operation, err)
	}

	return nil
}

// Remove entity from storage by entity.
func (s *Storage) Delete(ctx context.Context, entity *storage.Entities) error {
	const operation = "storage.mysql.Delete"

	query := `DELETE FROM entities WHERE name = ? AND value = ? AND description = ? AND flag = ?;`
	_, err := s.db.ExecContext(ctx, query, entity.Name, entity.Value, entity.Description, entity.Flag)
	if err != nil {
		return e.Wrap(operation, err)
	}
	return nil
}

// Remove entity from storage by Id
func (s *Storage) DeleteId(ctx context.Context, id int64) error {
	const operation = "storage.mysql.DeleteId"

	query := `DELETE FROM entities WHERE id = ?;`
	_, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return e.Wrap(operation, err)
	}
	return nil
}

// IsExists checks if entity exists in storage.
func (s *Storage) IsExists(ctx context.Context, entity *storage.Entities) (bool, error) {
	const operation = "storage.mysql.IsExists"

	var count int

	query := `SELECT COUNT(*) FROM entities WHERE name = ? AND value = ? AND description = ? AND flag = ?;`
	err := s.db.QueryRowContext(ctx, query, entity.Name, entity.Value, entity.Description, entity.Flag).Scan(&count)
	if err != nil {
		return false, e.Wrap(operation, err)
	}

	return count > 0, nil
}

// IsExistsById checks if entity exists in storage by Id.
func (s *Storage) IsExistsById(ctx context.Context, id int64) (bool, error) {
	const operation = "storage.mysql.IsExistsById"

	var count int

	query := `SELECT COUNT(*) FROM entities WHERE id = ?;`
	err := s.db.QueryRowContext(ctx, query, id).Scan(&count)
	if err != nil {
		return false, e.Wrap(operation, err)
	}

	return count > 0, nil
}
