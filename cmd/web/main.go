package main

import (
	"log"
	"net/http"

	"github.com/davidkuda/kudaai/internal/envcfg"
	"github.com/davidkuda/kudaai/internal/models"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type application struct {
	songs *models.SongModel
	users *models.UserModel

	JWT struct {
		Secret []byte
	}
}

func main() {
	app := &application{}

	c := envcfg.Get()

	app.JWT = c.JWT

	db, err := envcfg.DB()
	if err != nil {
		log.Fatalf("could not open DB: %v\n", err)
	}
	defer db.Close()

	app.songs = &models.SongModel{DB: db}
	app.users = &models.UserModel{DB: db}

	log.Print("Starting web server, listening on port 8873")
	err = http.ListenAndServe(":8873", app.routes())
	log.Fatal(err)
}
