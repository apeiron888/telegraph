package messages

import (
	"encoding/json"
	"net/http"
	"strconv"

	"telegraph/internal/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type Handler struct {
	service MessageService
}

func NewHandler(service MessageService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Routes() chi.Router {
	r := chi.NewRouter()

	// Messages for a specific channel
	r.Post("/{channelId}/messages", h.SendMessage)
	r.Get("/{channelId}/messages", h.GetMessages)

	// Individual message operations
	r.Put("/messages/{id}", h.EditMessage)
	r.Delete("/messages/{id}", h.DeleteMessage)
	r.Post("/messages/{id}/read", h.MarkAsRead)
	r.Post("/messages/{id}/delivered", h.MarkAsDelivered)
	
	// Typing indicators
	r.Post("/{channelId}/typing", h.SendTyping)
	r.Get("/unread", h.GetUnreadCounts)

	return r
}

func (h *Handler) SendMessage(w http.ResponseWriter, r *http.Request) {
	channelID, err := uuid.Parse(chi.URLParam(r, "channelId"))
	if err != nil {
		respondError(w, "invalid_channel_id", http.StatusBadRequest)
		return
	}

	var req SendMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, "invalid_request", http.StatusBadRequest)
		return
	}

	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		respondError(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	message, err := h.service.SendMessage(r.Context(), req, user.ID, channelID)
	if err != nil {
		if err == ErrNotChannelMember {
			respondError(w, err.Error(), http.StatusForbidden)
			return
		}
		if err == ErrContentTooLarge || err == ErrInvalidContentType {
			respondError(w, err.Error(), http.StatusBadRequest)
			return
		}
		respondError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(w, message, http.StatusCreated)
}

func (h *Handler) GetMessages(w http.ResponseWriter, r *http.Request) {
	channelID, err := uuid.Parse(chi.URLParam(r, "channelId"))
	if err != nil {
		respondError(w, "invalid_channel_id", http.StatusBadRequest)
		return
	}

	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		respondError(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse pagination parameters
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	messages, err := h.service.GetMessages(r.Context(), channelID, user.ID, limit, offset)
	if err != nil {
		if err == ErrNotChannelMember {
			respondError(w, err.Error(), http.StatusForbidden)
			return
		}
		respondError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(w, messages, http.StatusOK)
}

func (h *Handler) DeleteMessage(w http.ResponseWriter, r *http.Request) {
	messageID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, "invalid_message_id", http.StatusBadRequest)
		return
	}

	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		respondError(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	if err := h.service.DeleteMessage(r.Context(), messageID, user.ID, user.Role); err != nil {
		if err == ErrMessageNotFound {
			respondError(w, "message_not_found", http.StatusNotFound)
			return
		}
		if err == ErrNotSender {
			respondError(w, err.Error(), http.StatusForbidden)
			return
		}
		respondError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(w, map[string]string{"message": "deleted"}, http.StatusOK)
}

func (h *Handler) EditMessage(w http.ResponseWriter, r *http.Request) {
	messageID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, "invalid_message_id", http.StatusBadRequest)
		return
	}

	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		respondError(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		Content []byte `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, "invalid_request", http.StatusBadRequest)
		return
	}

	if err := h.service.EditMessage(r.Context(), messageID, user.ID, req.Content); err != nil {
		respondError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(w, map[string]string{"message": "edited"}, http.StatusOK)
}

func (h *Handler) MarkAsRead(w http.ResponseWriter, r *http.Request) {
	messageID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, "invalid_message_id", http.StatusBadRequest)
		return
	}

	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		respondError(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	if err := h.service.MarkAsRead(r.Context(), messageID, user.ID); err != nil {
		respondError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(w, map[string]string{"status": "read"}, http.StatusOK)
}

func (h *Handler) MarkAsDelivered(w http.ResponseWriter, r *http.Request) {
	messageID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, "invalid_message_id", http.StatusBadRequest)
		return
	}

	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		respondError(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	if err := h.service.MarkAsDelivered(r.Context(), messageID, user.ID); err != nil {
		respondError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(w, map[string]string{"status": "delivered"}, http.StatusOK)
}

func (h *Handler) SendTyping(w http.ResponseWriter, r *http.Request) {
	channelID, err := uuid.Parse(chi.URLParam(r, "channelId"))
	if err != nil {
		respondError(w, "invalid_channel_id", http.StatusBadRequest)
		return
	}
	
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		respondError(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		Typing bool `json:"typing"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, "invalid_request", http.StatusBadRequest)
		return
	}

	if err := h.service.BroadcastTyping(r.Context(), user.ID, channelID, req.Typing); err != nil {
		respondError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(w, map[string]string{"status": "ok"}, http.StatusOK)
}

func (h *Handler) GetUnreadCounts(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		respondError(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	counts, err := h.service.GetUnreadCounts(r.Context(), user.ID)
	if err != nil {
		respondError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(w, counts, http.StatusOK)
}

// Helper functions
func respondError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func respondJSON(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
