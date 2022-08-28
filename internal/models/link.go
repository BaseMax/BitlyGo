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

	query := "INSERT INTO links(owner_id, name, link, created_at, updated_at) VALUES($1, $2, $3, $4, $5)"
	values := []interface{}{newLink.OwnerId, newLink.Name, newLink.Link, newLink.CreatedAt, newLink.UpdatedAt}
	_, err := db.Exec(context.Background(), query, values...)
	if err != nil {
		return nil, err
	}
	return newLink, nil
}

func GetLinkByName(ctx context.Context, name string) *Link {
	db := ctx.Value(10).(*durable.Database)
	link := &Link{}

	query := "SELECT name, link FROM links WHERE name = $1"
	db.QueryRow(context.Background(), query, name).Scan(&link.Name, &link.Link)

	return link
}

func SearchLinkByName(ctx context.Context, name string, limit int) ([]*Link, error) {
	db := ctx.Value(10).(*durable.Database)
	links := []*Link{}

	query := fmt.Sprintf("SELECT name, link FROM links WHERE name LIKE '%%%v%%' LIMIT $1", name)
	rows, err := db.Query(context.Background(), query, limit)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		link := &Link{}

		err := rows.Scan(&link.Name, &link.Link)
		if err != nil {
			return nil, err
		}

		links = append(links, link)

	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return links, nil

}

func TopLinksByVisits(ctx context.Context, limit int) ([]*Link, error) {
	db := ctx.Value(10).(*durable.Database)
	links := []*Link{}

	query := "SELECT name, link, visits FROM links ORDER BY visits DESC LIMIT $1"
	rows, err := db.Query(context.Background(), query, limit)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		link := &Link{}

		err := rows.Scan(&link.Name, &link.Link, &link.Visits)
		if err != nil {
			return nil, err
		}

		links = append(links, link)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return links, nil
}

func UpdateLinkByName(ctx context.Context, name, newName, newLink string) (*Link, error) {
	db := ctx.Value(10).(*durable.Database)
	link := &Link{
		Name: newName,
		Link: newLink,
	}

	query := "UPDATE links SET name = COALESCE($1, name), link = $2 WHERE name = $3"
	values := []interface{}{link.Name, link.Link, name}
	_, err := db.Exec(context.Background(), query, values...)
	if err != nil {
		return nil, err
	}

	return link, nil
}

func DeleteLinkByName(ctx context.Context, name string) error {
	db := ctx.Value(10).(*durable.Database)

	query := "DELETE FROM links WHERE name = $1"
	_, err := db.Exec(context.Background(), query, name)

	return err
}
