package envcfg

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"log"
	"net/url"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type envcfg struct {
	JWT
}

type JWT struct {
	Secret []byte
}

type db struct {
	Scheme   string
	Address  string
	Name     string
	User     string
	Password string
}

func Get() (*envcfg, error) {
	cfg := envcfg{}

	// JWT:
	jwt := JWT{
		Secret: make([]byte, 0, 64),
	}
	jwtSecretBase64 := os.Getenv("JWT_SECRET_KEY")
	if jwtSecretBase64 == "" {
		return nil, fmt.Errorf("make sure to define JWT_SECRET_KEY in the environment.")
	}
	_, err := base64.RawStdEncoding.Decode(jwt.Secret, []byte(jwtSecretBase64))
	if err != nil {
		return nil, fmt.Errorf("could not decode base64 env var JWT_SECRET_KEY: %v\n", err)
	}
	cfg.JWT = jwt

	return &cfg, nil
}

func DB() (*sql.DB, error) {
	c := db{
		Scheme:   os.Getenv("DB_SCHEME"),
		Address:  os.Getenv("DB_ADDRESS"),
		Name:     os.Getenv("DB_NAME"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
	}

	var fail bool

	if c.Scheme == "" {
		fail = true
		log.Print("Could not parse env var DB_SCHEME (e.g. postgres or mysql)")
	}

	if c.Address == "" {
		fail = true
		log.Print("Could not parse env var DB_ADDRESS")
	}

	if c.Name == "" {
		fail = true
		log.Print("Could not parse env var DB_NAME")
	}

	if c.User == "" {
		fail = true
		log.Print("Could not parse env var DB_USER")
	}

	if c.Password == "" {
		log.Print("Warning: DB_PASSWORD not set")
	}

	if fail {
		os.Exit(1)
	}

	dsn := url.URL{
		Scheme: c.Scheme,
		Host:   c.Address,
		User:   url.UserPassword(c.User, c.Password),
		Path:   c.Name,
	}

	db, err := sql.Open("pgx", dsn.String())
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
