package clickhouse

import (
	"context"
	"fmt"
	"log"
	"sql-exp-test/internal/lib/e"
	"sql-exp-test/internal/storage"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/google/uuid"
)

type Storage struct {
	db driver.Conn
}

// New creates new ClickHouse storage.
func New(host, port, database, username, password string) (storage.Storage, error) {
	const operation = "storage.clickhouse.New"

	hostport := fmt.Sprintf("%s:%s", host, port)

	db, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{hostport},
		Auth: clickhouse.Auth{
			Database: database,
			Username: username,
			Password: password,
		},
		Debug:           true,
		DialTimeout:     time.Second * 10,
		MaxOpenConns:    10,
		MaxIdleConns:    5,
		ConnMaxLifetime: time.Hour,
	})

	if err != nil {
		return nil, e.Wrap(operation, err)
	}

	if err := db.Ping(context.Background()); err != nil {
		return nil, e.Wrap(operation, err)
	}

	return &Storage{db: db}, nil
}

// Close ClickHouse storage.
func (s *Storage) Close() error {
	const operation = "storage.clickhouse.Close"

	return e.WrapIfErr(operation, s.db.Close())
}

// Create table if not exists
func (s *Storage) Init(ctx context.Context) error {
	const operation = "storage.clickhouse.Init"

	query := `
	CREATE TABLE entities
	(
		id UUID,
		name Nullable (String),
		value Nullable (Float64),
		description Nullable (String),
		flag BOOL default 0
	)
		engine = MergeTree 
		primary key (id);
	`

	err := s.db.Exec(ctx, query)
	if err != nil {
		return e.Wrap(operation, err)
	}

	return nil
}

func (s *Storage) getId(ctx context.Context, entity *storage.Entities) (string, error) {
	const operation = "storage.clickhouse.getId"
	var id string

	query := "SELECT id FROM entities WHERE name = ? AND value = ? AND description = ? AND flag = ?;"

	err := s.db.QueryRow(ctx, query, entity.Name, entity.Value, entity.Description, entity.Flag).Scan(&id)
	if err != nil {
		return "", e.Wrap(operation, err)
	}

	return id, nil
}

// Create entity to storage.
func (s *Storage) Create(ctx context.Context, entity *storage.Entities) (any, error) {
	const operation = "storage.clickhouse.Create"

	query := "INSERT INTO entities VALUES (generateUUIDv4(), ?, ?, ?, ?);"

	err := s.db.Exec(ctx, query, entity.Name, entity.Value, entity.Description, entity.Flag)
	if err != nil {
		return 0, e.Wrap(operation, err)
	}
	id, err := s.getId(ctx, entity)
	if err != nil {
		return 0, e.Wrap(operation, err)
	}

	return id, nil
}

// Read entity from storage
func (s *Storage) Read(ctx context.Context, id any) (*storage.Entities, error) {
	const operation = "storage.clickhouse.Read"

	var entity storage.Entities
	var tId string

	query := "SELECT * FROM entities WHERE id = ? LIMIT 1;"

	err := s.db.QueryRow(ctx, query, id.(string)).Scan(&tId, &entity.Name, &entity.Value, &entity.Description, &entity.Flag)
	if err != nil {
		return nil, e.Wrap(operation, err)
	}
	entity.Id = tId

	return &entity, nil
}

// Update entity to storage
func (s *Storage) Update(ctx context.Context, entity *storage.Entities) error {
	// const operation = "storage.clickhouse.Update"

	// ClickHouse don't support update data.
	log.Println("ClickHouse don't support update data.")

	return nil
}

// Remove entity from storage by entity.
func (s *Storage) Delete(ctx context.Context, entity *storage.Entities) error {
	const operation = "storage.clickhouse.Delete"

	query := "DELETE FROM entities WHERE name = ? AND value = ? AND description = ? AND flag = ?;"
	err := s.db.Exec(ctx, query, entity.Name, entity.Value, entity.Description, entity.Flag)
	if err != nil {
		return e.Wrap(operation, err)
	}

	return nil
}

// Remove entity from storage by Id
func (s *Storage) DeleteId(ctx context.Context, id any) error {
	const operation = "storage.clickhouse.DeleteId"

	query := "DELETE FROM entities WHERE id = ?;"
	err := s.db.Exec(ctx, query, id)
	if err != nil {
		return e.Wrap(operation, err)
	}

	return nil
}

// IsExists checks if entity exists in storage.
func (s *Storage) IsExists(ctx context.Context, entity *storage.Entities) (bool, error) {
	const operation = "storage.clickhouse.IsExists"

	var count uint64

	query := "SELECT COUNT(*) FROM entities WHERE name = ? AND value = ? AND description = ? AND flag = ?;"
	err := s.db.QueryRow(ctx, query, entity.Name, entity.Value, entity.Description, entity.Flag).Scan(&count)
	if err != nil {
		return false, e.Wrap(operation, err)
	}

	return count > 0, nil
}

// IsExistsById checks if entity exists in storage by Id.
func (s *Storage) IsExistsById(ctx context.Context, id any) (bool, error) {
	const operation = "storage.clickhouse.IsExistsById"

	var count uint64

	query := "SELECT COUNT(*) FROM entities WHERE id = ?;"
	err := s.db.QueryRow(ctx, query, id).Scan(&count)
	if err != nil {
		return false, e.Wrap(operation, err)
	}

	return count > 0, nil
}

// Lots of records.
func (s *Storage) LotsOfRecords(ctx context.Context, entitis ...*storage.Entities) ([]any, error) {
	const operation = "storage.clickhouse.LotsOfRecords"

	var ids []any = []any{}

	query := "INSERT INTO entities VALUES (?, ?, ?, ?, ?);"

	batch, err := s.db.PrepareBatch(ctx, query)
	if err != nil {
		return nil, e.Wrap(operation, err)
	}

	for _, entity := range entitis {

		err := batch.Append(uuid.New(), entity.Name, entity.Value, entity.Description, entity.Flag)
		if err != nil {
			return nil, e.Wrap(operation, err)
		}

		// err := s.db.Exec(ctx, query, entity.Name, entity.Value, entity.Description, entity.Flag)
		// if err != nil {
		// 	return nil, e.Wrap(operation, err)
		// }
	}

	if err := batch.Send(); err != nil {
		return nil, e.Wrap(operation, err)
	}

	for _, entity := range entitis {
		id, err := s.getId(ctx, entity)
		if err != nil {
			return nil, e.Wrap(operation, err)
		}
		ids = append(ids, id)
	}

	return ids, nil
}
