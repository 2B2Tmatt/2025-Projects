package main

import (
	"cache-demo/internal/handlers"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	PORT := ":9000"
	r := mux.NewRouter()
	r.HandleFunc("/pages/home", handlers.Home)
	log.Println("Start server on port", PORT)
	err := http.ListenAndServe(PORT, r)
	if err != nil {
		log.Fatal("Listen and Serve:", err)
	}
}
