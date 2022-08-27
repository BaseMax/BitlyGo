package models

import (
	"context"
	"fmt"
	"time"

	"github.com/itsjoniur/bitlygo/internal/durable"
)

const (
	ExpireTime = time.Hour * 48
)

type Link struct {
	Id        int
	OwnerId   int
	Name      string
	Link      string
	Visits    int
	CreatedAt time.Time
	UpdatedAt time.Time
	ExpiredAt *time.Time
	DeletedAt *time.Time
}

func CreateLink(ctx context.Context, owner int, name, link string) (*Link, error) {
	db := ctx.Value(10).(*durable.Database)
	now := time.Now()
	newLink := &Link{
		OwnerId:   owner,
		Name:      name,
		Link:      link,
		CreatedAt: now,
		UpdatedAt: now,
	}

	fmt.Println(newLink)
	query := "INSERT INTO links(owner_id, name, link, created_at, updated_at) VALUES($1, $2, $3, $4, $5)"
	values := []interface{}{newLink.OwnerId, newLink.Name, newLink.Link, newLink.CreatedAt, newLink.UpdatedAt}
	_, err := db.Exec(context.Background(), query, values...)
	if err != nil {
		return nil, err
	}
	return newLink, nil
}
