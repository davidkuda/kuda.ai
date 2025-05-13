package main

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/davidkuda/kudaai/internal/models"
)

type templateData struct {
	Title         string
	NavItems      []NavItem
	Path          string
	RootPath      string
	HTML          template.HTML
	Songs         models.Songs
	Song          *models.Song
	TILs          models.TILs
	TIL           *models.TIL
	Pages         models.Pages
	Page          *models.Page
	Blogs         models.Blogs
	Blog          *models.Blog
	Form          any
	ShowUpdatedAt bool
	LoggedIn      bool
	HideNav       bool
	Sidebars      bool
	HighlightJS   bool
}

func (app *application) newTemplateData(r *http.Request) templateData {
	var err error

	var isAuthenticated bool
	err = app.checkJWTCookie(r)
	if err == nil {
		isAuthenticated = true
	}

	var rootPath, title string
	i := 1
	for i < len(r.URL.Path) && r.URL.Path[i] != '/' {
		i++
	}
	rootPath = r.URL.Path[0:i]
	title = strings.ToTitle(r.URL.Path[1:i])

	return templateData{
		NavItems: app.navItems,
		LoggedIn: isAuthenticated,
		Title:    title,
		RootPath: rootPath,
		Path:     r.URL.Path,
		Song:     &models.Song{},
		TIL:      &models.TIL{},
		Page:     &models.Page{},
		Blog:     &models.Blog{},
		Sidebars: true,
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

	funcs := template.FuncMap{
		"formatDate": formatDate,
	}

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

		tmpl := template.New("base").Funcs(funcs)
		t, err := tmpl.ParseFiles(files...)
		if err != nil {
			return nil, fmt.Errorf("Error parsing template files: %s", err.Error())
		}

		cache[name] = t
	}

	return cache, nil
}

func formatDate(t time.Time) string {
	return t.Format("January 2, 2006")
}
