
package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("./ui/static"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fs))

	mux.HandleFunc("GET /{$}", home)
	mux.HandleFunc("GET /songbook", getSongbook)
	mux.HandleFunc("GET /songbook/{song}", getSongbookSong)
	mux.HandleFunc("GET /songbook/add", getSongbookAdd)
	mux.HandleFunc("POST /songbook/add", postSongbookAdd)

	log.Print("Starting web server, listening on port 8873")

	err := http.ListenAndServe(":8873", mux)
	log.Fatal(err)
}

