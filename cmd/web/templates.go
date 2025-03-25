package main

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/davidkuda/kudaai/internal/models"
)

type templateData struct {
	Title    string
	NavItems []NavItem
	RootPath string
	HTML     template.HTML
	Songs    models.Songs
	Song     *models.Song
	LoggedIn bool
	HideNav  bool
}

func (app *application) newTemplateData() templateData {
	return templateData{
		NavItems: app.navItems,
	}
}

func (app *application) render(w http.ResponseWriter, r *http.Request, status int, page string, data *templateData) {
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("couldn't find template \"%s\" in app.templateCache", page)
		app.serverError(w, r, err)
		return
	}

	w.WriteHeader(status)

	err := ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		errMsg := fmt.Errorf("error executing templates: %s", err.Error())
		app.serverError(w, r, errMsg)
	}
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob("./ui/html/pages/*.tmpl.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		files := []string{
			"./ui/html/pages/base.tmpl.html",
			"./ui/html/partials/nav.tmpl.html",
			page,
		}

		t, err := template.ParseFiles(files...)
		if err != nil {
			return nil, fmt.Errorf("Error parsing template files: %s", err.Error())
		}

		cache[name] = t
	}

	return cache, nil
}
