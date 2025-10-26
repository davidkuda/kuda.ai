package models

import (
	"database/sql"
	"errors"
	"fmt"
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

	stmt := `INSERT INTO auth.users (email, hashed_password) VALUES($1, $2);`

	result, err := m.DB.Exec(stmt, email, string(hashedPassword))
	if err != nil {
		return err
	}
	result.RowsAffected()

	return nil
}

// Creates a new user in the database
func (m *UserModel) UpdatePassword(email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 15)
	if err != nil {
		return err
	}

	stmt := `UPDATE auth.users SET hashed_password = $1 WHERE email = $2;`

	result, err := m.DB.Exec(stmt, string(hashedPassword), email)
	if err != nil {
		return err
	}
	r, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if r != 1 {
		return fmt.Errorf("rows affected: %d", r)
	}

	return nil
}


func (m *UserModel) Authenticate(email, password string) error {
	var hashedPassword []byte

	stmt := "SELECT hashed_password FROM auth.users WHERE email = $1;"

	err := m.DB.QueryRow(stmt, email).Scan(&hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrInvalidCredentials
		} else {
			return fmt.Errorf("DB.QueryRow(): %v", err)
		}
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return ErrInvalidCredentials
		} else {
			return fmt.Errorf("failed to compare password: %v", err)
		}
	}

	return nil
}

func (m *UserModel) Exists(email string) (bool, error) {
	return false, nil
}

func (m *UserModel) GetUserIDByEmail(email string) (int, error) {
	stmt := "SELECT id FROM auth.users WHERE email = $1;"

	var userID int

	err := m.DB.QueryRow(stmt, email).Scan(&userID)
	if err != nil {
		return 0, fmt.Errorf("DB.QueryRow(): %v", err)
	}

	return userID, nil
}
