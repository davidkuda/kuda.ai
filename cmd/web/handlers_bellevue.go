package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/davidkuda/kudaai/internal/models"
)

// GET /admin/new-bellevue-activity
func (app *application) adminNewBellevueActivity(w http.ResponseWriter, r *http.Request) {
	t := app.newTemplateData(r)
	t.Title = "New Bellevue Activity"
	t.BellevueOfferings = models.NewBellevueOfferings()
	app.render(w, r, http.StatusOK, "admin.new_bellevue_activity.tmpl.html", &t)
}

// GET /bellevue-activities
func (app *application) bellevueActivities(w http.ResponseWriter, r *http.Request) {
	t := app.newTemplateData(r)
	bas, err := app.models.BellevueActivities.GetAllByUser(t.UserID)
	if err != nil {
		log.Println(fmt.Errorf("failed reading bellevue activities: %v", err))
	}
	t.BellevueActivityOverview.BellevueActivities = bas
	t.BellevueActivityOverview.CalculateTotalPrice()
	app.render(w, r, http.StatusOK, "bellevue_activities.tmpl.html", &t)
}

type bellevueActivityForm struct {
	BellevueActivity *models.BellevueActivity
	FieldErrors      map[string]string
}

// POST /admin/new-bellevue-activity
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

	// although the form has client side validation of integers by using
	// <input type="number">, a malicious actor could still place a POST
	// request not via the web form.
	breakfasts, err := strconv.Atoi(f.Get("bellevue-activity-breakfasts"))
	if err != nil {
		log.Printf("form.breakfasts: stconv.Atoi: someone wants to write non-integers: value: %v, err: %v", f.Get("bellevue-activity-breakfasts"), err)
		app.renderClientError(w, r, http.StatusUnprocessableEntity)
		return
	}
	lunches, err := strconv.Atoi(f.Get("bellevue-activity-lunches"))
	if err != nil {
		log.Printf("form.lunches: stconv.Atoi: someone wants to write non-integers: value: %v, err: %v", f.Get("bellevue-activity-lunches"), err)
		app.renderClientError(w, r, http.StatusUnprocessableEntity)
		return
	}
	dinners, err := strconv.Atoi(f.Get("bellevue-activity-dinners"))
	if err != nil {
		log.Printf("form.dinners: stconv.Atoi: someone wants to write non-integers: value: %v, err: %v", f.Get("bellevue-activity-dinners"), err)
		app.renderClientError(w, r, http.StatusUnprocessableEntity)
		return
	}
	coffees, err := strconv.Atoi(f.Get("bellevue-activity-coffees"))
	if err != nil {
		log.Printf("form.coffees: stconv.Atoi: someone wants to write non-integers: value: %v, err: %v", f.Get("bellevue-activity-coffees"), err)
		app.renderClientError(w, r, http.StatusUnprocessableEntity)
		return
	}
	saunas, err := strconv.Atoi(f.Get("bellevue-activity-saunas"))
	if err != nil {
		log.Printf("form.saunas: stconv.Atoi: someone wants to write non-integers: value: %v, err: %v", f.Get("bellevue-activity-saunas"), err)
		app.renderClientError(w, r, http.StatusUnprocessableEntity)
		return
	}
	lectures, err := strconv.Atoi(f.Get("bellevue-activity-lectures"))
	if err != nil {
		log.Printf("form.lectures: stconv.Atoi: someone wants to write non-integers: value: %v, err: %v", f.Get("bellevue-activity-lectures"), err)
		app.renderClientError(w, r, http.StatusUnprocessableEntity)
		return
	}

	var snacksCHF int
	snackCHFString := f.Get("bellevue-activity-snacks")
	if len(snackCHFString) > 0 {
		priceFloat, err := strconv.ParseFloat(snackCHFString, 64)
		if err != nil {
			log.Printf("failed parsing string \"%s\" to float:", snackCHFString)
			app.renderClientError(w, r, http.StatusInternalServerError)
			return
		}
		snacksCHF = int(math.Round(priceFloat * 100))
	}

	form := bellevueActivityForm{
		BellevueActivity: &models.BellevueActivity{
			Date:       date,
			Breakfasts: breakfasts,
			Lunches:    lunches,
			Dinners:    dinners,
			Coffees:    coffees,
			Saunas:     saunas,
			Lectures:   lectures,
			SnacksCHF:  snacksCHF,
			Comment:    f.Get("bellevue-activity-comment"),
		},
		FieldErrors: map[string]string{},
	}

	// TODO: FieldErrors?
	//       - [ ] if all counts are 0
	//       - [ ] if a count is negative
	if len(form.FieldErrors) > 0 {
		log.Println("field errors")
		return
	}

	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		log.Println("post /bellevue-activity: could not get userID from request.Context")
		// TODO: return 503
		return
	}
	form.BellevueActivity.UserID = userID

	err = app.models.BellevueActivities.Insert(form.BellevueActivity)
	if err != nil {
		log.Printf("app.bellevueActivities.Insert(): %v\n", err)
		// TODO: send some notification to the UI (failed submission)
		return
	}

	// TODO: send some notification (Toast) to the UI (successfully submitted)
	http.Redirect(w, r, "/bellevue-activities", http.StatusSeeOther)
	return
}
