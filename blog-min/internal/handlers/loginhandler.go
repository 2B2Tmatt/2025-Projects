package handlers

import (
	"blog-min/internal/encryption"
	makesql "blog-min/internal/sql"
	"log"
	"net/http"
)

func Login(w http.ResponseWriter, r *http.Request) {
	//Data Validation
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
		log.Println("Error parsing form", err)
		http.Error(w, "Error processing form", http.StatusInternalServerError)
		return
	}

	sqlStatement := `
	INSERT INTO users (display_name, email, password)
	VALUES ($1, $2, $3)
	`
	db := makesql.OpenDB()
	err = db.QueryRow(sqlStatement, 30)

	//Search for user if user actually exists

	//is user active

	// verify password

	//Generate Token

	//Send token as a response or cookie
}
