package handlers

import (
	"log"
	"net/http"
)

func Signup(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println("Error parsing form", err)
		http.Error(w, "Error processing form", http.StatusInternalServerError)
		return
	}
	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")

	sqlStatement := `
	INSERT INTO users (display_name, email, password)
	VALUES ($1, $2, $3)
	`
	err = db.QueryRow(sqlStatement, 30)

}
