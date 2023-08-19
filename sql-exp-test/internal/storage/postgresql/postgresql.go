package postgresql

import (
	"context"
	"sql-exp-test/internal/lib/e"
	"sql-exp-test/internal/storage"

	_ "database/sql"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Storage struct {
	db *sqlx.DB
}

// New creates new PostgreSQL storage.
func New(path string) (storage.Storage, error) {
	const operation = "storage.postgresql.New"

	db, err := sqlx.Open("postgres", path)
	if err != nil {
		return nil, e.Wrap(operation, err)
	}

	if err := db.Ping(); err != nil {
		return nil, e.Wrap(operation, err)
	}

	return &Storage{db: db}, nil
}

// Close PostgreSQL storage.
func (s *Storage) Close() error {
	const operation = "storage.postgresql.Close"

	return e.WrapIfErr(operation, s.db.Close())
}

// Create table if not exists
func (s *Storage) Init(ctx context.Context) error {
	const operation = "storage.postgresql.Init"

	query := `
	CREATE TABLE IF NOT EXISTS entities(
		id INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
		name TEXT NOT NULL,
		value REAL,
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

func (s *Storage) getId(ctx context.Context, entity *storage.Entities) (any, error) {
	const operation = "storage.postgresql.getId"
	var id any

	query := "SELECT id FROM entities WHERE name = $1 AND value = $2 AND description = $3 AND flag = $4;"
	stmt, err := s.db.Prepare(query)
	if err != nil {
		return 0, e.Wrap(operation, err)
	}
	err = stmt.QueryRowContext(ctx, entity.Name, entity.Value, entity.Description, entity.Flag).Scan(&id)
	if err != nil {
		return 0, e.Wrap(operation, err)
	}
	stmt.Close()

	return id, nil
}

// Create entity to storage.
func (s *Storage) Create(ctx context.Context, entity *storage.Entities) (any, error) {
	const operation = "storage.postgresql.Create"

	query := "INSERT INTO entities (name, value, description, flag) VALUES ($1, $2, $3, $4);"
	stmt, err := s.db.Prepare(query)
	if err != nil {
		return 0, e.Wrap(operation, err)
	}

	_, err = stmt.ExecContext(ctx, entity.Name, entity.Value, entity.Description, entity.Flag)
	if err != nil {
		return 0, e.Wrap(operation, err)
	}

	id, err := s.getId(ctx, entity)
	if err != nil {
		return 0, e.Wrap(operation, err)
	}

	stmt.Close()

	return id, nil
}

// Read entity from storage
func (s *Storage) Read(ctx context.Context, id any) (*storage.Entities, error) {
	const operation = "storage.postgresql.Read"

	var entity storage.Entities

	query := "SELECT * FROM entities WHERE id = $1 LIMIT 1;"
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
	const operation = "storage.postgresql.Update"

	query := "UPDATE entities SET name = $1, value = $2, description = $3, flag = $4 WHERE id = $5;"
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
	const operation = "storage.postgresql.Delete"

	query := "DELETE FROM entities WHERE name = $1 AND value = $2 AND description = $3 AND flag = $4;"
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
func (s *Storage) DeleteId(ctx context.Context, id any) error {
	const operation = "storage.postgresql.DeleteId"

	query := "DELETE FROM entities WHERE id = $1;"
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
	const operation = "storage.postgresql.IsExists"

	var count int

	query := "SELECT COUNT(*) FROM entities WHERE name = $1 AND value = $2 AND description = $3 AND flag = $4;"
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
func (s *Storage) IsExistsById(ctx context.Context, id any) (bool, error) {
	const operation = "storage.postgresql.IsExistsById"

	var count int

	query := "SELECT COUNT(*) FROM entities WHERE id = $1;"
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
func (s *Storage) LotsOfRecords(ctx context.Context, entitis ...*storage.Entities) ([]any, error) {
	const operation = "storage.postgresql.LotsOfRecords"

	var ids []any = []any{}

	query := "INSERT INTO entities (name, value, description, flag) VALUES ($1, $2, $3, $4);"

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
		_, err := stmt.ExecContext(ctx, entity.Name, entity.Value, entity.Description, entity.Flag)
		if err != nil {
			_ = tx.Rollback()
			return nil, e.Wrap(operation, err)
		}

		var id int64

		query := "SELECT id FROM entities WHERE name = $1 AND value = $2 AND description = $3 AND flag = $4;"
		stmt2, err := tx.Prepare(query)
		if err != nil {
			_ = tx.Rollback()
			return nil, e.Wrap(operation, err)
		}
		err = stmt2.QueryRowContext(ctx, entity.Name, entity.Value, entity.Description, entity.Flag).Scan(&id)
		if err != nil {
			_ = tx.Rollback()
			return nil, e.Wrap(operation, err)
		}
		stmt2.Close()

		ids = append(ids, id)
	}

	stmt.Close()

	if err := tx.Commit(); err != nil {
		return nil, e.Wrap(operation, err)
	}

	return ids, nil
}
