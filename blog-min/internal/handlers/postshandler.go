package handlers

import (
	"net/http"
)

func Posts(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/templates/posts.html")
}
