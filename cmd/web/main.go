package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/davidkuda/kudaai/internal/envcfg"
	"github.com/davidkuda/kudaai/internal/models"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type application struct {
	navItems []NavItem

	models models.Models
	songs  *models.SongModel
	users  *models.UserModel
	til    *models.TILModel
	blogs  *models.BlogModel

	idLimits map[string]int

	templateCache     map[string]*template.Template
	templateCacheHTMX map[string]*template.Template

	JWT struct {
		Secret       []byte
		CookieDomain string
	}
}

type NavItem struct {
	Name string
	Path string
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	addr := flag.String("addr", ":8873", "HTTP network address")
	// TODO: cookieDomain should be defined in envcfg or envcfg should be dropped. Can't decide now and want to keep focusing on other tasks (:
	cookieDomain := flag.String("cookie-domain", os.Getenv("COOKIE_DOMAIN"), "localhost or kuda.ai")
	if *cookieDomain == "" {
		log.Fatal("fail startup: make sure to either pass -cookie-domain [localhost|kuda.ai] or define env var COOKIE_DOMAIN")
	}

	app := &application{}

	app.navItems = []NavItem{
		{Name: "Now", Path: "/now"},
		{Name: "About", Path: "/about"},
		{Name: "Blog", Path: "/blog"},
		{Name: "Bookshelf", Path: "/bookshelf"},
		{Name: "Songbook", Path: "/songbook"},
		{Name: "TIL", Path: "/today-i-learned"},
	}

	c := envcfg.Get()

	app.JWT.Secret = c.JWT.Secret
	app.JWT.CookieDomain = *cookieDomain

	db, err := envcfg.DB()
	if err != nil {
		log.Fatalf("could not open DB: %v\n", err)
	}
	defer db.Close()

	app.models = models.New(db)
	// TODO: replace all individual models with calls to app.models.M.func
	// NOTE: some learning moment:
	// after thinking about it, pointers are not necessary for the models.
	// No mutation. No big, heavy structs.
	app.songs = &models.SongModel{DB: db}
	app.users = &models.UserModel{DB: db}
	app.til = &models.TILModel{DB: db}
	app.blogs = &models.BlogModel{DB: db}

	templateCache, err := newTemplateCache()
	if err != nil {
		log.Fatalf("could not initialise templateCache: %v\n", err)
	}
	app.templateCache = templateCache
	templateCacheHTMX, err := newTemplateCacheForHTMXPartials()
	if err != nil {
		log.Fatalf("could not initialise templateCache: %v\n", err)
	}
	app.templateCacheHTMX = templateCacheHTMX

	log.Print("Starting web server, listening on port 8873")
	err = http.ListenAndServe(*addr, app.routes())
	log.Fatal(err)
}
