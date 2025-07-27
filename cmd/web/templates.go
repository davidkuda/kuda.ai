package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/davidkuda/kudaai/internal/models"
)

type templateData struct {
	LoggedIn                 bool
	UserID                   int
	UserEmail                string
	IsAdmin                  bool
	Title                    string
	NavItems                 []NavItem
	Path                     string
	RootPath                 string
	HTML                     template.HTML
	Songs                    models.Songs
	Song                     *models.Song
	TILs                     models.TILs
	TIL                      *models.TIL
	Pages                    models.Pages
	Page                     *models.Page
	Blogs                    models.Blogs
	Blog                     *models.Blog
	BellevueActivityOverview models.BellevueActivityOverview
	BellevueActivity         *models.BellevueActivity
	BellevueOfferings        models.BellevueOfferings
	Form                     any
	ShowUpdatedAt            bool
	HideNav                  bool
	Sidebars                 bool
	HighlightJS              bool
	Error                    Error
}

type Error struct {
	HTTPStatusCode int
	HTTPStatusText string
	Method         string
	Path           string
}

func newError(r *http.Request, errorCode int) Error {
	return Error{
		HTTPStatusCode: errorCode,
		HTTPStatusText: http.StatusText(errorCode),
		Method:         r.Method,
		Path:           r.URL.Path,
	}
}

func (app *application) newTemplateData(r *http.Request) templateData {

	isAuthenticated, ok := r.Context().Value("isAuthenticated").(bool)
	if !ok {
		isAuthenticated = false
	}

	var userID int
	var userEmail string
	var isAdmin bool

	if isAuthenticated {
		userID, ok = r.Context().Value("userID").(int)
		if !ok {
			// TODO: what to do with this check?
			log.Println("newTemplateData: could not get userID from request.Context")
		}

		userEmail, ok = r.Context().Value("userEmail").(string)
		if !ok {
			// TODO: what to do with this check?
			log.Println("newTemplateData: could not get userEmail from request.Context")
		}

		isAdmin, ok = r.Context().Value("isAdmin").(bool)
		if !ok {
			isAdmin = false
		}
	}

	var rootPath, title string
	i := 1
	for i < len(r.URL.Path) && r.URL.Path[i] != '/' {
		i++
	}
	rootPath = r.URL.Path[0:i]
	title = strings.ToTitle(r.URL.Path[1:i])

	// TODO: using empty structs with a pointer seems so wrong here. How to fix it?
	// problem is that the templates will error on render.
	return templateData{
		LoggedIn:         isAuthenticated,
		UserID:           userID,
		UserEmail:        userEmail,
		IsAdmin:          isAdmin,
		NavItems:         app.navItems,
		Title:            title,
		RootPath:         rootPath,
		Path:             r.URL.Path,
		Song:             &models.Song{},
		TIL:              &models.TIL{},
		Page:             &models.Page{},
		Blog:             &models.Blog{},
		BellevueActivity: models.NewBellevueActivity(),
		Sidebars:         true,
	}
}

func (app *application) render(w http.ResponseWriter, r *http.Request, status int, page string, data *templateData) {
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("couldn't find template \"%s\" in app.templateCache", page)
		app.serverError(w, r, err)
		return
	}

	buf := bytes.Buffer{}

	err := ts.ExecuteTemplate(&buf, "base", data)
	if err != nil {
		errMsg := fmt.Errorf("error executing template %s: %s", page, err.Error())
		app.serverError(w, r, errMsg)
		return
	}

	w.WriteHeader(status)

	buf.WriteTo(w)
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	funcs := template.FuncMap{
		"formatDate":          formatDate,
		"formatDateFormInput": formatDateFormInput,
		"fmtDateNiceRead":     formatDateNiceRead,
		"fmtCHF":              formatCurrency,
	}

	pages, err := filepath.Glob("./ui/html/pages/*.tmpl.html")
	if err != nil {
		return nil, fmt.Errorf("failed filepath.Glob for pages: %v", err)
	}

	partials, err := filepath.Glob("./ui/html/partials/*.tmpl.html")
	if err != nil {
		return nil, fmt.Errorf("failed filepath.Glob for partials: %v", err)
	}

	for _, page := range pages {
		name := filepath.Base(page)

		N := 1 + len(partials) + 1
		files := make([]string, N)
		files[0] = "./ui/html/pages/base.tmpl.html"
		for i, partial := range partials {
			files[i+1] = partial
		}
		files[N-1] = page

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

func formatDateFormInput(t time.Time) string {
	return t.Format("2006-01-02")
}

func formatDateNiceRead(t time.Time) string {
	return t.Format("Mon 2.01.2006")
}

// formatCurrency converts an integer (in Rappen) to a currency string like "22.50 CHF".
func formatCurrency(value int) string {
	return fmt.Sprintf("%.2f CHF", float64(value)/100)
}
