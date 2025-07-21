package main

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/davidkuda/kudaai/internal/models"
)

func (app *application) adminNewBellevueActivity(w http.ResponseWriter, r *http.Request) {
	t := app.newTemplateData(r)
	t.Title = "New Bellevue Activity"
	app.render(w, r, http.StatusOK, "admin.new_bellevue_activity.tmpl.html", &t)
}

func (app *application) bellevueActivitiesGet(w http.ResponseWriter, r *http.Request) {

}

type bellevueActivityForm struct {
	BellevueActivity *models.BellevueActivity
	FieldErrors      map[string]string
}

func (app *application) bellevueActivityPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Printf("Failed parsing form: %v", err)
		// TODO: send status 400 Bad Request to the client
		return
	}

	f := r.PostForm

	// TODO: handle errors, and send feedback to client.
	dateStr := f.Get("bellevue-activity-date")
	date, _ := time.Parse("2006-01-02", dateStr)

	breakfasts, _ := strconv.Atoi(f.Get("bellevue-activity-breakfasts"))
	lunchDinners, _ := strconv.Atoi(f.Get("bellevue-activity-lunchdinners"))
	coffees, _ := strconv.Atoi(f.Get("bellevue-activity-coffees"))
	saunas, _ := strconv.Atoi(f.Get("bellevue-activity-saunas"))
	lectures, _ := strconv.Atoi(f.Get("bellevue-activity-lectures"))

	form := bellevueActivityForm{
		BellevueActivity: &models.BellevueActivity{
			Date:         date,
			Breakfasts:   breakfasts,
			LunchDinners: lunchDinners,
			Coffees:      coffees,
			Saunas:       saunas,
			Lectures:     lectures,
			Comment:      f.Get("bellevue-activity-comment"),
		},
		FieldErrors: map[string]string{},
	}

	// TODO: FieldErrors?
	//       - [ ] if all counts are 0
	if len(form.FieldErrors) > 0 {
		return
	}

	err = app.models.BellevueActivities.Insert(form.BellevueActivity)
	if err != nil {
		log.Printf("app.bellevueActivities.Insert(): %v\n", err)
		// TODO: send some notification to the UI (failed submission)
		return
	}

	// TODO: send some notification to the UI (successfully submitted)
	http.Redirect(w, r, "/admin/new-bellevue-activity", http.StatusSeeOther)
	return
}
