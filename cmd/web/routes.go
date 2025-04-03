package main

import "net/http"

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("GET /", app.home)

	// simple pages:
	mux.HandleFunc("GET /about", app.about)
	mux.HandleFunc("GET /blog", app.blog)
	mux.HandleFunc("GET /bookshelf", app.bookshelf)
	mux.HandleFunc("GET /today-i-learned", app.todayILearned)

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
	mux.Handle(
		"GET /admin/new-til",
		app.requireAuthentication(http.HandlerFunc(
			app.adminNewTIL,
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

func (app *application) requireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !app.isAuthenticated(r) {
			http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
			return
		}
		w.Header().Add("Cache-Control", "no-store")
		next.ServeHTTP(w, r)
	})
}
