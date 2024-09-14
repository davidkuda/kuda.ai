package models

import (
	"database/sql"
	"errors"
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

// Creates a new user in the database
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
    var hashedPassword []byte

    stmt := "SELECT email, hashed_password FROM users WHERE email = $1;"

    err := m.DB.QueryRow(stmt, email).Scan(&hashedPassword)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return ErrInvalidCredentials
        } else {
            return err
        }
    }

    err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
    if err != nil {
        if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
            return ErrInvalidCredentials
        } else {
            return err
        }
    }

    return nil
}

func (m *UserModel) Exists(email string) (bool, error) {
	return false, nil
}
