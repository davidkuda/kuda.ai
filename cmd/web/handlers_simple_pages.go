package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

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
		Title:   getTitleFromRequestPath(r),
		Content: template.HTML(htmlBytes),
	}

	getSimplePage(w, &pageData)
}

func getPageBlog(w http.ResponseWriter, r *http.Request) {
	md, err := os.ReadFile("./data/pages/blog.md")

	if err != nil {
		log.Printf("Error reading file: %s", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	htmlBytes := blackfriday.Run(md)
	pageData := Page{
		Title:   getTitleFromRequestPath(r),
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

	// TODO: How to highlight the nav element that is currently on?

	htmlBytes := blackfriday.Run(md)
	pageData := Page{
		Title:   getTitleFromRequestPath(r),
		Content: template.HTML(htmlBytes),
	}

	getSimplePage(w, &pageData)
}

func getPageCV(w http.ResponseWriter, r *http.Request) {
	// TODO: wouldn't it be nice to generate a beautiful PDF from this site?
	md, err := os.ReadFile("./data/pages/CV.md")

	if err != nil {
		log.Printf("Error reading file: %s", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	htmlBytes := blackfriday.Run(md)
	pageData := Page{
		Title:   "CV",
		Content: template.HTML(htmlBytes),
	}

	getSimplePage(w, &pageData)
}

func getTitleFromRequestPath(r *http.Request) string {
	// TODO: use "golang.org/x/text/cases" instead of strings
	return strings.Title(r.URL.Path[1:])
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
