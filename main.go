package main

import (
	"log"
	"net/http"
)
func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/{$}", home)
	mux.HandleFunc("/songbook", listSongs)
	mux.HandleFunc("/songbook/{song}", showSong)
	mux.HandleFunc("/songbook/add", addSong)

	log.Print("Starting web server, listening on port 8873")

	err := http.ListenAndServe(":8873", mux)
	log.Fatal(err)
}

func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello from kudaai"))
}

func listSongs(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display all songs"))
}

func showSong(w http.ResponseWriter, r *http.Request) {
	song := r.PathValue("song")
	w.Write([]byte("Requested Song " + song))
}

func addSong(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display a form to create a new song"))
}

