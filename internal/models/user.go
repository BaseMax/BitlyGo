package models

import (
	"context"
	"errors"
	"time"

	"github.com/itsjoniur/bitlygo/internal/durable"
	"github.com/itsjoniur/bitlygo/pkg/auth"
)

type User struct {
	ID        int
	Username  string
	Password  string
	CreatedAt time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time
}

func CreateUser(ctx context.Context, username, password string) (*User, error) {
	db := ctx.Value(10).(*durable.Database)
	now := time.Now()
	var err error
	var hashedPassword string

	hashedPassword, err = auth.HashPassword(password)
	if err != nil {
		return nil, err
	}

	newUser := &User{
		Username:  username,
		Password:  hashedPassword,
		CreatedAt: now,
	}
	query := "INSERT INTO users(username, password, created_at) VALUES($1, $2, $3) RETURNING id"
	values := []interface{}{newUser.Username, newUser.Password, newUser.CreatedAt}

	db.QueryRow(context.Background(), query, values...).Scan(&newUser.ID)
	if newUser.ID == 0 {
		return nil, errors.New("user not created")
	}
	return newUser, nil
}

func GetUserByUsername(ctx context.Context, username string) *User {
	db := ctx.Value(10).(*durable.Database)
	user := &User{}

	query := "SELECT id, username, password FROM users WHERE username = $1"
	db.QueryRow(context.Background(), query, username).Scan(&user.ID, &user.Username, &user.Password)

	if user.Username == "" && user.Password == "" {
		return nil
	}

	return user
}

func GetUserByApiKey(ctx context.Context, apiKey string) *User {
	db := ctx.Value(10).(*durable.Database)
	user := &User{}

	query := `
		SELECT id, username
		FROM users
		WHERE id = (
			SELECT user_id FROM api_keys WHERE key = $1
		);
	`
	db.QueryRow(context.Background(), query, apiKey).Scan(user.ID, user.Username)

	if user.Username == "" {
		return nil
	}

	return user
}

func UpdateUserByUsername(ctx context.Context, username, newUsername, newPassword string) (*User, error) {
	db := ctx.Value(10).(*durable.Database)
	user := &User{
		Username: newUsername,
		Password: newPassword,
	}

	query := "UPDATE users SET username = COALESCE(NULLIF($1, ''), username), password = COALESCE(NULLIF($2, ''), password) WHERE username = $3"
	values := []interface{}{user.Username, user.Password, username}
	_, err := db.Exec(context.Background(), query, values...)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func DeleteUserByUsername(ctx context.Context, username string) error {
	db := ctx.Value(10).(*durable.Database)

	query := "DELETE FROM users WHERE username = $1"
	_, err := db.Exec(context.Background(), query)

	return err
}
