package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/pascaldekloe/jwt"
)

func (app *application) admin(w http.ResponseWriter, r *http.Request) {
	t := app.newTemplateData(r)
	app.render(w, r, 200, "admin.tmpl.html", &t)
}

func (app *application) adminLogin(w http.ResponseWriter, r *http.Request) {
	if app.isAuthenticated(r) {
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
		return
	}
	t := app.newTemplateData(r)
	t = templateData{
		Title:    "Login",
		RootPath: "/admin",
	}

	app.render(w, r, 200, "admin.login.tmpl.html", &t)
}

func (app *application) adminLoginPost(w http.ResponseWriter, r *http.Request) {
	type userLoginForm struct {
		email    string
		password string
	}
	err := r.ParseForm()
	if err != nil {
		log.Printf("Failed parsing form: %v", err)
		http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
		return
	}
	form := userLoginForm{
		email:    r.PostForm.Get("email"),
		password: r.PostForm.Get("password"),
	}

	err = app.users.Authenticate(form.email, form.password)
	if err != nil {
		log.Printf("error authenticating user: %v\n", err)
		http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
		return
	}

	var claims jwt.Claims
	claims.Subject = form.email
	claims.Issued = jwt.NewNumericTime(time.Now())
	claims.NotBefore = jwt.NewNumericTime(time.Now())
	claims.Expires = jwt.NewNumericTime(time.Now().Add(24 * time.Hour))
	claims.Issuer = "kuda.ai"
	claims.Audiences = []string{"kuda.ai"}

	jwtBytes, err := claims.HMACSign(jwt.HS256, []byte(app.JWT.Secret))
	if err != nil {
		log.Printf("error signing jwt: %v\n", err)
		http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    string(jwtBytes),
		Domain:   app.JWT.CookieDomain,
		Expires:  time.Now().Add(24 * time.Hour),
		Secure:   true,
		HttpOnly: true,
		// SameSite: http.SameSiteNoneMode,
	})

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func (app *application) adminLogoutPost(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    "",
		Domain:   app.JWT.CookieDomain,
		Expires:  time.Now(),
		Secure:   true,
		HttpOnly: true,
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) adminNewSong(w http.ResponseWriter, r *http.Request) {

	err := app.checkJWTCookie(r)
	if err != nil {
		log.Printf("could not authenticate client: %v", err)
		http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
	} else {
		log.Println("Could authenticate client :)")
	}

	tmplFiles := []string{
		"./ui/html/pages/base.tmpl.html",
		"./ui/html/partials/nav.tmpl.html",
		"./ui/html/pages/admin.new_song.tmpl.html",
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

func (app *application) isAuthenticated(r *http.Request) bool {
	err := app.checkJWTCookie(r)
	if err == nil {
		return true
	} else {
		return false
	}
}

func (app *application) checkJWTCookie(r *http.Request) error {
	token, err := r.Cookie("session")
	if err != nil {
		return fmt.Errorf("couldn't find cookie: %v", err)
	}

	claims, err := jwt.HMACCheck([]byte(token.Value), []byte(app.JWT.Secret))
	if err != nil {
		return fmt.Errorf("detected invalid signature in jwtCookie: %v", err)
	}

	if !claims.Valid(time.Now()) {
		return fmt.Errorf("token no longer valid")
	}

	if claims.Issuer != "kuda.ai" {
		return fmt.Errorf("token has invalid issuer: %v", err)
	}

	if !claims.AcceptAudience("kuda.ai") {
		return fmt.Errorf("token is not in accepted audience: %v", err)
	}

	return nil
}
