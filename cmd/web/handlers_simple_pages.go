package main

import (
	"net/http"

	"github.com/davidkuda/kudaai/internal/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Server", "Go")
	w.Header().Add("Started-Working-On", "April-2024")
	t := app.newTemplateData(r)
	t.HTML = app.markdownHTMLCache["home.md"]
	app.render(w, r, 200, "home.tmpl.html", &t)
}

func (app *application) now(w http.ResponseWriter, r *http.Request) {
	t := app.newTemplateData(r)
	t.HTML = app.markdownHTMLCache["now.md"]
	app.render(w, r, 200, "simplePage.tmpl.html", &t)
}

func (app *application) about(w http.ResponseWriter, r *http.Request) {
	t := app.newTemplateData(r)
	t.HTML = app.markdownHTMLCache["about.md"]
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

type pageForm struct {
	Page        *models.Page
	FieldErrors map[string]string
}

func (app *application) adminNewPage(w http.ResponseWriter, r *http.Request) {
	t := app.newTemplateData(r)
	t.Form = pageForm{Page: t.Page}
	app.render(w, r, http.StatusOK, "admin.new_page.tmpl.html", &t)
}

func (app *application) adminPagesPage(w http.ResponseWriter, r *http.Request) {

}

func (app *application) pagesPost(w http.ResponseWriter, r *http.Request) {

}
