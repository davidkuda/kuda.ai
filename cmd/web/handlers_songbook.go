package main

import (
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/russross/blackfriday/v2"
)

type Songs []Song

type Song struct {
	ID     string
	Name   string
	Artist string
	Lyrics template.HTML
}

func getSongbook(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Server", "Go")
	w.Header().Add("Creation-Month-Year", "April-2024")

	allSongs := Songs{
		Song{
			ID:     "englishman-in-new-york",
			Name:   "Englishman In New York",
			Artist: "Sting",
		},
		Song{
			ID:     "englishman-in-new-york",
			Name:   "Englishman In New York",
			Artist: "Sting",
		},
	}

	tmplFiles := []string{
		"./ui/html/pages/base.tmpl.html",
		"./ui/html/partials/nav.tmpl.html",
		"./ui/html/pages/songbook.tmpl.html",
	}

	t, err := template.ParseFiles(tmplFiles...)
	if err != nil {
		log.Printf("Error parsing home.tmpl.html: %s", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = t.ExecuteTemplate(w, "base", allSongs)
	if err != nil {
		log.Printf("Error executing home.tmpl.html: %s", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func getSongbookSong(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Server", "Go")
	w.Header().Add("Creation-Month-Year", "August-2024")

	song := r.PathValue("song")

	availableSongs := map[string]bool{
		"englishman-in-new-york": true,
	}

	if !availableSongs[song] {
		http.NotFound(w, r)
	}

	lyrics, err := os.ReadFile("./data/songs/englishman-in-new-york.md")

	if err != nil {
		log.Printf("Error reading file: %s", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
	html := blackfriday.Run(lyrics)
	s := Song{
		ID:     song,
		Artist: "Sting",
		Name:   "Englishman In New York",
		Lyrics: template.HTML(html),
	}

	tmplFiles := []string{
		"./ui/html/pages/base.tmpl.html",
		"./ui/html/partials/nav.tmpl.html",
		"./ui/html/pages/song.tmpl.html",
	}

	t, err := template.ParseFiles(tmplFiles...)
	if err != nil {
		log.Printf("Error parsing home.tmpl.html: %s", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = t.ExecuteTemplate(w, "base", s)
	if err != nil {
		log.Printf("Error executing home.tmpl.html: %s", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

}

func getSongbookAdd(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display a form to create a new song"))
}

func postSongbookAdd(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Add a new song to the songbook ..."))
}
