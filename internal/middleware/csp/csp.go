package csp

import "net/http"

const CSP = "default-src 'none'; script-src 'self'; connect-src 'self'; img-src 'self'; style-src 'self'; base-uri 'self'; form-action 'self'"

func CspMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy", CSP)
		if r.Method == http.MethodOptions {
			return
		}
		next.ServeHTTP(w, r)
	})
}
