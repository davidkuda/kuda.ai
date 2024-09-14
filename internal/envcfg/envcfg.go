package envcfg

import (
	"database/sql"
	"log"
	"net/url"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type envcfg struct {
	DBScheme   string
	DBAddress  string
	DBName     string
	DBUser     string
	DBPassword string
}

func DB() (*sql.DB, error) {
	c := envcfg{
		DBScheme:   os.Getenv("DB_SCHEME"),
		DBAddress:  os.Getenv("DB_ADDRESS"),
		DBName:     os.Getenv("DB_NAME"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
	}

	var fail bool

	if c.DBScheme == "" {
		fail = true
		log.Print("Could not parse env var DB_SCHEME (e.g. postgres or mysql)")
	}

	if c.DBAddress == "" {
		fail = true
		log.Print("Could not parse env var DB_ADDRESS")
	}

	if c.DBName == "" {
		fail = true
		log.Print("Could not parse env var DB_NAME")
	}

	if c.DBUser == "" {
		fail = true
		log.Print("Could not parse env var DB_USER")
	}

	if c.DBPassword == "" {
		log.Print("Warning: DB_PASSWORD not set")
	}

	if fail {
		os.Exit(1)
	}

	dsn := url.URL{
		Scheme: c.DBScheme,
		Host:   c.DBAddress,
		User:   url.UserPassword(c.DBUser, c.DBPassword),
		Path:   c.DBName,
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
