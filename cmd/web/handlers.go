package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/russross/blackfriday/v2"
)

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Server", "Go")
	w.Header().Add("Creation-Month-Year", "April-2024")

	tmplFiles := []string{
		"./ui/html/pages/base.tmpl.html",
		"./ui/html/partials/nav.tmpl.html",
		"./ui/html/pages/home.tmpl.html",
	}

	t, err := template.ParseFiles(tmplFiles...)
	if err != nil {
		log.Printf("Error parsing home.tmpl.html: %s", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = t.ExecuteTemplate(w, "base", nil)
	if err != nil {
		log.Printf("Error executing home.tmpl.html: %s", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

type Song struct {
	Lyrics template.HTML
}

type Page struct {
	Title   string
	Content template.HTML
}

func getPageAbout(w http.ResponseWriter, r *http.Request) {
	md, err := os.ReadFile("./data/pages/about.md")

	if err != nil {
		log.Printf("Error reading file: %s", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	htmlBytes := blackfriday.Run(md)
	pageData := Page{
		Title:   "About",
		Content: template.HTML(htmlBytes),
	}

	getSimplePage(w, &pageData)
}

func getPageBookshelf(w http.ResponseWriter, r *http.Request) {
	md, err := os.ReadFile("./data/pages/bookshelf.md")

	if err != nil {
		log.Printf("Error reading file: %s", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	// TODO: get Title from path
	// TODO: How to highlight the nav element that is currently on?

	htmlBytes := blackfriday.Run(md)
	pageData := Page{
		Title:   "Bookshelf",
		Content: template.HTML(htmlBytes),
	}

	getSimplePage(w, &pageData)
}

func getSimplePage(w http.ResponseWriter, p *Page) {
	w.Header().Add("Server", "Go")
	w.Header().Add("Creation-Month-Year", "August-2024")

	tmplFiles := []string{
		"./ui/html/pages/base.tmpl.html",
		"./ui/html/partials/nav.tmpl.html",
		"./ui/html/pages/simplePage.tmpl.html",
	}

	t, err := template.ParseFiles(tmplFiles...)
	if err != nil {
		log.Printf("Error parsing home.tmpl.html: %s", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = t.ExecuteTemplate(w, "base", p)
	if err != nil {
		log.Printf("Error executing home.tmpl.html: %s", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func getSongbook(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Server", "Go")
	w.Header().Add("Creation-Month-Year", "April-2024")

	lyrics, err := os.ReadFile("./data/songs/englishman-in-new-york.md")

	if err != nil {
		log.Printf("Error reading file: %s", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
	html := blackfriday.Run(lyrics)
	s := Song{template.HTML(html)}

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
