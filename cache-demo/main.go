package main

import (
	"cache-demo/internal/cache"
	"cache-demo/internal/handlers"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	cache := cache.CreateCache()
	PORT := ":9000"
	r := mux.NewRouter()
	r.HandleFunc("/pages/home", func(w http.ResponseWriter, r *http.Request) { handlers.Home(w, r, cache) })
	log.Println("Start server on port", PORT)
	err := http.ListenAndServe(PORT, r)
	if err != nil {
		log.Fatal("Listen and Serve:", err)
	}
}
