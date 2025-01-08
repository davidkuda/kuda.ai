package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"

	"github.com/davidkuda/kudaai/internal/envcfg"
	"github.com/davidkuda/kudaai/internal/models"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type application struct {
	songs *models.SongModel
	users *models.UserModel

	templateCache map[string]*template.Template
	markdownHTMLCache map[string]template.HTML

	JWT struct {
		Secret []byte
	}
}

func main() {
	addr := flag.String("addr", ":8873", "HTTP network address")

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

	templateCache, err := newTemplateCache()
	if err != nil {
		log.Fatalf("could not initialise templateCache: %v\n", err)
	}
	app.templateCache = templateCache

	markdownHTMLCache, err := newMarkdownHTMLCache()
	if err != nil {
		log.Fatalf("could not initialise markdownHTMLCache: %v\n", err)
	}
	app.markdownHTMLCache = markdownHTMLCache

	log.Print("Starting web server, listening on port 8873")
	err = http.ListenAndServe(*addr, app.routes())
	log.Fatal(err)
}
