package handlers

import (
	sessions "blog-min/internal/session"
	"database/sql"
	"log"
	"net/http"
)

func Post(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	if r.Method == http.MethodGet {
		cookie, err := r.Cookie("session_token")
		if err != nil {
			log.Println("Error 1")
			http.ServeFile(w, r, "web/templates/blocked.html")
			return
		}

		uid, err := sessions.CheckSession(db, cookie.Value)
		if err != nil {
			log.Println("DB error:", err)
			http.Error(w, "server error", 500)
			return
		}
		if uid == 0 {
			log.Println("No valid session for sid", cookie.Value)
			http.Error(w, "invalid session", 401)
			return
		}
		http.ServeFile(w, r, "web/templates/blocked.html")
		http.ServeFile(w, r, "web/templates/post.html")
		return
	}
	//POSTS
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
