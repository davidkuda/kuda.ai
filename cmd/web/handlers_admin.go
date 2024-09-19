package main

import (
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/pascaldekloe/jwt"
)

func (app *application) adminLogin(w http.ResponseWriter, r *http.Request) {
	tmplFiles := []string{
		"./ui/html/pages/base.tmpl.html",
		"./ui/html/partials/nav.tmpl.html",
		"./ui/html/pages/admin/login.tmpl.html",
	}

	t, err := template.ParseFiles(tmplFiles...)
	if err != nil {
		log.Printf("Error parsing template files: %s", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = t.ExecuteTemplate(w, "base", nil)
	if err != nil {
		log.Printf("Error executing templates: %s", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (app *application) adminLoginPost(w http.ResponseWriter, r *http.Request) {
	type userLoginForm struct {
		email    string
		password string
	}
	err := r.ParseForm()
	if err != nil {
		log.Printf("Failed parsing form: %v", err)
		return
	}
	form := userLoginForm{
		email:    r.PostForm.Get("email"),
		password: r.PostForm.Get("password"),
	}

	err = app.users.Authenticate(form.email, form.password)
	if err != nil {
		log.Printf("error authenticating user: %v\n", err)
		return
	}

	var claims jwt.Claims
	claims.Subject = form.email
	claims.Issued = jwt.NewNumericTime(time.Now())
	claims.NotBefore = jwt.NewNumericTime(time.Now())
	claims.Expires = jwt.NewNumericTime(time.Now().Add(24 * time.Hour))
	claims.Issuer = "kuda.ai"
	claims.Audiences = []string{"kuda.ai"}

	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    "jwt",
		Domain:   "lyricsapi.kuda.ai",
		Expires:  time.Now().Add(24 * time.Hour),
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	})
}

func (app *application) adminNewSong(w http.ResponseWriter, r *http.Request) {

	tmplFiles := []string{
		"./ui/html/pages/base.tmpl.html",
		"./ui/html/partials/nav.tmpl.html",
		"./ui/html/pages/admin/new_song.tmpl.html",
	}

	t, err := template.ParseFiles(tmplFiles...)
	if err != nil {
		log.Printf("Error parsing home.tmpl.html: %s", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = t.ExecuteTemplate(w, "base", nil)
	if err != nil {
		log.Printf("Error executing home.tmpl.html: %s", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
