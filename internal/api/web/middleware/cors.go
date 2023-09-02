package middleware

import (
	"net/http"
)

func (m middleware) EnableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000" /*see 'Origin' in request object*/)

		/*
		 * To expose content-disposition header which can contains filename to a client when it's downloading a file.
		 * Otherwise, the client doesn't see the filename which it's downloading in the browser. This header is simply not shown:
		 * content-disposition: "attachment; filename=25351.pptx"
		 *
		 * That's security restriction of the browser, not client. And this is related to CORS.
		 */
		w.Header().Set("Access-Control-Expose-Headers", "*")

		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, X-CSRF-Token, Authorization")
		} else {
			// For client app, which is running inside webpack, to be removed once dev phase is complete
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			next.ServeHTTP(w, r)
		}
	})
}