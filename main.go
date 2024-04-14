package main

import (
	"log"
	"net/http"
)

func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello from kudaai"))
}

func listSongs(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display all songs"))
}

func createSong(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display a form to create a new song"))
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/{$}", home)
	mux.HandleFunc("/songs/view", listSongs)
	mux.HandleFunc("/songs/create", createSong)

	log.Print("Starting web server, listening on port 8873")

	err := http.ListenAndServe(":8873", mux)
	log.Fatal(err)
}
