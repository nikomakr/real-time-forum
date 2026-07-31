package main

import (
	"log"
	"net/http"
	"real-time-forum/db"
	"real-time-forum/handlers"
	"real-time-forum/utils"
	"strings"
)

func main() {
	db.Init("./forum.db")

	http.Handle("/", http.FileServer(http.Dir("./static")))

	http.HandleFunc("/api/register", handlers.Register)
	http.HandleFunc("/api/login", handlers.Login)
	http.HandleFunc("/api/logout", handlers.RequireAuth(handlers.Logout))
	http.HandleFunc("/api/posts/", handlers.RequireAuth(handlers.GetPost))

	// The following endpoint is protected by the RequireAuth middleware, which checks for a valid session cookie and ensures the user is authenticated before allowing access to the /api/me endpoint. If the user is authenticated, their user ID is returned in the response.
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

	log.Println("server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
