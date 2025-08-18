package handlers

import (
	"blog-min/internal/encryption"
	"database/sql"
	"log"
	"net/http"
)

func Login(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	//Data Validation
	if r.Method == http.MethodGet {
		http.ServeFile(w, r, "blog-min/web/templates/login.html")
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
	password, err = encryption.GeneratePassword(password)
	if err != nil {
		log.Println("Error creating hash", err)
		http.Error(w, "Error processing information", http.StatusInternalServerError)
		return
	}

	sqlStatement := `
	SELECT password_hash FROM users WHERE display_name=$1
	`
	var dbpass string
	db.QueryRow(sqlStatement, username).Scan(&dbpass)
	if password != dbpass {
		http.Error(w, "username or password incorrect", http.StatusUnauthorized)
		return
	}
	http.ServeFile(w, r, "blog-min/web/templates/redirect.html")

	//Search for user if user actually exists

	//is user active

	// verify password

	//Generate Token

	//Send token as a response or cookie
}
