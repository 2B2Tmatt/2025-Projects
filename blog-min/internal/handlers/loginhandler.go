package handlers

import (
	"blog-min/internal/encryption"
	sessions "blog-min/internal/session"
	"database/sql"
	"log"
	"net/http"
	"time"
)

func Login(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	//Data Validation
	if r.Method == http.MethodGet {
		http.ServeFile(w, r, "web/templates/login.html")
		return
	}
	err := r.ParseForm()
	if err != nil {
		log.Println("Error parsing form", err)
		http.Error(w, "Error processing form", http.StatusInternalServerError)
		return
	}
	username := r.FormValue("username")
	password := r.FormValue("password")

	sqlStatement := `
	SELECT password_hash,id FROM users WHERE display_name=$1
	`
	var dbhash string
	var uid int64
	db.QueryRow(sqlStatement, username).Scan(&dbhash, &uid)
	if !encryption.CheckPasswordHash(password, dbhash) {
		http.Error(w, "username or password incorrect", http.StatusUnauthorized)
		return
	}
	_, err = sessions.CreateSession(w, db, uid, time.Hour*24)
	if err != nil {
		log.Println("Error creating session", err)
		http.Error(w, "Error creating session", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/pages", http.StatusSeeOther)
	//Search for user if user actually exists

	//is user active

	// verify password

	//Generate Token

	//Send token as a response or cookie
}
