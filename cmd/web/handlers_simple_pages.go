package main

import (
	"log"
	"net/http"
	"regexp"

	"github.com/davidkuda/kudaai/internal/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
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
	var err error

	path := r.PathValue("page")

	page, err := app.pages.GetByPath(path)
	if err != nil {
		log.Printf("app.pages.GetByPath(%v): %v", path, err)
		http.NotFound(w, r)
		return
	}

	t := app.newTemplateData(r)
	t.Page = page
	t.Form = pageForm{}
	app.render(w, r, http.StatusOK, "admin.new_page.tmpl.html", &t)
}

func (app *application) pagesPost(w http.ResponseWriter, r *http.Request) {
	var err error

	err = r.ParseForm()
	if err != nil {
		log.Printf("Failed parsing form: %v", err)
		return
	}

	f := r.PostForm

	form := pageForm{
		Page: &models.Page{
			Path:    f.Get("page-path"),
			Title:   f.Get("page-title"),
			Content: f.Get("page-content"),
		},
		FieldErrors: map[string]string{},
	}

	// TODO: This is used several times. make a function.
	var rxPat = regexp.MustCompile(`[^a-z\-]`)
	if rxPat.MatchString(form.Page.Path) {
		form.FieldErrors["pathfmt"] = "path may only contain lowercase characters and hyphens"
	}

	if len(form.FieldErrors) > 0 {
		t := app.newTemplateData(r)
		t.Form = form
		t.Page = form.Page
		app.render(w, r, http.StatusUnprocessableEntity, "admin.new_page.tmpl.html", &t)
		return
	}

	err = app.pages.Insert(form.Page)
	if err != nil {
		log.Printf("app.pages.Insert(form.Page): %v\n", err)
		return
	}
	http.Redirect(w, r, "/" + form.Page.Path, http.StatusSeeOther)
	return

}
