package auth

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type UserRepository interface {
	FindByUsername(ctx context.Context, username string) (User, error)
	FindByID(ctx context.Context, uid int) (UserPayload, error)
}

type UserRepositoryMySQL struct {
	DBConnection *sqlx.DB
}

func (repository UserRepositoryMySQL) FindByUsername(ctx context.Context, username string) (User, error) {
	var user User
	query := `SELECT id, username, first_name, last_name, password FROM users WHERE username = ?`
	err := repository.DBConnection.GetContext(ctx, &user, query, username)
	return user, err
}

func (repository UserRepositoryMySQL) FindByID(ctx context.Context, uid int) (UserPayload, error) {
	var user UserPayload
	query := `SELECT id, username, first_name, last_name FROM users WHERE id = ?`
	err := repository.DBConnection.GetContext(ctx, &user, query, uid)
	return user, err
}
