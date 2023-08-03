package storage

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
)

type Storage interface {
	Save(ctx context.Context, p *Page) error
	PickRandom(ctx context.Context, userName string) (*Page, error)
	Remove(ctx context.Context, p *Page) error
	IsExists(ctx context.Context, p *Page) (bool, error)
}

var ErrNoSavedPage = errors.New("no saved page")

type Page struct {
	URL      string
	UserName string
}

func (p Page) Hash() (string, error) {
	// h := sha1.New()

	// if _, err := io.WriteString(h, p.URL); err != nil {
	// 	return "", e.Wrap("can't calculated hash", err)
	// }

	// if _, err := io.WriteString(h, p.UserName); err != nil {
	// 	return "", e.Wrap("can't calculated hash", err)
	// }

	// return string(h.Sum(nil)), nil
	return GenerateId(), nil
}

func GenerateId() string {
	b := make([]byte, 32)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
