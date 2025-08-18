package handlers

import (
	"database/sql"
	"log"
	"net/http"
)

func Post(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method == http.MethodGet {
		http.ServeFile(w, r, "web/templates/post.html")
		return
	}

	err := r.ParseForm()
	if err != nil {
		log.Println("Error parsing form", err)
		http.Error(w, "Error processing form", http.StatusInternalServerError)
		return
	}
	title := r.FormValue("post_title")
	body := r.FormValue("blog")

	sqlStatement := `
	INSERT INTO posts (user_id, title, body)
	VALUES ($1, $2, $3)
	`
	_, err = db.Exec(sqlStatement, 1, title, body)
	if err != nil {
		log.Println("Error processing form", err)
		http.Error(w, "Error processing form", http.StatusInternalServerError)
		return
	}
	http.ServeFile(w, r, "web/templates/redirect.html")

}
