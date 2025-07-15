package handlers

import (
	"net/http"
	"os"
	"path/filepath"
)

type ImageHandler struct{}

func (h *ImageHandler) ServeServiceImage(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get(":filename")
	if filename == "" {
		http.Error(w, "filename is required", http.StatusBadRequest)
		return
	}

	imagePath := filepath.Join("cmd/uploads/reviews", filename)

	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		http.Error(w, "image not found", http.StatusNotFound)
		return
	}

	switch ext := filepath.Ext(imagePath); ext {
	case ".jpg", ".jpeg":
		w.Header().Set("Content-Type", "image/jpeg")
	case ".png":
		w.Header().Set("Content-Type", "image/png")
	case ".gif":
		w.Header().Set("Content-Type", "image/gif")
	default:
		w.Header().Set("Content-Type", "application/octet-stream")
	}

	http.ServeFile(w, r, imagePath)
}
