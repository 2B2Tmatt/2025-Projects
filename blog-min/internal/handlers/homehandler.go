package handlers

import (
	"blog-min/internal/helper"
	"database/sql"
	"log"
	"net/http"
)

func Home(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	log.Println("This was reached")
	anon := &helper.User{User: "Guest"}
	uid, err := helper.GetId(db, r)
	if err != nil {
		_ = helper.Render(w, anon, "web/templates/home.html")
		return
	}
	user, err := helper.GetUsername(db, uid)
	if err != nil || user == nil {
		_ = helper.Render(w, anon, "web/templates/home.html")
		return
	}
	helper.Render(w, user, "web/templates/home.html")
}
