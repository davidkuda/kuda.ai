package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/russross/blackfriday/v2"
)

func (app *application) renderClientError(w http.ResponseWriter, r *http.Request, errorCode int) {
	ts, ok := app.templateCache["error.tmpl.html"]
	if !ok {
		err := fmt.Errorf("couldn't find template error.tmpl.html in app.templateCache")
		app.serverError(w, r, err)
		return
	}
	
	buf := bytes.Buffer{}

	data := app.newTemplateData(r)
	data.Error = newError(r, errorCode)

	err := ts.ExecuteTemplate(&buf, "base", data)
	if err != nil {
		errMsg := fmt.Errorf("error executing templates: %s", err.Error())
		app.serverError(w, r, errMsg)
	}

	w.WriteHeader(errorCode)
	buf.WriteTo(w)
}

// The serverError helper writes a log entry at Error level (including the request
// method and URI as attributes), then sends a generic 500 Internal Server Error
// response to the user.
func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
	)

	// TODO: implement a logger
	// app.logger.Error(err.Error(), "method", method, "uri", uri)

	log.Printf("%v (method: %v, uri: %v)\n", err.Error(), method, uri)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
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

func newMarkdownHTMLCache() (map[string]template.HTML, error) {
	pages := make(map[string]template.HTML)

	files := []string{
		"home.md",
		"now.md",
		"about.md",
		"blog.md",
		"bookshelf.md",
	}

	for _, file := range files {
		md, err := os.ReadFile("./data/pages/" + file)
		if err != nil {
			return nil, fmt.Errorf("could not read file: %v", err)
		}

		htmlBytes := blackfriday.Run(md)
		content := template.HTML(htmlBytes)

		pages[file] = content
	}

	return pages, nil
}

func isSameDay(a, b time.Time) bool {
	return a.Year() == b.Year() &&
		a.Month() == b.Month() &&
		a.Day() == b.Day()
}
