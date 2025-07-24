package main

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/davidkuda/kudaai/internal/models"
)

type blogForm struct {
	Blog        *models.Blog
	FieldErrors map[string]string
}

func (app *application) blog(w http.ResponseWriter, r *http.Request) {
	t := app.newTemplateData(r)
	// TODO: Here, we should probably not get HTML into the struct, because a lot to fetch and keep in memory and not necessary for the overview.
	blogs, err := app.blogs.GetAll()
	if err != nil {
		// TODO: how to handle gracefully?
		log.Println(err)
	}
	t.Blogs = blogs
	app.render(w, r, 200, "blog.tmpl.html", &t)
}

func (app *application) blogPath(w http.ResponseWriter, r *http.Request) {

	path := r.PathValue("path")

	blog, err := app.blogs.GetByPath(path)
	if err != nil {
		log.Printf("app.blogs.GetByPath(%v): %v", path, err)
		// TODO: Show a nice 404 page.
		http.NotFound(w, r)
		return
	}

	t := app.newTemplateData(r)
	t.Blog = blog
	t.HighlightJS = true

	if !isSameDay(t.Blog.CreatedAt, t.Blog.UpdatedAt) {
		t.ShowUpdatedAt = true
	}

	app.render(w, r, 200, "blog.entry.tmpl.html", &t)

}

func (app *application) adminNewBlog(w http.ResponseWriter, r *http.Request) {
	t := app.newTemplateData(r)
	t.Form = blogForm{Blog: t.Blog}
	app.render(w, r, http.StatusOK, "admin.new_blog.tmpl.html", &t)
}

func (app *application) adminBlogPath(w http.ResponseWriter, r *http.Request) {
	var err error

	path := r.PathValue("path")
	blog := &models.Blog{}

	blog, err = app.blogs.GetByPath(path)
	if err != nil {
		log.Printf("app.til.GetBy(%v): %v", path, err)
		// TODO: Show a nice 404 page.
		http.NotFound(w, r)
		return
	}

	t := app.newTemplateData(r)
	t.Blog = blog
	t.Form = blogForm{Blog: blog}
	app.render(w, r, http.StatusOK, "admin.new_blog.tmpl.html", &t)

}

func (app *application) blogPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Printf("Failed parsing form: %v", err)
		// TODO: send status 400 Bad Request to the client
		return
	}

	f := r.PostForm

	var id int
	if f.Get("blog-id") != "" {
		id, err = strconv.Atoi(f.Get("blog-id"))
		if err != nil {
			log.Printf("failed converting ascii to int: %v", err)
			// TODO: send status 400 Bad Request to the client
			return
		}
	}
	var isUpdate bool
	if id > 0 {
		isUpdate = true
	}

	form := blogForm{
		Blog: &models.Blog{
			ID:       id,
			Path:     f.Get("blog-path"),
			Title:    f.Get("blog-title"),
			Summary:  f.Get("blog-summary"),
			Content:     strings.ReplaceAll(f.Get("blog-text"), "\r\n", "\n"),
		},
		FieldErrors: map[string]string{},
	}

	// regex for valid URL path; TIL.Path will be used in the URL.
	// Therefore, it should only contain letters and hyphens.
	var rxPat = regexp.MustCompile(`[^a-z\-0-9]`)
	if rxPat.MatchString(form.Blog.Path) {
		form.FieldErrors["pathfmt"] = "id may only contain lowercase characters, numbers and hyphens"
	}

	// TODO: Abstract away into a Validator
	if !isUpdate {
		uniq, err := app.blogs.PathIsUnique(form.Blog.Path)
		if err != nil {
			log.Printf("app.til.PathIsUnique(Path): %v", err)
			// TODO: send status 400 Bad Request to the client
			return
		}
		if !uniq {
			form.FieldErrors["pathuniq"] = "path is not unique"

		}
	}

	for i := range form.Blog.Summary {
		if form.Blog.Summary[i] == '\n' {
			form.FieldErrors["summary"] = "please write everything in one paragraph without new lines"
		}
	}

	if len(form.FieldErrors) > 0 {
		t := app.newTemplateData(r)
		t.Form = form
		t.Blog = form.Blog
		app.render(w, r, http.StatusUnprocessableEntity, "admin.new_blog.tmpl.html", &t)
		return
	}

	if id == 0 {
		err = app.blogs.Insert(form.Blog)
		if err != nil {
			log.Printf("app.blogs.Insert(): %v\n", err)
			// TODO: send some notification to the UI
		}
	} else {
		err = app.blogs.UpdateExisting(form.Blog)
		if err != nil {
			log.Printf("app.blogs.UpdateExisting(): %v\n", err)
			// TODO: send some notification to the UI
		}
	}
	http.Redirect(w, r, fmt.Sprintf("/blog/%v", form.Blog.Path), http.StatusSeeOther)
	return
}
