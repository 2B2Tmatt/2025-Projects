package main

import (
	makesql "blog-min/internal/sql"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	r := mux.NewRouter()
	PORT := ":8080"

	db := makesql.OpenDB()
	r.Handle("/", http.HandlerFunc(Redirect))
	r.Handle("/pages/", http.HandlerFunc(Home))
	r.Handle("/pages/signup", http.HandlerFunc(Signup))
	r.Handle("/pages/login", http.HandlerFunc(Login))
	r.Handle("/pages/post", http.HandlerFunc(Post))
	r.Handle("/pages/posts/", http.HandlerFunc(Posts))
	err = http.ListenAndServeTLS(PORT, "server.crt", "server.key", r)
	if err != nil {
		log.Fatal("Listen and Serve:", err)
	}
}
