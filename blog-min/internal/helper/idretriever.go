package helper

import (
	sessions "blog-min/internal/session"
	"database/sql"
	"log"
	"net/http"
)

func GetId(db *sql.DB, r *http.Request) (int64, error) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		log.Println("No cookie")
		return 0, err
	}
	uid, exists, _ := sessions.CheckSession(db, cookie.Value)
	if !exists {
		log.Println("Cookie not found in DB")
		return 0, err
	}
	return uid, nil
}
