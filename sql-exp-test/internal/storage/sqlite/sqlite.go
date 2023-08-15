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
func New(path string) (storage.Storage, error) {
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

// Close SQLite storage.
func (s *Storage) Close() error {
	const operation = "storage.sqlite.Close"

	return e.WrapIfErr(operation, s.db.Close())
}

// Create table if not exists
func (s *Storage) Init(ctx context.Context) error {
	const operation = "storage.sqlite.Init"

	query := `
	CREATE TABLE IF NOT EXISTS entities(
		id INTEGER PRIMARY KEY,
		name TEXT NOT NULL,
		value DOUBLE,
		description TEXT,
		flag BOOL);
	`

	stmt, err := s.db.Prepare(query)

	if err != nil {
		return e.Wrap(operation, err)
	}

	if _, err = stmt.ExecContext(ctx); err != nil {
		return e.Wrap(operation, err)
	}
	stmt.Close()

	return nil
}

// Create entity to storage.
func (s *Storage) Create(ctx context.Context, entity *storage.Entities) (int64, error) {
	const operation = "storage.sqlite.Create"

	query := "INSERT INTO entities (name, value, description, flag) VALUES (?, ?, ?, ?);"
	stmt, err := s.db.Prepare(query)
	if err != nil {
		return 0, e.Wrap(operation, err)
	}

	result, err := stmt.ExecContext(ctx, entity.Name, entity.Value, entity.Description, entity.Flag)
	if err != nil {
		return 0, e.Wrap(operation, err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, e.Wrap(operation, err)
	}
	stmt.Close()

	return id, nil
}

// Read entity from storage
func (s *Storage) Read(ctx context.Context, id int64) (*storage.Entities, error) {
	const operation = "storage.sqlite.Read"

	var entity storage.Entities

	query := "SELECT * FROM entities WHERE id = ? LIMIT 1;"
	stmt, err := s.db.Prepare(query)
	if err != nil {
		return nil, e.Wrap(operation, err)
	}

	err = stmt.QueryRowContext(ctx, id).Scan(&entity.Id, &entity.Name, &entity.Value, &entity.Description, &entity.Flag)
	if err != nil {
		return nil, e.Wrap(operation, err)
	}
	stmt.Close()

	return &entity, nil
}

// Update entity to storage
func (s *Storage) Update(ctx context.Context, entity *storage.Entities) error {
	const operation = "storage.sqlite.Update"

	query := "UPDATE entities SET name = ?, value = ?, description = ?, flag = ? WHERE id = ?;"
	stmt, err := s.db.Prepare(query)
	if err != nil {
		return e.Wrap(operation, err)
	}
	_, err = stmt.ExecContext(ctx, entity.Name, entity.Value, entity.Description, entity.Flag, entity.Id)
	if err != nil {
		return e.Wrap(operation, err)
	}
	stmt.Close()

	return nil
}

// Remove entity from storage by entity.
func (s *Storage) Delete(ctx context.Context, entity *storage.Entities) error {
	const operation = "storage.sqlite.Delete"

	query := "DELETE FROM entities WHERE name = ? AND value = ? AND description = ? AND flag = ?;"
	stmt, err := s.db.Prepare(query)
	if err != nil {
		return e.Wrap(operation, err)
	}
	_, err = stmt.ExecContext(ctx, entity.Name, entity.Value, entity.Description, entity.Flag)
	if err != nil {
		return e.Wrap(operation, err)
	}
	stmt.Close()

	return nil
}

// Remove entity from storage by Id
func (s *Storage) DeleteId(ctx context.Context, id int64) error {
	const operation = "storage.sqlite.DeleteId"

	query := "DELETE FROM entities WHERE id = ?;"
	stmt, err := s.db.Prepare(query)
	if err != nil {
		return e.Wrap(operation, err)
	}
	_, err = stmt.ExecContext(ctx, id)
	if err != nil {
		return e.Wrap(operation, err)
	}
	stmt.Close()

	return nil
}

// IsExists checks if entity exists in storage.
func (s *Storage) IsExists(ctx context.Context, entity *storage.Entities) (bool, error) {
	const operation = "storage.sqlite.IsExists"

	var count int

	query := "SELECT COUNT(*) FROM entities WHERE name = ? AND value = ? AND description = ? AND flag = ?;"
	stmt, err := s.db.Prepare(query)
	if err != nil {
		return false, e.Wrap(operation, err)
	}
	err = stmt.QueryRowContext(ctx, entity.Name, entity.Value, entity.Description, entity.Flag).Scan(&count)
	if err != nil {
		return false, e.Wrap(operation, err)
	}
	stmt.Close()

	return count > 0, nil
}

// IsExistsById checks if entity exists in storage by Id.
func (s *Storage) IsExistsById(ctx context.Context, id int64) (bool, error) {
	const operation = "storage.sqlite.IsExistsById"

	var count int

	query := "SELECT COUNT(*) FROM entities WHERE id = ?;"
	stmt, err := s.db.Prepare(query)
	if err != nil {
		return false, e.Wrap(operation, err)
	}
	err = stmt.QueryRowContext(ctx, id).Scan(&count)
	if err != nil {
		return false, e.Wrap(operation, err)
	}
	stmt.Close()

	return count > 0, nil
}

// Lots of records.
func (s *Storage) LotsOfRecords(ctx context.Context, entitis ...*storage.Entities) ([]int64, error) {
	const operation = "storage.sqlite.LotsOfRecords"

	var ids []int64 = []int64{}

	query := "INSERT INTO entities (name, value, description, flag) VALUES (?, ?, ?, ?);"

	tx, err := s.db.Begin()
	if err != nil {
		return nil, e.Wrap(operation, err)
	}

	stmt, err := tx.Prepare(query)
	if err != nil {
		_ = tx.Rollback()
		return nil, e.Wrap(operation, err)
	}

	for _, entity := range entitis {
		result, err := stmt.ExecContext(ctx, entity.Name, entity.Value, entity.Description, entity.Flag)
		if err != nil {
			_ = tx.Rollback()
			return nil, e.Wrap(operation, err)
		}

		id, err := result.LastInsertId()
		if err != nil {
			_ = tx.Rollback()
			return nil, e.Wrap(operation, err)
		}

		ids = append(ids, id)
	}

	stmt.Close()

	if err := tx.Commit(); err != nil {
		return nil, e.Wrap(operation, err)
	}

	return ids, nil
}
