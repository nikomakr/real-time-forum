package main

import (
	"log"
	"net/http"
	"real-time-forum/db"
)

func main() {
	db.Init("./forum.db")
	log.Println("server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}