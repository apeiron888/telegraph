package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/google/uuid"
)

// RefreshTokenManager manages refresh tokens
type RefreshTokenManager struct {
	repo RefreshTokenRepo
	ttl  time.Duration
}

func NewRefreshTokenManager(repo RefreshTokenRepo, ttl time.Duration) *RefreshTokenManager {
	return &RefreshTokenManager{repo: repo, ttl: ttl}
}

func (m *RefreshTokenManager) Generate(ctx context.Context, userID uuid.UUID) (string, error) {
	// Generate random token
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", err
	}
	token := hex.EncodeToString(tokenBytes)
	
	// Hash for storage
	tokenHash := hashToken(token)
	
	expiresAt := time.Now().Add(m.ttl)
	return token, m.repo.Create(ctx, userID, tokenHash, expiresAt)
}

func (m *RefreshTokenManager) Verify(ctx context.Context, token string) (uuid.UUID, error) {
	tokenHash := hashToken(token)
	
	refreshToken, err := m.repo.GetByHash(ctx, tokenHash)
	if err != nil || refreshToken == nil {
		return uuid.UUID{}, err
	}
	
	// Check expiry
	if time.Now().After(refreshToken.ExpiresAt) {
		return uuid.UUID{}, err
	}
	
	return refreshToken.UserID, nil
}

func (m *RefreshTokenManager) Revoke(ctx context.Context, token string) error {
	tokenHash := hashToken(token)
	return m.repo.Revoke(ctx, tokenHash)
}

func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}
