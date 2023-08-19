package mongodb

import (
	"context"
	"fmt"
	"sql-exp-test/internal/lib/e"
	"sql-exp-test/internal/storage"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Storage struct {
	db          *mongo.Client
	colEntities *mongo.Collection
}

type tEnt struct {
	Id          any     `bson:"_id"`
	Name        string  `bson:"name"`
	Value       float64 `bson:"value"`
	Description string  `bson:"description"`
	Flag        bool    `bson:"flag"`
}

// New creates new MongoDB storage.
func New(path string) (storage.Storage, error) {
	const operation = "storage.mongodb.New"

	db, err := mongo.Connect(context.Background(), options.Client().ApplyURI(path))
	if err != nil {
		return nil, e.Wrap(operation, err)
	}

	err = db.Ping(context.Background(), readpref.Primary())
	if err != nil {
		return nil, e.Wrap(operation, err)
	}

	return &Storage{
		db: db,
	}, nil
}

// Close MongoDB storage.
func (s *Storage) Close() error {
	const operation = "storage.mongodb.Close"

	return e.WrapIfErr(operation, s.db.Disconnect(context.Background()))
}

// Create table if not exists
func (s *Storage) Init(ctx context.Context) error {
	// const operation = "storage.mongodb.Init"

	colEntities := s.db.Database("entities").Collection("entities")
	s.colEntities = colEntities

	return nil
}

// Create entity to storage.
func (s *Storage) Create(ctx context.Context, entity *storage.Entities) (any, error) {
	const operation = "storage.mongodb.Create"

	res, err := s.colEntities.InsertOne(ctx,
		bson.M{
			"name":        entity.Name,
			"value":       entity.Value,
			"description": entity.Description,
			"flag":        entity.Flag,
		})
	if err != nil {
		return nil, e.Wrap(operation, err)
	}

	id := res.InsertedID

	return id, nil
}

// Read entity from storage
func (s *Storage) Read(ctx context.Context, id any) (*storage.Entities, error) {
	const operation = "storage.mongodb.Read"

	var ent tEnt

	err := s.colEntities.FindOne(ctx, bson.M{
		"_id": id,
	}).Decode(&ent)
	if err == mongo.ErrNoDocuments {
		return nil, e.Wrap(operation, fmt.Errorf("record does not exist"))
	} else if err != nil {
		return nil, e.Wrap(operation, err)
	}

	return &storage.Entities{
		Id:          ent.Id,
		Name:        ent.Name,
		Value:       ent.Value,
		Description: ent.Description,
		Flag:        ent.Flag,
	}, nil
}

// Update entity to storage
func (s *Storage) Update(ctx context.Context, entity *storage.Entities) error {
	const operation = "storage.mongodb.Update"

	_, err := s.colEntities.UpdateOne(ctx,
		bson.M{
			"_id": entity.Id,
		},
		bson.M{
			"$set": bson.M{
				"name":        entity.Name,
				"value":       entity.Value,
				"description": entity.Description,
				"flag":        entity.Flag,
			},
		})
	if err != nil {
		return e.Wrap(operation, err)
	}

	return nil
}

// Remove entity from storage by entity.
func (s *Storage) Delete(ctx context.Context, entity *storage.Entities) error {
	const operation = "storage.mongodb.Delete"

	_, err := s.colEntities.DeleteOne(ctx, bson.M{
		"name":        entity.Name,
		"value":       entity.Value,
		"description": entity.Description,
		"flag":        entity.Flag,
	})
	if err != nil {
		return e.Wrap(operation, err)
	}

	return nil
}

// Remove entity from storage by Id
func (s *Storage) DeleteId(ctx context.Context, id any) error {
	const operation = "storage.mongodb.DeleteId"

	_, err := s.colEntities.DeleteOne(ctx, bson.M{
		"_id": id,
	})
	if err != nil {
		return e.Wrap(operation, err)
	}

	return nil
}

// IsExists checks if entity exists in storage.
func (s *Storage) IsExists(ctx context.Context, entity *storage.Entities) (bool, error) {
	const operation = "storage.mongodb.IsExists"

	res, err := s.colEntities.Find(ctx, bson.M{
		"name":        entity.Name,
		"value":       entity.Value,
		"description": entity.Description,
		"flag":        entity.Flag,
	})
	if err != nil {
		return false, e.Wrap(operation, err)
	}

	count := res.RemainingBatchLength()

	return count > 0, nil
}

// IsExistsById checks if entity exists in storage by Id.
func (s *Storage) IsExistsById(ctx context.Context, id any) (bool, error) {
	const operation = "storage.mongodb.IsExistsById"

	res, err := s.colEntities.Find(ctx, bson.M{
		"_id": id,
	})
	if err != nil {
		return false, e.Wrap(operation, err)
	}

	count := res.RemainingBatchLength()

	return count > 0, nil
}

// Lots of records.
func (s *Storage) LotsOfRecords(ctx context.Context, entitis ...*storage.Entities) ([]any, error) {
	const operation = "storage.mongodb.LotsOfRecords"

	var ids []any = []any{}

	for _, entity := range entitis {

		res, err := s.colEntities.InsertOne(ctx,
			bson.M{
				"name":        entity.Name,
				"value":       entity.Value,
				"description": entity.Description,
				"flag":        entity.Flag,
			})
		if err != nil {
			return nil, e.Wrap(operation, err)
		}

		id := res.InsertedID

		ids = append(ids, id)
	}

	return ids, nil
}
