package main

import (
	"net/http"

	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	standard := alice.New(logRequest, commonHeaders)
	protected := alice.New(app.requireAuthentication)

	mux.HandleFunc("GET /", app.home)

	// simple pages:
	mux.HandleFunc("GET /now", app.now)
	mux.HandleFunc("GET /about", app.about)
	mux.HandleFunc("GET /blog", app.blog)
	mux.HandleFunc("GET /bookshelf", app.bookshelf)
	// protected:
	mux.Handle("GET /admin/new-page", protected.ThenFunc(app.adminNewPage))
	mux.Handle("GET /admin/pages/{page}", protected.ThenFunc(app.adminPagesPage))
	mux.Handle("POST /pages", protected.ThenFunc(app.pagesPost))

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

	// admin:
	mux.HandleFunc("GET /admin/login", app.adminLogin)
	mux.HandleFunc("POST /admin/login", app.adminLoginPost)
	// protected:
	mux.Handle("GET /admin", protected.ThenFunc(app.admin))
	mux.Handle("GET /admin/logout", protected.ThenFunc(app.adminLogoutPost))

	// finances:
	// protected:
	mux.Handle("GET /finances", protected.ThenFunc(app.finances))

	return standard.Then(mux)
}
