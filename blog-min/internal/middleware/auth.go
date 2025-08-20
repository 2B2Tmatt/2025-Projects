package middleware

import (
	sessions "blog-min/internal/session"
	"database/sql"
	"log"
	"net/http"
)

const SessionCookieName = "session_token"

type AuthedHandler func(w http.ResponseWriter, r *http.Request, uid int64)

func RequireSession(db *sql.DB, next AuthedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(SessionCookieName)
		if err != nil {
			http.ServeFile(w, r, "web/templates/blocked.html")
			log.Println("No cookie")
			return
		}
		uid, exists, _ := sessions.CheckSession(db, cookie.Value)
		if !exists {
			log.Println("Cookie not found in DB")
			http.ServeFile(w, r, "web/templates/blocked.html")
			return
		}
		log.Println("Line reached")
		next(w, r, uid)
	}
}
