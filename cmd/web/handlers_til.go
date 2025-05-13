package main

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/davidkuda/kudaai/internal/models"
)

type tilForm struct {
	TIL         *models.TIL
	FieldErrors map[string]string
}

func (app *application) todayILearned(w http.ResponseWriter, r *http.Request) {
	t := app.newTemplateData(r)
	t.Title = "Today I Learned"
	t.RootPath = "/today-i-learned"
	tils, err := app.til.GetAll()
	if err != nil {
		// TODO: how to handle gracefully?
		log.Println(err)
	}
	t.TILs = tils
	t.Sidebars = false
	app.render(w, r, 200, "tils.tmpl.html", &t)
}

func (app *application) todayILearnedPath(w http.ResponseWriter, r *http.Request) {

	path := r.PathValue("path")

	til, err := app.til.GetBy(path)
	if err != nil {
		log.Printf("app.til.GetBy(%v): %v", path, err)
		// TODO: Show a nice 404 page.
		http.NotFound(w, r)
		return
	}

	t := app.newTemplateData(r)
	t.Title = "Today I Learned"
	t.RootPath = "/today-i-learned"
	t.TIL = til
	t.HighlightJS = true

	if !isSameDay(t.TIL.CreatedAt, t.TIL.UpdatedAt) {
		t.ShowUpdatedAt = true
	}

	app.render(w, r, 200, "tils.til.tmpl.html", &t)
}

func (app *application) adminNewTIL(w http.ResponseWriter, r *http.Request) {
	t := app.newTemplateData(r)
	t.Form = tilForm{TIL: &models.TIL{}}
	app.render(w, r, http.StatusOK, "admin.new_til.tmpl.html", &t)
}

func (app *application) adminTILSTIL(w http.ResponseWriter, r *http.Request) {
	var err error

	path := r.PathValue("path")
	til := &models.TIL{}

	til, err = app.til.GetBy(path)
	if err != nil {
		log.Printf("app.til.GetBy(%v): %v", path, err)
		// TODO: Show a nice 404 page.
		http.NotFound(w, r)
		return
	}

	t := app.newTemplateData(r)
	t.TIL = til
	t.Form = tilForm{TIL: &models.TIL{}}
	app.render(w, r, http.StatusOK, "admin.new_til.tmpl.html", &t)

}

func (app *application) tilPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Printf("Failed parsing form: %v", err)
		// TODO: send status 400 Bad Request to the client
		return
	}

	f := r.PostForm

	var id int
	if f.Get("til-id") != "" {
		id, err = strconv.Atoi(f.Get("til-id"))
		if err != nil {
			log.Printf("failed converting ascii to int: %v", err)
			// TODO: send status 400 Bad Request to the client
			return
		}
	}
	var isUpdate bool
	if id > 0 {
		isUpdate = true
	}

	form := tilForm{
		TIL: &models.TIL{
			ID:       id,
			Path:     f.Get("til-path"),
			Title:    f.Get("til-title"),
			Category: f.Get("til-category"),
			Summary:  f.Get("til-summary"),
			Text:     strings.ReplaceAll(f.Get("til-text"), "\r\n", "\n"),
		},
		FieldErrors: map[string]string{},
	}

	// regex for valid URL path; TIL.Path will be used in the URL.
	// Therefore, it should only contain letters and hyphens.
	var rxPat = regexp.MustCompile(`[^a-z\-]`)
	if rxPat.MatchString(form.TIL.Path) {
		form.FieldErrors["pathfmt"] = "id may only contain lowercase characters and hyphens"
	}

	if !isUpdate {
		uniq, err := app.til.PathIsUnique(form.TIL.Path)
		if err != nil {
			log.Printf("app.til.PathIsUnique(Path): %v", err)
			// TODO: send status 400 Bad Request to the client
			return
		}
		if !uniq {
			form.FieldErrors["pathuniq"] = "path is not unique"

		}
	}

	for i := range form.TIL.Summary {
		if form.TIL.Summary[i] == '\n' {
			form.FieldErrors["summary"] = "please write everything in one paragraph without new lines"
		}
	}

	if len(form.FieldErrors) > 0 {
		t := app.newTemplateData(r)
		t.Form = form
		t.TIL = form.TIL
		app.render(w, r, http.StatusUnprocessableEntity, "admin.new_til.tmpl.html", &t)
		return
	}

	if id == 0 {
		err = app.til.Insert(form.TIL)
		if err != nil {
			log.Printf("app.til.Insert(): %v\n", err)
			// TODO: send some notification to the UI
		}
	} else {
		err = app.til.UpdateExisting(form.TIL)
		if err != nil {
			log.Printf("app.til.UpdateExisting(): %v\n", err)
			// TODO: send some notification to the UI
		}
	}
	http.Redirect(w, r, fmt.Sprintf("/today-i-learned/%v", form.TIL.Path), http.StatusSeeOther)
	return
}
