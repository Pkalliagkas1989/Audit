package utils

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/google/uuid"
)

// GenerateUUID generates a new UUID string
func GenerateUUID() string {
	return uuid.New().String()
}

// GenerateSessionToken creates a new session token/ID
func GenerateSessionToken() string {
<<<<<<< HEAD
	return uuid.New().String()
=======
	return GenerateUUID()
>>>>>>> 289cb6929aae576611050ad1261732033fc1f610
}

// CalculateSessionExpiry calculates the expiry time for a session
// Default session lifetime is 24 hours
func CalculateSessionExpiry() time.Time {
	return time.Now().Add(24 * time.Hour)
}

func GenerateCSRFToken() string {
	bytes := make([]byte, 32) // 256 bits of randomness
	if _, err := rand.Read(bytes); err != nil {
		// fallback (very rare)
		return GenerateUUID()
	}
	return hex.EncodeToString(bytes)
}
