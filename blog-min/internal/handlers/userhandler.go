package handlers

import (
	"blog-min/internal/helper"
	"database/sql"
	"net/http"

	"github.com/gorilla/mux"
)

func User(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	vars := mux.Vars(r)
	username := vars["username"]
	page := helper.PageData{
		User:  &helper.User{User: "Guest"},
		Posts: nil,
	}
	uid, err := helper.GetId(db, r)
	if err != nil {
		_ = helper.Render(w, page, "web/templates/posts.html")
		return
	}
	user, err := helper.GetUsername(db, uid)
	if err != nil || user == nil {
		_ = helper.Render(w, page, "web/templates/posts.html")
		return
	}

	posts, err := helper.GetUserPost(db, username)
	if err != nil {
		_ = helper.Render(w, page, "web/templates/posts.html")
		return
	}
	page.User = user
	page.Posts = posts
	helper.Render(w, page, "web/templates/posts.html")

}
