package main

import (
	"fmt"
	"net/http"
)

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Server", "Go")
	w.Header().Add("Creation-Month-Year", "April-2024")
	w.Write([]byte("Hello from kudaai"))
}

func getSongbook(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display all songs"))
}

func getSongbookSong(w http.ResponseWriter, r *http.Request) {
	song := r.PathValue("song")

	availableSongs := map[string]bool{
		"englishman-in-new-york": true,
	}

	if !availableSongs[song] {
		http.NotFound(w, r)
	}

	fmt.Fprintf(w, "Requested Song: %s", song)
}

func getSongbookAdd(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display a form to create a new song"))
}

func postSongbookAdd(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Add a new song to the songbook ..."))
}
