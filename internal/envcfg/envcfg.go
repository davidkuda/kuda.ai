package envcfg

import (
	"database/sql"
	"encoding/base64"
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

func Get() *envcfg {
	cfg := envcfg{}
	cfg.JWT = parseJWT()
	return &cfg
}

func parseJWT() JWT {
	jwt := JWT{}

	jwtSecretBase64 := os.Getenv("JWT_SECRET_KEY")
	if jwtSecretBase64 == "" {
		log.Fatal("make sure to define JWT_SECRET_KEY in the environment.")
	}

	s, err := base64.StdEncoding.DecodeString(jwtSecretBase64)
	if err != nil {
		log.Fatalf("could not decode base64 env var JWT_SECRET_KEY: %v\n", err)
	}

	jwt.Secret = s

	return jwt
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
