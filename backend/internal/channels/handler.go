package channels

import (
	"encoding/json"
	"net/http"

	"telegraph/internal/middleware"
	"telegraph/internal/users"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type Handler struct {
	service ChannelService
	userSvc users.UserService
}

func NewHandler(service ChannelService, userSvc users.UserService) *Handler {
	return &Handler{service: service, userSvc: userSvc}
}

func (h *Handler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Post("/", h.CreateChannel)
	r.Get("/", h.ListMyChannels)
	r.Get("/{id}", h.GetChannel)
	r.Post("/{id}/members", h.AddMember)
	r.Delete("/{id}/members/{userId}", h.RemoveMember)
	r.Post("/{id}/members/{userId}/promote", h.PromoteMember)
	r.Post("/{id}/members/{userId}/demote", h.DemoteMember)
	r.Delete("/{id}", h.DeleteChannel)

	return r
}

func (h *Handler) CreateChannel(w http.ResponseWriter, r *http.Request) {
	var req CreateChannelRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, "invalid_request: "+err.Error(), http.StatusBadRequest)
		return
	}

	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		respondError(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	channel, err := h.service.CreateChannel(r.Context(), req, user.ID, user.Role)
	if err != nil {
		if err == ErrBroadcastRestricted {
			respondError(w, err.Error(), http.StatusForbidden)
			return
		}
		respondError(w, err.Error(), http.StatusBadRequest)
		return
	}

	respondJSON(w, channel, http.StatusCreated)
}

func (h *Handler) GetChannel(w http.ResponseWriter, r *http.Request) {
	channelID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, "invalid_channel_id", http.StatusBadRequest)
		return
	}

	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		respondError(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	channel, err := h.service.GetChannel(r.Context(), channelID)
	if err != nil {
		if err == ErrChannelNotFound {
			respondError(w, "channel_not_found", http.StatusNotFound)
			return
		}
		respondError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if user is a member
	isMember, err := h.service.IsMember(r.Context(), channelID, user.ID)
	if err != nil {
		respondError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !isMember {
		respondError(w, "not_a_member", http.StatusForbidden)
		return
	}

	respondJSON(w, channel, http.StatusOK)
}

func (h *Handler) ListMyChannels(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		respondError(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	channels, err := h.service.GetUserChannels(r.Context(), user.ID)
	if err != nil {
		respondError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(w, channels, http.StatusOK)
}

func (h *Handler) AddMember(w http.ResponseWriter, r *http.Request) {
	channelID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, "invalid_channel_id", http.StatusBadRequest)
		return
	}

	var req AddMemberRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, "invalid_request", http.StatusBadRequest)
		return
	}

	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		respondError(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	memberID := req.UserID
	if memberID == uuid.Nil {
		// Try lookup by email or phone
		identifier := req.Email
		if identifier == "" {
			identifier = req.Phone
		}
		if identifier == "" {
			respondError(w, "user_id, email, or phone is required", http.StatusBadRequest)
			return
		}

		foundUser, err := h.userSvc.GetByEmailOrPhone(r.Context(), identifier)
		if err != nil {
			respondError(w, "user_not_found", http.StatusNotFound)
			return
		}
		memberID = foundUser.ID
	}

	if err := h.service.AddMember(r.Context(), channelID, user.ID, memberID); err != nil {
		if err == ErrNotChannelOwner {
			respondError(w, err.Error(), http.StatusForbidden)
			return
		}
		respondError(w, err.Error(), http.StatusBadRequest)
		return
	}

	respondJSON(w, map[string]string{"message": "member_added"}, http.StatusOK)
}

func (h *Handler) RemoveMember(w http.ResponseWriter, r *http.Request) {
	channelID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, "invalid_channel_id", http.StatusBadRequest)
		return
	}

	memberID, err := uuid.Parse(chi.URLParam(r, "userId"))
	if err != nil {
		respondError(w, "invalid_user_id", http.StatusBadRequest)
		return
	}

	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		respondError(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	if err := h.service.RemoveMember(r.Context(), channelID, user.ID, memberID); err != nil {
		if err == ErrNotChannelOwner {
			respondError(w, err.Error(), http.StatusForbidden)
			return
		}
		respondError(w, err.Error(), http.StatusBadRequest)
		return
	}

	respondJSON(w, map[string]string{"message": "member_removed"}, http.StatusOK)
}

func (h *Handler) DeleteChannel(w http.ResponseWriter, r *http.Request) {
	channelID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, "invalid_channel_id", http.StatusBadRequest)
		return
	}

	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		respondError(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	if err := h.service.DeleteChannel(r.Context(), channelID, user.ID); err != nil {
		if err == ErrNotChannelOwner {
			respondError(w, err.Error(), http.StatusForbidden)
			return
		}
		if err == ErrChannelNotFound {
			respondError(w, "channel_not_found", http.StatusNotFound)
			return
		}
		respondError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(w, map[string]string{"message": "channel_deleted"}, http.StatusOK)
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


func (h *Handler) PromoteMember(w http.ResponseWriter, r *http.Request) {
	channelID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, "invalid_channel_id", http.StatusBadRequest)
		return
	}

	memberID, err := uuid.Parse(chi.URLParam(r, "userId"))
	if err != nil {
		respondError(w, "invalid_user_id", http.StatusBadRequest)
		return
	}

	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		respondError(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	if err := h.service.PromoteToAdmin(r.Context(), channelID, user.ID, memberID); err != nil {
		if err == ErrNotChannelOwner {
			respondError(w, err.Error(), http.StatusForbidden)
			return
		}
		respondError(w, err.Error(), http.StatusBadRequest)
		return
	}

	respondJSON(w, map[string]string{"message": "member_promoted"}, http.StatusOK)
}

func (h *Handler) DemoteMember(w http.ResponseWriter, r *http.Request) {
	channelID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, "invalid_channel_id", http.StatusBadRequest)
		return
	}

	memberID, err := uuid.Parse(chi.URLParam(r, "userId"))
	if err != nil {
		respondError(w, "invalid_user_id", http.StatusBadRequest)
		return
	}

	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		respondError(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	if err := h.service.DemoteAdmin(r.Context(), channelID, user.ID, memberID); err != nil {
		if err == ErrNotChannelOwner {
			respondError(w, err.Error(), http.StatusForbidden)
			return
		}
		respondError(w, err.Error(), http.StatusBadRequest)
		return
	}

	respondJSON(w, map[string]string{"message": "member_demoted"}, http.StatusOK)
}