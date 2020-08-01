package middlewares

import (
	"github.com/gorilla/context"
	"mypatterngo/app/helpers"
	"net/http"
	"os"
	"strings"
)

//Guest Middleware
func GuestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		var key = r.Header.Get("API-Key")
		key = strings.TrimSpace(key)

		if key != os.Getenv("API_KEY") {
			helpers.Error(w, http.StatusBadRequest, "Missing Authorization")
			return
		}
		next.ServeHTTP(w, r)
	})
}


//With JWT Access Middleware
func AuthenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var key = r.Header.Get("API-Key")
		key = strings.TrimSpace(key)

		if key != os.Getenv("API_KEY") {
			helpers.Error(w, http.StatusBadRequest, "Missing Authorization")
			return
		}

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		err := helpers.TokenValid(r)
		if err != nil {
			helpers.Error(w, http.StatusUnauthorized, "You can't access this page")
			return
		}
		next.ServeHTTP(w, r)
	})
}

// Admin Middleware
func AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var key = r.Header.Get("API-Key")
		key = strings.TrimSpace(key)

		if key != os.Getenv("API_KEY") {
			helpers.Error(w, http.StatusBadRequest, "Missing Authorization")
			return
		}

		roleName := context.Get(r, "RoleName")

		if roleName != "High Admin" {
			helpers.Error(w, http.StatusUnauthorized, "You can't access this page")
			return
		}

		next.ServeHTTP(w, r)
	})
}
