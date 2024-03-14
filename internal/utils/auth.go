package utils

import (
	"shopifyx/configs"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// GenerateAccessToken generates a JWT access token for the provided username.
func GenerateAccessToken(username string, userID string) (string, error) {
	config, err := configs.LoadConfig()
	if err != nil {
		return "", err
	}
	var (
		// Define a secret key for signing the JWT token.
		// Ensure to keep this key secure and don't expose it.
		// You may want to use an environment variable to store it.
		secretKey = []byte(config.JWTSecret)
	)
	// Define the token expiration time.
	expirationTime := time.Now().Add(2 * time.Minute) // 2 minutes for testing (adjust as needed)

	// Create a new token object with the appropriate claims.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"user_id":  userID,
		"exp":      expirationTime.Unix(),
	})

	// Sign the token with the secret key.
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
