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
			uagent = r.Header.Get("User-Agent")
			platf  = r.Header.Get("Sec-Ch-Ua-Platform")
		)

		// caddy will set X-Forwarded-For with original src IP when reverse proxying.
		// r.RemoteAddr will be localhost, in that case.
		xff := r.Header.Get("X-Forwarded-For")
		if xff != "" {
			ip = xff
		}

		log.Printf("msg=ReceivedRequest ip=%v proto=%v method=%v uri=%v platf=%v user-agent=%v", ip, proto, method, uri, platf, uagent)

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

		highlightJS := "'sha256-KuW8nrMYej09eTtZkBNDwTy8Yn05dABB5v2dLSEPgTY='"

		w.Header().Set(
			"Content-Security-Policy",
			"default-src 'self';"+
				"img-src 'self' images.ctfassets.net;"+
				"script-src 'self' cdnjs.cloudflare.com " + highlightJS + ";"+
				"style-src 'self' cdnjs.cloudflare.com fonts.googleapis.com;"+
				"font-src fonts.gstatic.com",
		)
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")

		w.Header().Set("Server", "Go")

		next.ServeHTTP(w, r)
	})
}
