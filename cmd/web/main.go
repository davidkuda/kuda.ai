package main

import (
	"database/sql"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/davidkuda/kudaai/internal/models"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type application struct {
	songs *models.SongModel
	users *models.UserModel
}

type ConfigFromEnv struct {
	ListenAddress string

	DBScheme   string
	DBAddress  string
	DBName     string
	DBUser     string
	DBPassword string
}

func main() {
	c := getConfigFromEnv()

	db, err := openDB(c)
	if err != nil {
		log.Fatal(err.Error())
		os.Exit(1)
	}
	defer db.Close()

	app := &application{
		songs: &models.SongModel{DB: db},
		users: &models.UserModel{DB: db},
	}

	log.Print("Starting web server, listening on port 8873")
	err = http.ListenAndServe(":8873", app.routes())
	log.Fatal(err)
}

func getConfigFromEnv() *ConfigFromEnv {
	c := ConfigFromEnv{
		ListenAddress: os.Getenv("LISTEN_ADDRESS"),

		DBScheme:   os.Getenv("DB_SCHEME"),
		DBAddress:  os.Getenv("DB_ADDRESS"),
		DBName:     os.Getenv("DB_NAME"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
	}

	var fail bool

	if c.ListenAddress == "" {
		c.ListenAddress = ":8873"
	}

	if c.DBScheme == "" {
		fail = true
		log.Fatal("Could not parse env var DB_SCHEME (e.g. postgres or mysql)")
	}

	if c.DBAddress == "" {
		fail = true
		log.Fatal("Could not parse env var DB_ADDRESS")
	}

	if c.DBName == "" {
		fail = true
		log.Fatal("Could not parse env var DB_NAME")
	}

	if c.DBUser == "" {
		fail = true
		log.Fatal("Could not parse env var DB_USER")
	}

	// if c.DBPassword == "" {
	// 	fail = true
	// 	log.Fatal("Could not parse env var DB_PASSWORD")
	// }

	if fail {
		os.Exit(1)
	}

	return &c
}

func openDB(c *ConfigFromEnv) (*sql.DB, error) {
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
