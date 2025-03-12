package main

import (
	"net/http"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Server", "Go")
	w.Header().Add("Started-Working-On", "April-2024")

	t := templateData{
		Title:    "kuda.ai",
		HideNav:  true,
	}

	app.render(w, r, 200, "home.tmpl.html", &t)
}

func (app *application) about(w http.ResponseWriter, r *http.Request) {
	t := templateData{
		Title:    "About",
		RootPath: "/about",
		HTML:     app.markdownHTMLCache["about.md"],
	}
	app.render(w, r, 200, "simplePage.tmpl.html", &t)
}

func (app *application) blog(w http.ResponseWriter, r *http.Request) {
	t := templateData{
		Title:    "Blog",
		RootPath: "/blog",
		HTML:     app.markdownHTMLCache["blog.md"],
	}
	app.render(w, r, 200, "simplePage.tmpl.html", &t)
}

func (app *application) bookshelf(w http.ResponseWriter, r *http.Request) {
	t := templateData{
		Title:    "Bookshelf",
		RootPath: "/bookshelf",
		HTML:     app.markdownHTMLCache["bookshelf.md"],
	}
	app.render(w, r, 200, "simplePage.tmpl.html", &t)
}

func (app *application) cv(w http.ResponseWriter, r *http.Request) {
	// TODO: wouldn't it be nice to generate a beautiful PDF from this site?
	t := templateData{
		Title:    "CV",
		RootPath: "/cv",
		HTML:     app.markdownHTMLCache["cv.md"],
	}
	app.render(w, r, 200, "simplePage.tmpl.html", &t)
}

func (app *application) til(w http.ResponseWriter, r *http.Request) {
	t := templateData{
		Title:    "Today I Learned",
		RootPath: "/today-i-learned",
		HTML:     app.markdownHTMLCache["til.md"],
	}
	app.render(w, r, 200, "simplePage.tmpl.html", &t)
}
