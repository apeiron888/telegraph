package media

import (
	"encoding/json"
	"net/http"

	"telegraph/internal/users"
)

type Handler struct {
	media MediaHandler
}

func NewHandler(media MediaHandler) *Handler {
	return &Handler{media: media}
}

// UploadFile handles file uploads
func (h *Handler) UploadFile(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form (max 50MB)
	if err := r.ParseMultipartForm(50 << 20); err != nil {
		http.Error(w, "file too large", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "no file provided", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Get uploader ID from context
	uploaderID := users.UserIDFromContext(r.Context())

	// Upload file
	metadata, err := h.media.Upload(file, header.Filename, header.Header.Get("Content-Type"), uploaderID, header.Size)
	if err != nil {
		http.Error(w, "upload failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(metadata)
}

// ServeFile serves uploaded files
func (h *Handler) ServeFile(w http.ResponseWriter, r *http.Request) {
	fileID := r.URL.Query().Get("id")
	if fileID == "" {
		http.Error(w, "file id required", http.StatusBadRequest)
		return
	}

	filePath, err := h.media.Serve(fileID)
	if err != nil {
		http.Error(w, "file not found", http.StatusNotFound)
		return
	}

	http.ServeFile(w, r, filePath)
}
