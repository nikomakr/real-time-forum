package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"path"
	"strings"

	"real-time-forum/db"
	"real-time-forum/utils"
)

type postResponse struct {
	ID           string   `json:"id"`
	Title        string   `json:"title"`
	Content      string   `json:"content"`
	Author       string   `json:"author"`
	Categories   []string `json:"categories"`
	CommentCount int      `json:"comment_count"`
	LikeCount    int      `json:"like_count"`
	Liked        bool     `json:"liked"`
	CreatedAt    string   `json:"created_at"`
}

type postDetailResponse struct {
	ID           string   `json:"id"`
	Title        string   `json:"title"`
	Content      string   `json:"content"`
	Author       string   `json:"author"`
	Categories   []string `json:"categories"`
	CommentCount int      `json:"comment_count"`
	LikeCount    int      `json:"like_count"`
	Liked        bool     `json:"liked"`
	CreatedAt    string   `json:"created_at"`
}

func GetPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// Extract post ID from URL path — /api/posts/{id}
	id := path.Base(r.URL.Path)
	if id == "" || id == "/" || id == "posts" {
		utils.WriteError(w, http.StatusBadRequest, "post ID is required")
		return
	}

	var post postDetailResponse
	var categoryStr string
	var liked int

	err := db.DB.QueryRow(`
		SELECT
			p.id,
			p.title,
			p.content,
			p.created_at,
			u.nickname,
			COALESCE((
				SELECT GROUP_CONCAT(c.name, ',')
				FROM post_categories pc
				JOIN categories c ON pc.category_id = c.id
				WHERE pc.post_id = p.id
			), '') AS categories,
			(SELECT COUNT(*) FROM comments WHERE post_id = p.id) AS comment_count,
			(SELECT COUNT(*) FROM post_likes WHERE post_id = p.id) AS like_count,
			EXISTS(SELECT 1 FROM post_likes WHERE post_id = p.id AND user_id = ?) AS liked
		FROM posts p
		JOIN users u ON p.author_id = u.id
		WHERE p.id = ?
	`, GetUserID(r), id).Scan(
		&post.ID,
		&post.Title,
		&post.Content,
		&post.CreatedAt,
		&post.Author,
		&categoryStr,
		&post.CommentCount,
		&post.LikeCount,
		&liked,
	)
	post.Liked = liked != 0

	if errors.Is(err, sql.ErrNoRows) {
		utils.WriteError(w, http.StatusNotFound, "post not found")
		return
	}

	if err != nil {
		log.Printf("[ERROR] [GetPost DB Query]: %v", err)
		utils.WriteError(w, http.StatusInternalServerError, "could not fetch post")
		return
	}

	if categoryStr != "" {
		post.Categories = utils.SplitAndTrim(categoryStr, ",")
	} else {
		post.Categories = []string{}
	}

	if err := utils.WriteJSON(w, http.StatusOK, post); err != nil {
		log.Printf("[ERROR] [GetPost Response JSON]: %v", err)
	}
}

func GetPosts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	category := r.URL.Query().Get("category")

	query := `
		SELECT
			p.id,
			p.title,
			p.content,
			p.created_at,
			u.nickname,
			COALESCE((
				SELECT GROUP_CONCAT(c.name, ',')
				FROM post_categories pc
				JOIN categories c ON pc.category_id = c.id
				WHERE pc.post_id = p.id
			), '') AS categories,
			(SELECT COUNT(*) FROM comments WHERE post_id = p.id) AS comment_count,
			(SELECT COUNT(*) FROM post_likes WHERE post_id = p.id) AS like_count,
			EXISTS(SELECT 1 FROM post_likes WHERE post_id = p.id AND user_id = ?) AS liked
		FROM posts p
		JOIN users u ON p.author_id = u.id
	`

	args := []interface{}{GetUserID(r)}

	if category != "" {
		query += `
		WHERE EXISTS (
			SELECT 1 FROM post_categories pc2
			JOIN categories c2 ON pc2.category_id = c2.id
			WHERE pc2.post_id = p.id AND c2.name = ?
		)
		`
		args = append(args, category)
	}

	query += `
		ORDER BY p.created_at DESC
		LIMIT 50
	`

	rows, err := db.DB.Query(query, args...)
	if err != nil {
		log.Printf("[ERROR] [GetPosts DB Query]: %v", err)
		utils.WriteError(w, http.StatusInternalServerError, "could not fetch posts")
		return
	}
	defer rows.Close()

	posts := []postResponse{}

	for rows.Next() {
		var post postResponse
		var categoryStr string
		var liked int

		if err := rows.Scan(
			&post.ID,
			&post.Title,
			&post.Content,
			&post.CreatedAt,
			&post.Author,
			&categoryStr,
			&post.CommentCount,
			&post.LikeCount,
			&liked,
		); err != nil {
			log.Printf("[ERROR] [GetPosts Scan]: %v", err)
			utils.WriteError(w, http.StatusInternalServerError, "error processing post data")
			return
		}
		post.Liked = liked != 0

		if categoryStr != "" {
			post.Categories = utils.SplitAndTrim(categoryStr, ",")
		} else {
			post.Categories = []string{}
		}

		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		log.Printf("[ERROR] [GetPosts Rows Loop]: %v", err)
		utils.WriteError(w, http.StatusInternalServerError, "error reading posts stream")
		return
	}

	if err := utils.WriteJSON(w, http.StatusOK, posts); err != nil {
		log.Printf("[ERROR] [GetPosts Response JSON]: %v", err)
	}
}

type createPostPayload struct {
	Title      string   `json:"title"`
	Content    string   `json:"content"`
	Categories []string `json:"categories"`
}

func CreatePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1024*1024)
	defer r.Body.Close()

	var payload createPostPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		var maxBytesErr *http.MaxBytesError
		if errors.As(err, &maxBytesErr) {
			utils.WriteError(w, http.StatusRequestEntityTooLarge, "request body too large")
			return
		}
		utils.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	payload.Title = strings.TrimSpace(payload.Title)
	payload.Content = strings.TrimSpace(payload.Content)

	if payload.Title == "" {
		utils.WriteError(w, http.StatusBadRequest, "title is required")
		return
	}
	if payload.Content == "" {
		utils.WriteError(w, http.StatusBadRequest, "content is required")
		return
	}
	if len(payload.Categories) == 0 {
		utils.WriteError(w, http.StatusBadRequest, "at least one category is required")
		return
	}

	// Get authenticated user ID from context
	authorID := GetUserID(r)
	if authorID == "" {
		utils.WriteError(w, http.StatusUnauthorized, "authentication required")
		return
	}

	postID, err := utils.NewUUID()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "could not generate post ID")
		return
	}

	// Insert post and categories atomically
	tx, err := db.DB.Begin()
	if err != nil {
		log.Printf("[ERROR] [CreatePost Begin TX]: %v", err)
		utils.WriteError(w, http.StatusInternalServerError, "could not create post")
		return
	}

	_, err = tx.Exec(
		`INSERT INTO posts (id, author_id, title, content) VALUES (?, ?, ?, ?)`,
		postID, authorID, payload.Title, payload.Content,
	)
	if err != nil {
		tx.Rollback()
		log.Printf("[ERROR] [CreatePost Insert Post]: %v", err)
		utils.WriteError(w, http.StatusInternalServerError, "could not create post")
		return
	}

	for _, categoryID := range payload.Categories {
		if strings.TrimSpace(categoryID) == "" {
			continue
		}
		_, err = tx.Exec(
			`INSERT INTO post_categories (post_id, category_id) VALUES (?, ?)`,
			postID, strings.TrimSpace(categoryID),
		)
		if err != nil {
			tx.Rollback()
			log.Printf("[ERROR] [CreatePost Insert Category]: %v", err)
			utils.WriteError(w, http.StatusBadRequest, "invalid category")
			return
		}
	}

	if err := tx.Commit(); err != nil {
		log.Printf("[ERROR] [CreatePost Commit TX]: %v", err)
		utils.WriteError(w, http.StatusInternalServerError, "could not create post")
		return
	}

	if err := utils.WriteJSON(w, http.StatusCreated, map[string]string{
		"message": "post created successfully",
		"id":      postID,
	}); err != nil {
		log.Printf("[ERROR] [CreatePost Response JSON]: %v", err)
	}
}
