package main

import (
	"log"
	"net/http"
	"real-time-forum/db"
	"real-time-forum/handlers"
)

func main() {
	db.Init("./forum.db")

	http.HandleFunc("/api/register", handlers.Register)

	log.Println("server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
