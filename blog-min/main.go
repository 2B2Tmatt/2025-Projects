package main

import (
	"blog-min/internal/dbconn"
	"blog-min/internal/middleware"
	"log"
	"net/http"

	"blog-min/internal/handlers"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func Working(w http.ResponseWriter, r *http.Request) {
	log.Println("I work")
}

func main() {
	r := mux.NewRouter()
	PORT := ":9000"

	db, err := dbconn.OpenDB()
	if err != nil {
		log.Fatal("Error opening db")
	}
	err = dbconn.InitSchema(db)
	if err != nil {
		log.Fatal("Error creating tables")
	}
	defer db.Close()
	r.HandleFunc("/", handlers.Redirect)
	r.HandleFunc("/pages", func(w http.ResponseWriter, r *http.Request) { handlers.Home(w, r, db) })
	r.HandleFunc("/pages/signup", func(w http.ResponseWriter, r *http.Request) { handlers.Signup(w, r, db) })
	r.HandleFunc("/pages/login", func(w http.ResponseWriter, r *http.Request) { handlers.Login(w, r, db) })
	// r.HandleFunc("/pages/post", Working)
	r.HandleFunc("/pages/post",
		middleware.RequireSession(db, func(w http.ResponseWriter, r *http.Request, uid int64) {
			handlers.Post(w, r, db, uid)
		}),
	).Methods(http.MethodGet, http.MethodPost)
	r.HandleFunc("/pages/posts", func(w http.ResponseWriter, r *http.Request) { handlers.Posts(w, r, db) })
	r.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) { handlers.Logout(w, r, db) })
	r.HandleFunc("/pages/{username}", func(w http.ResponseWriter, r *http.Request) { handlers.User(w, r, db) })
	log.Println("Running on Port:", PORT)
	// err := http.ListenAndServeTLS(PORT, "server.crt", "server.key", r)
	wrapped := middleware.NoStore(r)
	wrapped = middleware.RequestLog(wrapped)
	err = http.ListenAndServe(PORT, wrapped)
	if err != nil {
		log.Fatal("Listen and Serve:", err)
	}
}
