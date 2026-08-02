package main

import (
	"log"
	"net/http"
	"strings"

	"real-time-forum/db"
	"real-time-forum/handlers"
	"real-time-forum/utils"
	"real-time-forum/ws"
)

func main() {
	db.Init("./forum.db")

	hub := ws.NewHub()

	http.Handle("/", http.FileServer(http.Dir("./static")))

	http.HandleFunc("/api/register", handlers.Register)
	http.HandleFunc("/api/login", handlers.Login)
	http.HandleFunc("/api/logout", handlers.RequireAuth(handlers.Logout))

	http.HandleFunc("/api/me", handlers.RequireAuth(func(w http.ResponseWriter, r *http.Request) {
		userID := handlers.GetUserID(r)
		if err := utils.WriteJSON(w, http.StatusOK, map[string]string{"user_id": userID}); err != nil {
			log.Printf("could not write response: %v", err)
		}
	}))

	http.HandleFunc("/api/posts", handlers.RequireAuth(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			handlers.GetPosts(w, r)
		} else if r.Method == http.MethodPost {
			handlers.CreatePost(w, r)
		} else {
			utils.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	}))

	http.HandleFunc("/api/posts/", handlers.RequireAuth(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/comments") {
			if r.Method == http.MethodGet {
				handlers.GetComments(w, r)
			} else if r.Method == http.MethodPost {
				handlers.CreateComment(w, r)
			} else {
				utils.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
			}
		} else {
			handlers.GetPost(w, r)
		}
	}))

	http.HandleFunc("/api/messages/{id}", handlers.RequireAuth(handlers.GetMessages))
	http.HandleFunc("/api/users", handlers.RequireAuth(handlers.GetUsers(hub)))

	// WebSocket endpoint — authenticated
	http.HandleFunc("/ws", handlers.RequireAuth(func(w http.ResponseWriter, r *http.Request) {
		handlers.ServeWS(hub, w, r)
	}))

	log.Println("server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
