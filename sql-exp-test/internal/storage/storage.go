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
	Create(ctx context.Context, entity *Entities) error
	Read(ctx context.Context, id int64) (*Entities, error)
	Update(ctx context.Context, entity *Entities) error
	Delete(ctx context.Context, entity *Entities) error
	DeleteId(ctx context.Context, id int64) error
	IsExists(ctx context.Context, entity *Entities) (bool, error)
}

type Entities struct {
	Id          int64
	Name        string
	Value       float64
	Description string
	Flag        bool
}
