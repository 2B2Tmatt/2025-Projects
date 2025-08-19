package handlers

import (
	sessions "blog-min/internal/session"
	"database/sql"
	"log"
	"net/http"
)

func Logout(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		http.Error(w, "Error signing out", http.StatusInternalServerError)
		log.Println("Cookie doesn't exist")
		return
	}
	sid := cookie.Value
	err = sessions.EndSession(w, db, sid)
	if err != nil {
		http.Error(w, "Error signing out", http.StatusInternalServerError)
		log.Println("Error deleting session")
		return
	}
	log.Println("Sucessfully logged out")
	http.Redirect(w, r, "/pages", http.StatusSeeOther)
}
