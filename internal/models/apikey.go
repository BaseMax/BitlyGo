package models

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/itsjoniur/bitlygo/internal/durable"
)

type ApiKey struct {
	UserID    int
	Key       uuid.UUID
	CreatedAt time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time
}

func CreateApiKey(ctx context.Context, user *User) (*ApiKey, error) {
	db := ctx.Value(10).(*durable.Database)
	now := time.Now()
	newApiKey := &ApiKey{
		UserID:    user.ID,
		Key:       uuid.New(),
		CreatedAt: now,
	}

	query := "INSERT INTO api_keys(user_id, key) VALUES($1, $2)"
	values := []interface{}{newApiKey.UserID, newApiKey.Key}
	_, err := db.Exec(context.Background(), query, values...)
	if err != nil {
		return nil, err
	}

	return newApiKey, nil
}
