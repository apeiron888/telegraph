package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/google/uuid"
)

type MFARepo interface {
	Store(ctx context.Context, uid uuid.UUID, code string, exp time.Time) error
	Find(ctx context.Context, uid uuid.UUID) (string, time.Time, error)
	Delete(ctx context.Context, uid uuid.UUID) error
}

type MFAManager struct {
	repo   MFARepo
	sender EmailSender
}

func NewMFAManager(repo MFARepo, sender EmailSender) *MFAManager {
	return &MFAManager{repo: repo, sender: sender}
}

func (m *MFAManager) SendOTP(ctx context.Context, uid uuid.UUID, email string) (string, error) {
	buf := make([]byte, 3) // 6 hex chars
	rand.Read(buf)
	code := hex.EncodeToString(buf)

	exp := time.Now().Add(10 * time.Minute)
	err := m.repo.Store(ctx, uid, code, exp)
	if err != nil {
		return "", err
	}

	err = m.sender.Send(email, "Your OTP Code", "Your code: "+code)
	if err != nil {
		return "", err
	}

	return code, nil
}

func (m *MFAManager) VerifyOTP(ctx context.Context, uid uuid.UUID, code string) error {
	storedCode, exp, err := m.repo.Find(ctx, uid)
	if err != nil {
		return err
	}

	if time.Now().After(exp) {
		_ = m.repo.Delete(ctx, uid)
		return errors.New("expired")
	}

	if storedCode != code {
		return errors.New("invalid_code")
	}

	return m.repo.Delete(ctx, uid)
}
