package main

import (
	"net/http"

	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	// standard := alice.New(logRequest, commonHeaders)
	standard := alice.New(commonHeaders, app.identify)
	usersOnly := alice.New(app.requireAuthentication)
	adminsOnly := alice.New(app.requireAuthentication, app.requireAdmin)

	mux.HandleFunc("GET /", app.home)

	// simple pages:
	mux.HandleFunc("GET /now", app.now)
	mux.HandleFunc("GET /about", app.about)
	mux.HandleFunc("GET /bookshelf", app.bookshelf)
	// protected:
	mux.Handle("GET /admin/new-page", adminsOnly.ThenFunc(app.adminNewPage))
	mux.Handle("GET /admin/pages/{page}", adminsOnly.ThenFunc(app.adminPagesPage))
	mux.Handle("POST /pages", adminsOnly.ThenFunc(app.pagesPost))

	// blogs:
	mux.HandleFunc("GET /blog", app.blog)
	mux.HandleFunc("GET /blog/{path}", app.blogPath)
	// protected:
	mux.Handle("POST /blog", adminsOnly.ThenFunc(app.blogPost))
	mux.Handle("GET /admin/new-blog", adminsOnly.ThenFunc(app.adminNewBlog))
	mux.Handle("GET /admin/blog/{path}", adminsOnly.ThenFunc(app.adminBlogPath))

	// til:
	mux.HandleFunc("GET /today-i-learned", app.todayILearned)
	mux.HandleFunc("GET /today-i-learned/{path}", app.todayILearnedPath)
	// protected:
	mux.Handle("POST /til", adminsOnly.ThenFunc(app.tilPost))
	mux.Handle("GET /admin/new-til", adminsOnly.ThenFunc(app.adminNewTIL))
	mux.Handle("GET /admin/tils/{path}", adminsOnly.ThenFunc(app.adminTILSTIL))

	//songbook:
	mux.HandleFunc("GET /songbook", app.songbook)
	mux.HandleFunc("GET /songbook/{song}", app.songbookSong)
	// protected:
	mux.Handle("POST /songbook", adminsOnly.ThenFunc(app.songbookPost))
	mux.Handle("GET /admin/new-song", adminsOnly.ThenFunc(app.adminNewSong))
	mux.Handle("GET /admin/songbook/{song}", adminsOnly.ThenFunc(app.adminSongbookSong))

	// Bellevue Activities (all protected):
	mux.Handle("GET /admin/new-bellevue-activity", usersOnly.ThenFunc(app.adminNewBellevueActivity))
	mux.Handle("GET /bellevue-activities", usersOnly.ThenFunc(app.bellevueActivities))
	mux.Handle("POST /bellevue-activities", usersOnly.ThenFunc(app.bellevueActivityPost))
	mux.Handle("PUT /bellevue-activities/{id}", usersOnly.ThenFunc(app.bellevueActivityPut))
	mux.Handle("DELETE /bellevue-activities/{id}", usersOnly.ThenFunc(app.bellevueActivityDelete))
	// HTMX partials
	mux.Handle("GET /bellevue-activities/{id}/edit", usersOnly.ThenFunc(app.bellevueActivityIDEdit))

	// admin:
	mux.HandleFunc("GET /admin/login", app.adminLogin)
	mux.HandleFunc("POST /admin/login", app.adminLoginPost)
	// protected:
	mux.Handle("GET /admin", usersOnly.ThenFunc(app.admin))
	mux.Handle("GET /admin/logout", usersOnly.ThenFunc(app.adminLogoutPost))

	// finances:
	// protected:
	mux.Handle("GET /finances", adminsOnly.ThenFunc(app.finances))

	return standard.Then(mux)
}
