package main

import (
	"cache-demo/internal/cache"
	"cache-demo/internal/handlers"
	"log"
	"net/http"
	"time"

	mw "cache-demo/internal/middleware"

	"github.com/gorilla/mux"
)

func main() {
	cache := cache.CreateCache()
	go cache.BackgroundProcessor()
	PORT := ":9000"
	r := mux.NewRouter()
	r.HandleFunc("/pages/home", func(w http.ResponseWriter, r *http.Request) { handlers.Home(w, r, cache) })
	log.Println("Starting server on port", PORT)
	rateLimiter := mw.NewRateLimiter(120, time.Second)
	wrapped := mw.LimitRates(r, rateLimiter)
	err := http.ListenAndServe(PORT, wrapped)
	if err != nil {
		log.Fatal("Listen and Serve:", err)
	}
}
