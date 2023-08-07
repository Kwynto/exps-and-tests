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
	Save(ctx context.Context, p *Page) error
	// PickRandom(ctx context.Context, userName string) (*Page, error)
	Remove(ctx context.Context, p *Page) error
	IsExists(ctx context.Context, p *Page) (bool, error)
}

type Page struct {
	URL      string
	UserName string
}
