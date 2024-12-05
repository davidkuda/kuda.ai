package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/davidkuda/kudaai/internal/models"
)

type templateData struct {
	Songs           models.Songs
	Song            models.Song
	CurrentRootPath string
}

func (app *application) render(w http.ResponseWriter, r *http.Request, status int, page string, data *templateData) {
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
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
