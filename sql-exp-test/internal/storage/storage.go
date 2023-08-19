package storage

import (
	"context"
	"errors"
)

var (
	ErrURLNotFound = errors.New("url not found")
	ErrURLExists   = errors.New("url exists")
	ErrNoSavedPage = errors.New("no saved page")
)

type Storage interface {
	Close() error
	Create(ctx context.Context, entity *Entities) (any, error)
	Init(ctx context.Context) error
	Read(ctx context.Context, id any) (*Entities, error)
	Update(ctx context.Context, entity *Entities) error
	Delete(ctx context.Context, entity *Entities) error
	DeleteId(ctx context.Context, id any) error
	IsExists(ctx context.Context, entity *Entities) (bool, error)
	IsExistsById(ctx context.Context, id any) (bool, error)
	LotsOfRecords(ctx context.Context, entitis ...*Entities) ([]any, error)
}

type Entities struct {
	Id          any
	Name        string
	Value       float64
	Description string
	Flag        bool
}
