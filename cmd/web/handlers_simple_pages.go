package main

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/davidkuda/kudaai/internal/models"
)

// GET /
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Started-Working-On", "April-2024")
	t := app.newTemplateData(r)
	page, err := app.models.Pages.Get("home")
	if err != nil {
		app.serverError(w, r, fmt.Errorf("app.models.pages.Get(\"home\"): %v", err))
		return
	}
	t.HTML = page.HTMLContent
	app.render(w, r, 200, "simplePage.tmpl.html", &t)
}

// GET /now
func (app *application) now(w http.ResponseWriter, r *http.Request) {
	t := app.newTemplateData(r)
	page, err := app.models.Pages.Get("now")
	if err != nil {
		app.serverError(w, r, fmt.Errorf("app.models.pages.Get(\"now\"): %v", err))
		return
	}
	t.HTML = page.HTMLContent
	app.render(w, r, 200, "simplePage.tmpl.html", &t)
}

// GET /about
func (app *application) about(w http.ResponseWriter, r *http.Request) {
	t := app.newTemplateData(r)
	page, err := app.models.Pages.Get("about")
	if err != nil {
		app.serverError(w, r, fmt.Errorf("app.models.pages.Get(\"about\"): %v", err))
		return
	}
	t.HTML = page.HTMLContent
	app.render(w, r, 200, "simplePage.tmpl.html", &t)
}

// GET /bookshelf
func (app *application) bookshelf(w http.ResponseWriter, r *http.Request) {
	t := app.newTemplateData(r)
	page, err := app.models.Pages.Get("bookshelf")
	if err != nil {
		app.serverError(w, r, fmt.Errorf("app.models.pages.Get(\"bookshelf\"): %v", err))
		return
	}
	t.HTML = page.HTMLContent
	app.render(w, r, 200, "simplePage.tmpl.html", &t)
}

type pageForm struct {
	Page        *models.Page
	FieldErrors map[string]string
}

// GET /admin/new-page
func (app *application) adminNewPage(w http.ResponseWriter, r *http.Request) {
	t := app.newTemplateData(r)
	t.Form = pageForm{Page: t.Page}
	app.render(w, r, http.StatusOK, "admin.new_page.tmpl.html", &t)
}

// GET /admin/pages/:page
func (app *application) adminPagesPage(w http.ResponseWriter, r *http.Request) {
	var err error

	name := r.PathValue("page")

	page, err := app.models.Pages.Get(name)
	if err != nil {
		log.Printf("app.pages.GetByPath(%v): %v", name, err)
		http.NotFound(w, r)
		return
	}

	t := app.newTemplateData(r)
	t.Page = page
	t.Form = pageForm{}
	app.render(w, r, http.StatusOK, "admin.new_page.tmpl.html", &t)
}

// POST /admin/pages
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
			Name:    f.Get("page-path"),
			Title:   f.Get("page-title"),
			Content: strings.ReplaceAll(f.Get("page-content"), "\r\n", "\n"),
		},
		FieldErrors: map[string]string{},
	}

	// TODO: This is used several times. make a function.
	var rxPat = regexp.MustCompile(`[^a-z\-]`)
	if rxPat.MatchString(form.Page.Name) {
		form.FieldErrors["pathfmt"] = "path/name may only contain lowercase characters and hyphens"
	}

	if len(form.FieldErrors) > 0 {
		t := app.newTemplateData(r)
		t.Form = form
		t.Page = form.Page
		app.render(w, r, http.StatusUnprocessableEntity, "admin.new_page.tmpl.html", &t)
		return
	}

	err = app.models.Pages.Insert(form.Page)
	if err != nil {
		log.Printf("app.pages.Insert(form.Page): %v\n", err)
		return
	}
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
	return

}
