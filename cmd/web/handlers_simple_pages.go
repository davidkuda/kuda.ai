package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/russross/blackfriday/v2"
)

type Page struct {
	Title       string
	Content     template.HTML
	CurrentPath string
}

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

	p := Page{
		Title:       "kuda.ai | home",
		Content:     "hello world",
		CurrentPath: getRootPath(r.URL.Path),
	}

	err = t.ExecuteTemplate(w, "base", p)
	if err != nil {
		log.Printf("Error executing home.tmpl.html: %s", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
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

func getTitleFromRequestPath(r *http.Request) string {
	// TODO: use "golang.org/x/text/cases" instead of strings
	return strings.Title(r.URL.Path[1:])
}

func getRootPath(path string) string {
	var i int
	for i = 1; i < len(path); i++ {
		if path[i] == '/' {
			break
		}
	}
	return path[0:i]
}

func renderPageSimple(w http.ResponseWriter, p *Page) {
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
