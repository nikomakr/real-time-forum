package main

import (
	"log"
	"net/http"
	"real-time-forum/db"
	"real-time-forum/handlers"
	"real-time-forum/utils"
)

func main() {
	db.Init("./forum.db")

	http.HandleFunc("/api/register", handlers.Register)
	http.HandleFunc("/api/login", handlers.Login)

	// The following endpoint is protected by the RequireAuth middleware, which checks for a valid session cookie and ensures the user is authenticated before allowing access to the /api/me endpoint. If the user is authenticated, their user ID is returned in the response.
	http.HandleFunc("/api/me", handlers.RequireAuth(func(w http.ResponseWriter, r *http.Request) {
		userID := handlers.GetUserID(r)
		if err := utils.WriteJSON(w, http.StatusOK, map[string]string{"user_id": userID}); err != nil {
			log.Printf("could not write response: %v", err)
		}
	}))

	log.Println("server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
