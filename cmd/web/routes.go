package main

import "net/http"

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("GET /", app.home)

	// simple pages:
	mux.HandleFunc("GET /now", app.now)
	mux.HandleFunc("GET /about", app.about)
	mux.HandleFunc("GET /blog", app.blog)
	mux.HandleFunc("GET /bookshelf", app.bookshelf)
	mux.Handle(
		"GET /admin/new-page",
		app.requireAuthentication(http.HandlerFunc(
			app.adminNewPage,
		)),
	)
	mux.Handle(
		"GET /admin/pages/{page}",
		app.requireAuthentication(http.HandlerFunc(
			app.adminPagesPage,
		)),
	)
	mux.Handle(
		"POST /pages",
		app.requireAuthentication(http.HandlerFunc(
			app.pagesPost,
		)),
	)

	// til:
	mux.HandleFunc("GET /today-i-learned", app.todayILearned)
	mux.HandleFunc("GET /today-i-learned/{path}", app.todayILearnedPath)
	mux.Handle(
		"POST /til",
		app.requireAuthentication(http.HandlerFunc(
			app.tilPost,
		)),
	)
	mux.Handle(
		"GET /admin/new-til",
		app.requireAuthentication(http.HandlerFunc(
			app.adminNewTIL,
		)),
	)
	mux.Handle(
		"GET /admin/tils/{path}",
		app.requireAuthentication(http.HandlerFunc(
			app.adminTILSTIL,
		)),
	)

	//songbook:
	mux.HandleFunc("GET /songbook", app.songbook)
	mux.HandleFunc("GET /songbook/{song}", app.songbookSong)
	// protected:
	mux.Handle(
		"POST /songbook",
		app.requireAuthentication(http.HandlerFunc(
			app.songbookPost,
		)),
	)
	mux.Handle(
		"GET /admin/new-song",
		app.requireAuthentication(http.HandlerFunc(
			app.adminNewSong,
		)),
	)
	mux.Handle(
		"GET /admin/songbook/{song}",
		app.requireAuthentication(http.HandlerFunc(
			app.adminSongbookSong,
		)),
	)

	// admin:
	mux.HandleFunc("GET /admin/login", app.adminLogin)
	mux.HandleFunc("POST /admin/login", app.adminLoginPost)
	// protected:
	mux.Handle(
		"GET /admin",
		app.requireAuthentication(http.HandlerFunc(
			app.admin,
		)),
	)
	mux.Handle(
		"GET /admin/logout",
		app.requireAuthentication(http.HandlerFunc(
			app.adminLogoutPost,
		)),
	)

	// finances:
	// protected:
	mux.Handle(
		"GET /finances",
		app.requireAuthentication(http.HandlerFunc(
			app.finances,
		)),
	)

	return mux
}
