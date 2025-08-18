package handlers

import (
	"net/http"
	"time"
)

func Redirect(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "blog-min/web/templates/redirect01.html")
	time.Sleep(time.Second)
	http.ServeFile(w, r, "blog-min/web/templates/home/html")
}
