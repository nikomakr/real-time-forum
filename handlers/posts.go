package handlers

import (
	"log"
	"net/http"

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
	CreatedAt    string   `json:"created_at"`
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
			(SELECT COUNT(*) FROM comments WHERE post_id = p.id) AS comment_count
		FROM posts p
		JOIN users u ON p.author_id = u.id
	`

	var args []interface{}

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

		if err := rows.Scan(
			&post.ID,
			&post.Title,
			&post.Content,
			&post.CreatedAt,
			&post.Author,
			&categoryStr,
			&post.CommentCount,
		); err != nil {
			log.Printf("[ERROR] [GetPosts Scan]: %v", err)
			utils.WriteError(w, http.StatusInternalServerError, "error processing post data")
			return
		}

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