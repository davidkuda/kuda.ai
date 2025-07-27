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
	t.Form = bellevueActivityForm{}
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
		Date:        date,
		Breakfasts:  breakfasts,
		Lunches:     lunches,
		Dinners:     dinners,
		Coffees:     coffees,
		Saunas:      saunas,
		Lectures:    lectures,
		SnacksCHF:   snacksCHF,
		Comment:     f.Get("bellevue-activity-comment"),
		FieldErrors: map[string]string{},
	}

	if form.hasNegativeNumbers() {
		form.FieldErrors["negatives"] = "you may not send negative numbers."
	}
	if form.hasOnlyZeroes() {
		form.FieldErrors["zeroes"] = "you have all 0 and therefore not any activity to upload."
	}
	if len(form.FieldErrors) > 0 {
		log.Println("field errors")
		data := app.newTemplateData(r)
		data.BellevueActivity = form.toModel()
		data.BellevueActivity.PopulateItems()
		data.BellevueOfferings = models.NewBellevueOfferings()
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "admin.new_bellevue_activity.tmpl.html", &data)
		return
	}

	b := form.toModel()

	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		log.Println("post /bellevue-activity: could not get userID from request.Context")
		app.serverError(w, r, err)
		return
	}
	b.UserID = userID

	err = app.models.BellevueActivities.Insert(b)
	if err != nil {
		log.Printf("app.bellevueActivities.Insert(): %v\n", err)
		app.serverError(w, r, err)
		return
	}

	// TODO: send some notification (Toast) to the UI (successfully submitted)
	http.Redirect(w, r, "/bellevue-activities", http.StatusSeeOther)
	return
}

type bellevueActivityForm struct {
	ID          int
	UserID      int
	Date        time.Time
	Breakfasts  int
	Lunches     int
	Dinners     int
	Coffees     int
	Saunas      int
	Lectures    int
	SnacksCHF   int
	Comment     string
	FieldErrors map[string]string
}

func (b *bellevueActivityForm) hasNegativeNumbers() bool {
	if b.Breakfasts < 0 {
		return true
	}
	if b.Lunches < 0 {
		return true
	}
	if b.Dinners < 0 {
		return true
	}
	if b.Coffees < 0 {
		return true
	}
	if b.Saunas < 0 {
		return true
	}
	if b.Lectures < 0 {
		return true
	}
	if b.SnacksCHF < 0 {
		return true
	}
	return false
}

func (b *bellevueActivityForm) hasOnlyZeroes() bool {
	if b.Breakfasts > 0 {
		return false
	}
	if b.Lunches > 0 {
		return false
	}
	if b.Dinners > 0 {
		return false
	}
	if b.Coffees > 0 {
		return false
	}
	if b.Saunas > 0 {
		return false
	}
	if b.Lectures > 0 {
		return false
	}
	if b.SnacksCHF > 0 {
		return false
	}
	return true
}

func (b *bellevueActivityForm) toModel() *models.BellevueActivity {
	return &models.BellevueActivity{
		ID:         b.ID,
		UserID:     b.UserID,
		Date:       b.Date,
		Breakfasts: b.Breakfasts,
		Lunches:    b.Lunches,
		Dinners:    b.Dinners,
		Coffees:    b.Coffees,
		Saunas:     b.Saunas,
		SnacksCHF:  b.SnacksCHF,
		Lectures:   b.Lectures,
		Comment:    b.Comment,
	}
}
