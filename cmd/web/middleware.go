package main

import (
	"log"
	"net/http"
)

func logRequest(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        var (
            ip     = r.RemoteAddr
            proto  = r.Proto
            method = r.Method
            uri    = r.URL.RequestURI()
        )

        log.Printf("msg=received request ip=%v proto=%v method=%v uri=%v", ip, proto, method, uri)

        next.ServeHTTP(w, r)
    })
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

// see https://owasp.org/www-project-secure-headers/
// see https://developer.mozilla.org/en-US/docs/Web/HTTP/Guides/CSP
func commonHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: rm cdn.jsdelivr.net soon!
		w.Header().Set(
			"Content-Security-Policy",
			"default-src 'self'; style-src 'self' cdn.jsdelivr.net fonts.googleapis.com; font-src fonts.gstatic.com",
		)
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")

		w.Header().Set("Server", "Go")

		next.ServeHTTP(w, r)
	})
}
