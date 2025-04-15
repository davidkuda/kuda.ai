package main

import (
	"net/http"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Server", "Go")
	w.Header().Add("Started-Working-On", "April-2024")
	t := app.newTemplateData(r)
	t.HTML = app.markdownHTMLCache["home.md"]
	app.render(w, r, 200, "home.tmpl.html", &t)
}

func (app *application) about(w http.ResponseWriter, r *http.Request) {
	t := app.newTemplateData(r)
	t.HTML = app.markdownHTMLCache["about.md"]
	// Title:    "About",
	// RootPath: "/about",
	app.render(w, r, 200, "simplePage.tmpl.html", &t)
}

func (app *application) blog(w http.ResponseWriter, r *http.Request) {
	t := app.newTemplateData(r)
	t.HTML = app.markdownHTMLCache["blog.md"]
	app.render(w, r, 200, "simplePage.tmpl.html", &t)
}

func (app *application) bookshelf(w http.ResponseWriter, r *http.Request) {
	t := app.newTemplateData(r)
	t.HTML = app.markdownHTMLCache["bookshelf.md"]
	app.render(w, r, 200, "simplePage.tmpl.html", &t)
}
