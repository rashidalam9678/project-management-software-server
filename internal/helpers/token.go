package helpers

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// GenerateToken generates a unique token
func GenerateToken() (string, error) {
    // Generate random bytes
    tokenBytes := make([]byte, 16) // Adjust the length as needed
    if _, err := rand.Read(tokenBytes); err != nil {
        return "", err
    }

    // Convert bytes to hex
    tokenHex := hex.EncodeToString(tokenBytes)

    // Hash the token for added security
    hashedToken := sha256.Sum256([]byte(tokenHex))

    // Convert the hashed token to a string
    token := fmt.Sprintf("%x", hashedToken)

    return token, nil
}