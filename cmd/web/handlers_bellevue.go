package main

import (
	"net/http"
)

func (app *application) adminBellevueActivity(w http.ResponseWriter, r *http.Request) {
	t := app.newTemplateData(r)
	t.Title = "New Bellevue Activity"
	app.render(w, r, http.StatusOK, "admin.new_bellevue_activity.tmpl.html", &t)
}
