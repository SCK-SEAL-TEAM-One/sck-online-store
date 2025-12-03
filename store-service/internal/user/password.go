package user

import (
	"golang.org/x/crypto/bcrypt"
)

type PasswordHelper interface {
	CheckPasswordHash(password, hashed string) bool
	HashPassword(password string) (string, error)
}

type BcryptPasswordChecker struct{}

func (BcryptPasswordChecker) CheckPasswordHash(password, hashed string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
	return err == nil
}

func (BcryptPasswordChecker) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}
