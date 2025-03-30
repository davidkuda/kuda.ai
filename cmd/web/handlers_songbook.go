package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/russross/blackfriday/v2"

	"github.com/davidkuda/kudaai/internal/models"
)
func (app *application) songbook(w http.ResponseWriter, r *http.Request) {
	allSongs, err := app.songs.GetAllSongs()
	if err != nil {
		log.Printf("Failed getting all songs: %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	t := app.newTemplateData(r)
	t.Songs = allSongs

	app.render(w, r, 200, "songbook.tmpl.html", &t)
}

func (app *application) songbookSong(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Server", "Go")
	w.Header().Add("Creation-Month-Year", "August-2024")

	songName := r.PathValue("song")

	song, err := app.songs.Get(songName)
	if err != nil {
		log.Printf("could not get song \"%v\": %v", songName, err)
		// TODO: Show a nice 404 page.
		http.NotFound(w, r)
		return
	}

	song.HTML.Lyrics = template.HTML(blackfriday.Run([]byte(song.Lyrics)))
	song.HTML.Chords = template.HTML(blackfriday.Run([]byte(song.Chords)))

	t := app.newTemplateData(r)
	t.Song = song
	t.Title = "Songbook: " + song.Name + " (" + song.Artist + ")"

	app.render(w, r, 200, "song.tmpl.html", &t)
}

func (app *application) adminNewSong(w http.ResponseWriter, r *http.Request) {
	t := app.newTemplateData(r)
	app.render(w, r, http.StatusOK, "admin.new_song.tmpl.html", &t)
}

type songbookSongForm struct {
	Song        models.Song
	FieldErrors map[string]string
}

func (app *application) songbookPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Printf("Failed parsing form: %v", err)
		// TODO: send status 400 Bad Request to the client
		return
	}

	f := r.PostForm

	form := songbookSongForm{
		Song: models.Song{
			ID:        f.Get("song-id"),
			Artist:    f.Get("song-artist"),
			Name:      f.Get("song-name"),
			Lyrics:    strings.ReplaceAll(f.Get("song-lyrics"), "\r\n", "\n"),
			Chords:    strings.ReplaceAll(f.Get("song-chords"), "\r\n", "\n"),
			Copyright: f.Get("song-copyright"),
			MyCover:   f.Get("song-my-cover"),
		},
		FieldErrors: map[string]string{},
	}

	// regex for valid URL path; song.ID will be used in the URL.
	// Therefore, it should only contain letters and hyphens.
	var rxPat = regexp.MustCompile(`[^a-z\-]*`)
	if !rxPat.MatchString(form.Song.ID) {
		form.FieldErrors["song-id"] = "SongID may only contain lowercase characters and hyphens"
	}

	 if len(form.FieldErrors) > 0 {
        t := app.newTemplateData(r)
        t.Form = form
        app.render(w, r, http.StatusUnprocessableEntity, "create.tmpl", &t)
        return
    }

	app.songs.Insert(&form.Song)
	// w.WriteHeader(http.StatusCreated)
	http.Redirect(w, r, fmt.Sprintf("/songbook/%v", form.Song.ID), http.StatusSeeOther)
	return
}
