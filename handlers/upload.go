package handlers

import (
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"real-time-forum/utils"
)

const (
	maxImageBytes = 5 << 20 // 5MB
	uploadDir     = "static/uploads"
)

var allowedImageExtensions = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
	"image/gif":  ".gif",
	"image/webp": ".webp",
}

// UploadImage accepts a multipart image upload and returns its public URL — POST /api/messages/upload
func UploadImage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	if GetUserID(r) == "" {
		utils.WriteError(w, http.StatusUnauthorized, "authentication required")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxImageBytes)
	if err := r.ParseMultipartForm(maxImageBytes); err != nil {
		var maxBytesErr *http.MaxBytesError
		if errors.As(err, &maxBytesErr) {
			utils.WriteError(w, http.StatusRequestEntityTooLarge, "image too large (max 5MB)")
			return
		}
		utils.WriteError(w, http.StatusBadRequest, "invalid image upload")
		return
	}

	file, _, err := r.FormFile("image")
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "image file is required")
		return
	}
	defer file.Close()

	// Sniff the actual content type from the file bytes — never trust the client's declared type
	head := make([]byte, 512)
	n, err := file.Read(head)
	if err != nil && err != io.EOF {
		utils.WriteError(w, http.StatusBadRequest, "could not read image")
		return
	}
	head = head[:n]

	ext, ok := allowedImageExtensions[http.DetectContentType(head)]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, "unsupported image type — only JPEG, PNG, GIF, and WEBP are allowed")
		return
	}

	id, err := utils.NewUUID()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "could not generate image ID")
		return
	}

	if err := os.MkdirAll(uploadDir, 0o755); err != nil {
		log.Printf("[ERROR] [UploadImage MkdirAll]: %v", err)
		utils.WriteError(w, http.StatusInternalServerError, "could not save image")
		return
	}

	filename := id + ext
	dst, err := os.Create(filepath.Join(uploadDir, filename))
	if err != nil {
		log.Printf("[ERROR] [UploadImage Create]: %v", err)
		utils.WriteError(w, http.StatusInternalServerError, "could not save image")
		return
	}
	defer dst.Close()

	if _, err := dst.Write(head); err != nil {
		log.Printf("[ERROR] [UploadImage Write Header]: %v", err)
		utils.WriteError(w, http.StatusInternalServerError, "could not save image")
		return
	}
	if _, err := io.Copy(dst, file); err != nil {
		log.Printf("[ERROR] [UploadImage Write Body]: %v", err)
		utils.WriteError(w, http.StatusInternalServerError, "could not save image")
		return
	}

	if err := utils.WriteJSON(w, http.StatusCreated, map[string]string{
		"image_url": "/uploads/" + filename,
	}); err != nil {
		log.Printf("[ERROR] [UploadImage Response JSON]: %v", err)
	}
}
