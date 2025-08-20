package handlers

import (
	"blog-min/internal/encryption"
	"database/sql"
	"log"
	"net/http"
)

func Signup(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method == http.MethodGet {
		http.ServeFile(w, r, "web/templates/signup.html")
		return
	}
	err := r.ParseForm()
	if err != nil {
		log.Println("Error parsing form", err)
		http.Error(w, "Error processing form", http.StatusInternalServerError)
		return
	}
	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")
	password, err = encryption.HashPassword(password)
	if err != nil {
		log.Println("Error creating hash", err)
		http.Error(w, "Error processing information", http.StatusInternalServerError)
		return
	}

	sqlStatement := `
	INSERT INTO users (display_name, email, password_hash)
	VALUES ($1, $2, $3)
	`
	_, err = db.Exec(sqlStatement, username, email, password)
	if err != nil {
		log.Println("Error processing form", err)
		http.Error(w, "Error processing form", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/pages", http.StatusSeeOther)
}
