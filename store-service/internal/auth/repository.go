package auth

import (
	"github.com/jmoiron/sqlx"
)

type UserRepository interface {
	FindByUsername(username string) (User, error)
	FindByID(uid int) (UserPayload, error)
}

type UserRepositoryMySQL struct {
	DBConnection *sqlx.DB
}

func (repository UserRepositoryMySQL) FindByUsername(username string) (User, error) {
	var user User
	query := `SELECT id, username, first_name, last_name, password FROM users WHERE username = ?`
	err := repository.DBConnection.Get(&user, query, username)
	return user, err
}

func (repository UserRepositoryMySQL) FindByID(uid int) (UserPayload, error) {
	var user UserPayload
	query := `SELECT id, username, first_name, last_name FROM users WHERE id = ?`
	err := repository.DBConnection.Get(&user, query, uid)
	return user, err
}
