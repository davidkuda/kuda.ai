package main

import "net/http"

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("GET /{$}", home)

	// simple pages:
	mux.HandleFunc("GET /about", app.about)
	mux.HandleFunc("GET /blog", app.blog)
	mux.HandleFunc("GET /bookshelf", app.bookshelf)
	mux.HandleFunc("GET /cv", app.cv)
	mux.HandleFunc("GET /today-i-learned", app.til)

	//songbook:
	mux.HandleFunc("GET /songbook", app.songbook)
	mux.HandleFunc("GET /songbook/{song}", app.songbookSong)
	mux.HandleFunc("POST /songbook", app.songbookPost)

	// admin area:
	mux.HandleFunc("GET /admin", app.admin)
	mux.HandleFunc("GET /admin/login", app.adminLogin)
	mux.HandleFunc("POST /admin/login", app.adminLoginPost)
	mux.HandleFunc("GET /admin/new-song", app.adminNewSong)

	// finances:
	mux.HandleFunc("GET /finances", app.finances)

	return mux
}
