package handlers

import (
	"database/sql"
	"net/http"
)

func Post(w http.ResponseWriter, r *http.Request, db *sql.DB, uid int64) {
	if r.Method == http.MethodGet {
		http.ServeFile(w, r, "web/templates/post.html")
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form", http.StatusInternalServerError)
		return
	}
	title := r.FormValue("post_title")
	body := r.FormValue("blog")

	_, err := db.Exec(`INSERT INTO posts (user_id, title, body) VALUES ($1,$2,$3)`,
		uid, title, body)
	if err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/pages", http.StatusSeeOther)
}
