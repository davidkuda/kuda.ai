package main

import (
	"net/http"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Server", "Go")
	w.Header().Add("Started-Working-On", "April-2024")

	t := app.newTemplateData()
	t = templateData{
		NavItems: t.NavItems,
		Title:    "kuda.ai",
		HideNav:  true,
	}

	app.render(w, r, 200, "home.tmpl.html", &t)
}

func (app *application) about(w http.ResponseWriter, r *http.Request) {
	t := app.newTemplateData()
	t = templateData{
		NavItems: t.NavItems,
		Title:    "About",
		RootPath: "/about",
		HTML:     app.markdownHTMLCache["about.md"],
	}
	app.render(w, r, 200, "simplePage.tmpl.html", &t)
}

func (app *application) blog(w http.ResponseWriter, r *http.Request) {
	t := app.newTemplateData()
	t = templateData{
		NavItems: t.NavItems,
		Title:    "Blog",
		RootPath: "/blog",
		HTML:     app.markdownHTMLCache["blog.md"],
	}
	app.render(w, r, 200, "simplePage.tmpl.html", &t)
}

func (app *application) bookshelf(w http.ResponseWriter, r *http.Request) {
	t := app.newTemplateData()
	t = templateData{
		NavItems: t.NavItems,
		Title:    "Bookshelf",
		RootPath: "/bookshelf",
		HTML:     app.markdownHTMLCache["bookshelf.md"],
	}
	app.render(w, r, 200, "simplePage.tmpl.html", &t)
}

//# TODO: implement? or just link to LinkedIn?
func (app *application) cv(w http.ResponseWriter, r *http.Request) {
	// TODO: wouldn't it be nice to generate a beautiful PDF from this site?
	t := templateData{
		Title:    "CV",
		RootPath: "/cv",
		HTML:     app.markdownHTMLCache["cv.md"],
	}
	app.render(w, r, 200, "simplePage.tmpl.html", &t)
}

func (app *application) todayILearned(w http.ResponseWriter, r *http.Request) {
	t := app.newTemplateData()
	t = templateData{
		NavItems: t.NavItems,
		Title:    "Today I Learned",
		RootPath: "/today-i-learned",
		HTML:     app.markdownHTMLCache["til.md"],
	}
	app.render(w, r, 200, "simplePage.tmpl.html", &t)
}
