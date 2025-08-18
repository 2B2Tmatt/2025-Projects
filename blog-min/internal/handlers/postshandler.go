package handlers

import "net/http"

func Posts(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "blog-min/web/templates/posts.html")
}
