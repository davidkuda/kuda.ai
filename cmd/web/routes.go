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
	standard := alice.New(commonHeaders)
	identity := alice.New(app.requireAuthentication)
	protected := alice.New(app.requireAuthentication, app.requireAdmin)

	mux.HandleFunc("GET /", app.home)

	// simple pages:
	mux.HandleFunc("GET /now", app.nowFromDB)
	mux.HandleFunc("GET /about", app.about)
	mux.HandleFunc("GET /bookshelf", app.bookshelf)
	// protected:
	mux.Handle("GET /admin/new-page", protected.ThenFunc(app.adminNewPage))
	mux.Handle("GET /admin/pages/{page}", protected.ThenFunc(app.adminPagesPage))
	mux.Handle("POST /pages", protected.ThenFunc(app.pagesPost))

	// blogs:
	mux.HandleFunc("GET /blog", app.blog)
	mux.HandleFunc("GET /blog/{path}", app.blogPath)
	// protected:
	mux.Handle("POST /blog", protected.ThenFunc(app.blogPost))
	mux.Handle("GET /admin/new-blog", protected.ThenFunc(app.adminNewBlog))
	mux.Handle("GET /admin/blog/{path}", protected.ThenFunc(app.adminBlogPath))

	// til:
	mux.HandleFunc("GET /today-i-learned", app.todayILearned)
	mux.HandleFunc("GET /today-i-learned/{path}", app.todayILearnedPath)
	// protected:
	mux.Handle("POST /til", protected.ThenFunc(app.tilPost))
	mux.Handle("GET /admin/new-til", protected.ThenFunc(app.adminNewTIL))
	mux.Handle("GET /admin/tils/{path}", protected.ThenFunc(app.adminTILSTIL))

	//songbook:
	mux.HandleFunc("GET /songbook", app.songbook)
	mux.HandleFunc("GET /songbook/{song}", app.songbookSong)
	// protected:
	mux.Handle("POST /songbook", protected.ThenFunc(app.songbookPost))
	mux.Handle("GET /admin/new-song", protected.ThenFunc(app.adminNewSong))
	mux.Handle("GET /admin/songbook/{song}", protected.ThenFunc(app.adminSongbookSong))

	// Bellevue Activities (all protected):
	mux.Handle("GET /admin/new-bellevue-activity", identity.ThenFunc(app.adminNewBellevueActivity))
	mux.Handle("GET /bellevue-activities", identity.ThenFunc(app.bellevueActivities))
	mux.Handle("POST /bellevue-activities", identity.ThenFunc(app.bellevueActivityPost))

	// admin:
	mux.HandleFunc("GET /admin/login", app.adminLogin)
	mux.HandleFunc("POST /admin/login", app.adminLoginPost)
	// protected:
	mux.Handle("GET /admin", identity.ThenFunc(app.admin))
	mux.Handle("GET /admin/logout", identity.ThenFunc(app.adminLogoutPost))

	// finances:
	// protected:
	mux.Handle("GET /finances", protected.ThenFunc(app.finances))

	return standard.Then(mux)
}
