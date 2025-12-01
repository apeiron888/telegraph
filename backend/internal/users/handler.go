package users

import (
	"encoding/json"
	"net/http"
	"strings"
	
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type Handler struct {
	svc UserService
	jwt *JWTManager
}

func NewHandler(svc UserService, jwt *JWTManager) *Handler {
	return &Handler{svc: svc, jwt: jwt}
}

func (h *Handler) Routes() http.Handler {
	r := chi.NewRouter()
	r.Post("/register", h.Register)
	return r
}

// Register
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "bad_request", 400)
		return
	}

	u := &User{
		Username: body.Username,
		Email:    strings.TrimSpace(body.Email),
		Bio:      "",
	}

	if err := h.svc.Register(r.Context(), u, body.Password); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	json.NewEncoder(w).Encode(map[string]any{
		"id":       u.ID,
		"username": u.Username,
		"email":    u.Email,
	})
}

// Me
func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	userID := UserIDFromContext(r.Context())
	id, _ := uuid.Parse(userID)

	u, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, "not_found", 404)
		return
	}

	json.NewEncoder(w).Encode(u)
}

// GetUser
func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	id, _ := uuid.Parse(chi.URLParam(r, "id"))
	u, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, "not_found", 404)
		return
	}
	json.NewEncoder(w).Encode(u)
}

// Update
func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id, _ := uuid.Parse(chi.URLParam(r, "id"))

	var body struct{ Bio string `json:"bio"` }
	json.NewDecoder(r.Body).Decode(&body)

	u, _ := h.svc.GetByID(r.Context(), id)
	u.Bio = body.Bio

	if err := h.svc.UpdateUser(r.Context(), u); err != nil {
		http.Error(w, "update_failed", 400)
		return
	}

	json.NewEncoder(w).Encode(u)
}

// Delete
func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, _ := uuid.Parse(chi.URLParam(r, "id"))
	err := h.svc.DeleteUser(r.Context(), id)
	if err != nil {
		http.Error(w, "not_found", 404)
		return
	}
	w.WriteHeader(204)
}
