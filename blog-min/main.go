package main

import (
	makesql "blog-min/internal/sql"
	"log"
	"net/http"

	"blog-min/internal/handlers"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	r := mux.NewRouter()
	PORT := ":8080"

	db := makesql.OpenDB()
	defer db.Close()
	r.HandleFunc("/", handlers.Redirect)
	r.HandleFunc("/pages", handlers.Home)
	r.HandleFunc("/pages/signup", func(w http.ResponseWriter, r *http.Request) { handlers.Signup(w, r, db) })
	r.HandleFunc("/pages/login", func(w http.ResponseWriter, r *http.Request) { handlers.Login(w, r, db) })
	r.HandleFunc("/pages/post", func(w http.ResponseWriter, r *http.Request) { handlers.Post(w, r, db) })
	r.HandleFunc("/pages/posts", handlers.Posts)

	log.Println("Running on Port:", PORT)
	// err := http.ListenAndServeTLS(PORT, "server.crt", "server.key", r)
	err := http.ListenAndServe(PORT, r)
	if err != nil {
		log.Fatal("Listen and Serve:", err)
	}
}
