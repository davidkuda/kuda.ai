package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/russross/blackfriday/v2"

	"github.com/davidkuda/kudaai/internal/models"
)

func (app *application) getSongbook(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Server", "Go")
	w.Header().Add("Creation-Month-Year", "April-2024")

	allSongs, err := app.songs.GetAllSongs()
	if err != nil {
		log.Printf("Failed getting all songs: %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
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
	s := models.Song{
		ID:     song,
		Artist: "Sting",
		Name:   "Englishman In New York",
		// Lyrics: template.HTML(html),
		Lyrics: string(html),
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

func (app *application) songbookPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Printf("Failed parsing form: %v", err)
		// TODO: send status 400 Bad Request to the client
		return
	}

	f := r.PostForm
	s := models.Song{
		ID:        f.Get("song-id"),
		Artist:    f.Get("song-artist"),
		Name:      f.Get("song-name"),
		Lyrics:    f.Get("song-lyrics"),
		Chords:    f.Get("song-chords"),
		Copyright: f.Get("song-copyright"),
		MyCover:   f.Get("song-my-cover"),
	}

	app.songs.Insert(&s)
	// w.WriteHeader(http.StatusCreated)
	http.Redirect(w, r, fmt.Sprintf("/songbook/%v", s.ID), http.StatusSeeOther)
	return
}
