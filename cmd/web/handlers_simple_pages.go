package main

import (
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/russross/blackfriday/v2"
)

type Page struct {
	Title       string
	Content     template.HTML
	CurrentPath string
}

func newPageFromMarkdown(markdown []byte, r *http.Request) *Page {
	htmlBytes := blackfriday.Run(markdown)
	content := template.HTML(htmlBytes)
	pageData := newPage(content, r)
	return &pageData
}

func newPage(content template.HTML, r *http.Request) Page {
	return Page{
		Title:       getTitleFromRequestPath(r),
		Content:     content,
		CurrentPath: getRootPath(r.URL.Path),
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Server", "Go")
	w.Header().Add("Creation-Month-Year", "April-2024")
	http.Redirect(w, r, "/about", http.StatusSeeOther)
}

func getPageAbout(w http.ResponseWriter, r *http.Request) {
	md, err := os.ReadFile("./data/pages/about.md")
	if err != nil {
		log.Printf("Error reading file: %s", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	pageData := newPageFromMarkdown(md, r)
	renderPageSimple(w, pageData)
}

func (app *application) blog(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, 200, "simplePage.html.tmpl", nil)
}

func getPageBlog(w http.ResponseWriter, r *http.Request) {
	md, err := os.ReadFile("./data/pages/blog.md")
	if err != nil {
		log.Printf("Error reading file: %s", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	pageData := newPageFromMarkdown(md, r)
	renderPageSimple(w, pageData)
}

func getPageBookshelf(w http.ResponseWriter, r *http.Request) {
	md, err := os.ReadFile("./data/pages/bookshelf.md")
	if err != nil {
		log.Printf("Error reading file: %s", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	pageData := newPageFromMarkdown(md, r)
	renderPageSimple(w, pageData)
}

func getPageCV(w http.ResponseWriter, r *http.Request) {
	// TODO: wouldn't it be nice to generate a beautiful PDF from this site?
	md, err := os.ReadFile("./data/pages/CV.md")
	if err != nil {
		log.Printf("Error reading file: %s", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	pageData := newPageFromMarkdown(md, r)
	pageData.Title = "CV"
	renderPageSimple(w, pageData)
}

func getPageTIL(w http.ResponseWriter, r *http.Request) {
	md, err := os.ReadFile("./data/pages/til.md")
	if err != nil {
		log.Printf("Error reading file: %s", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	pageData := newPageFromMarkdown(md, r)
	pageData.Title = "Today I Learned"
	renderPageSimple(w, pageData)
}
