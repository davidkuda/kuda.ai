package main

import "net/http"

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("GET /{$}", home)
	mux.HandleFunc("GET /about", getPageAbout)
	mux.HandleFunc("GET /blog", getPageBlog)
	mux.HandleFunc("GET /bookshelf", getPageBookshelf)
	mux.HandleFunc("GET /cv", getPageCV)
	mux.HandleFunc("GET /songbook", app.getSongbook)
	mux.HandleFunc("GET /songbook/{song}", app.songbookSong)
	mux.HandleFunc("POST /songbook", app.songbookPost)
	mux.HandleFunc("GET /til", getPageTIL)
	mux.HandleFunc("GET /today-i-learned", getPageTIL)

	mux.HandleFunc("GET /admin/new-song", app.adminNewSong)

	return mux
}
