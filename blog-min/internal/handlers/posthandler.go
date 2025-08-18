package handlers

import (
	"database/sql"
	"net/http"
)

func Post(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method == http.MethodGet {
		http.ServeFile(w, r, "blog-min/web/templates/post.html")
	}

}
