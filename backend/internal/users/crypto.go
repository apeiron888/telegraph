package users

import (
	"crypto/rand"
	"encoding/base64"
	"strings"
	"crypto/subtle"


	"golang.org/x/crypto/argon2"
)

func HashPassword(pw string) string {
	salt := make([]byte, 16)
	rand.Read(salt)

	hash := argon2.IDKey([]byte(pw), salt, 1, 64*1024, 4, 32)

	return base64.RawStdEncoding.EncodeToString(salt) + ":" +
		base64.RawStdEncoding.EncodeToString(hash)
}

func VerifyPassword(stored, pw string) bool {
	parts := strings.Split(stored, ":")
	if len(parts) != 2 {
		return false
	}

	salt, _ := base64.RawStdEncoding.DecodeString(parts[0])
	hash, _ := base64.RawStdEncoding.DecodeString(parts[1])

	newHash := argon2.IDKey([]byte(pw), salt, 1, 64*1024, 4, 32)
	return subtle.ConstantTimeCompare(hash, newHash) == 1
}
