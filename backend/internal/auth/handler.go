package auth

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/google/uuid"

	"telegraph/internal/users"
)

type Handler struct {
	userSvc users.UserService
	refresh *RefreshTokenManager
	jwt     *users.JWTManager
	mfa     *MFAManager
}

func NewHandler(userSvc users.UserService, refresh *RefreshTokenManager, jwt *users.JWTManager, mfa *MFAManager) *Handler {
	return &Handler{userSvc: userSvc, refresh: refresh, jwt: jwt, mfa: mfa}
}

func (h *Handler) Routes() http.Handler {
	r := chi.NewRouter()

	r.Post("/login", h.Login)
	r.Post("/refresh", h.Refresh)
	r.Post("/logout", h.Logout)
	r.Post("/mfa/send", h.SendMFA)
	r.Post("/mfa/verify", h.VerifyMFA)

	return r
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	json.NewDecoder(r.Body).Decode(&body)

	user, err := h.userSvc.Login(r.Context(), body.Email, body.Password)
	if err != nil {
		http.Error(w, "invalid_credentials", 401)
		return
	}

	// access token
	access, err := h.jwt.Generate(user.ID.String(), user.Role)
	if err != nil {
		http.Error(w, "jwt_error", 500)
		return
	}

	// refresh
	refreshToken, err := h.refresh.Generate(r.Context(), user.ID)
	if err != nil {
		http.Error(w, "refresh_error", 500)
		return
	}

	json.NewEncoder(w).Encode(map[string]any{
		"access_token":  access,
		"refresh_token": refreshToken,
		"user": map[string]any{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
	})
}

func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Token string `json:"refresh_token"`
	}
	json.NewDecoder(r.Body).Decode(&body)

	uid, err := h.refresh.Verify(r.Context(), body.Token)
	if err != nil {
		http.Error(w, "invalid_refresh", 401)
		return
	}

	// rotate
	_ = h.refresh.Revoke(r.Context(), body.Token)
	newToken, _ := h.refresh.Generate(r.Context(), uid)

	u, _ := h.userSvc.GetByID(r.Context(), uid)

	access, _ := h.jwt.Generate(u.ID.String(), u.Role)

	json.NewEncoder(w).Encode(map[string]any{
		"access_token":  access,
		"refresh_token": newToken,
	})
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Token string `json:"refresh_token"`
	}
	json.NewDecoder(r.Body).Decode(&body)
	_ = h.refresh.Revoke(r.Context(), body.Token)
	w.WriteHeader(204)
}

func (h *Handler) SendMFA(w http.ResponseWriter, r *http.Request) {
	var body struct {
		UserID string `json:"user_id"`
	}
	json.NewDecoder(r.Body).Decode(&body)

	uid, _ := uuid.Parse(body.UserID)
	u, err := h.userSvc.GetByID(r.Context(), uid)
	if err != nil {
		http.Error(w, "not_found", 404)
		return
	}

	_, err = h.mfa.SendOTP(r.Context(), u.ID, u.Email)
	if err != nil {
		http.Error(w, "send_failed", 500)
		return
	}

	w.WriteHeader(204)
}

func (h *Handler) VerifyMFA(w http.ResponseWriter, r *http.Request) {
	var body struct {
		UserID string `json:"user_id"`
		Code   string `json:"code"`
	}
	json.NewDecoder(r.Body).Decode(&body)

	uid, _ := uuid.Parse(body.UserID)
	err := h.mfa.VerifyOTP(r.Context(), uid, body.Code)
	if err != nil {
		http.Error(w, "invalid_code", 400)
		return
	}

	w.WriteHeader(204)
}
