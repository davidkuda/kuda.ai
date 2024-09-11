package main

import "net/http"

func routes() *http.ServeMux {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("GET /{$}", home)
	mux.HandleFunc("GET /about", getPageAbout)
	mux.HandleFunc("GET /blog", getPageBlog)
	mux.HandleFunc("GET /bookshelf", getPageBookshelf)
	mux.HandleFunc("GET /cv", getPageCV)
	mux.HandleFunc("GET /songbook", getSongbook)
	mux.HandleFunc("GET /songbook/{song}", getSongbookSong)
	mux.HandleFunc("POST /songbook/add", postSongbookAdd)
	mux.HandleFunc("GET /til", getPageTIL)
	mux.HandleFunc("GET /today-i-learned", getPageTIL)

	mux.HandleFunc("GET /admin/new-song", app.getAdminNewSong)

	return mux
}
