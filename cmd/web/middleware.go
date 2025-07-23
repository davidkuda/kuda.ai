package main

import (
	"context"
	"fmt"
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

func (app *application) identify(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Cache-Control", "no-store")

		var userEmail string
		var userID int

		userEmail, err := app.extractUserFromJWTCookie(r)
		if err != nil {
			ctx := r.Context()
			ctx = context.WithValue(ctx, "isAuthenticated", false)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
			return
		}

		userID, err = app.users.GetUserIDByEmail(userEmail)
		if err != nil {
			err = fmt.Errorf("could not get email from user with email %s: %v\n", userEmail, err)
			app.serverError(w, r, err)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, "isAuthenticated", true)
		ctx = context.WithValue(ctx, "userEmail", userEmail)
		ctx = context.WithValue(ctx, "userID", userID)

		// TODO: implement nice permission management...
		if userID == 1 {
			ctx = context.WithValue(ctx, "isAdmin", true)
		}

		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func (app *application) requireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: think this over, this makes another JWT decoding
		if !app.isAuthenticated(r) {
			http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (app *application) requireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		isAdmin, ok := r.Context().Value("isAdmin").(bool)
		if !ok {
			app.renderError(w, r, http.StatusForbidden)
			return
		}
		if !isAdmin {
			app.renderError(w, r, http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// see https://owasp.org/www-project-secure-headers/
// see https://developer.mozilla.org/en-US/docs/Web/HTTP/Guides/CSP
func commonHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		setTheme := "'sha256-lQ7hEV5vzkCFUSGHetH4H+fWaAnYTAAiEMxdfS0bTeU='"
		highlightJS := "'sha256-kesGYQCKRT1Io3waiBp5a4n4ZLg1Xbn8ldhKQWp/hco='"
		allowedInlineJS := fmt.Sprintf("%s %s", setTheme, highlightJS)

		w.Header().Set(
			"Content-Security-Policy",
			"default-src 'self';"+
				"img-src 'self' images.ctfassets.net;"+
				"script-src 'self' cdnjs.cloudflare.com "+allowedInlineJS+";"+
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
