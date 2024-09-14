package models

import (
	"database/sql"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Email          string
	HashedPassword []byte
	Created        time.Time
}

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 15)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO users (email, hashed_password, created)
    VALUES($1, $2, UTC_TIMESTAMP());`

	_, err = m.DB.Exec(stmt, email, string(hashedPassword))
	if err != nil {
		return err
	}

	return nil
}

func (m *UserModel) Authenticate(email, password string) error {
	return nil
}

func (m *UserModel) Exists(email string) (bool, error) {
	return false, nil
}
