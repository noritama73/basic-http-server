package util

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

// JWT related errors
var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token has expired")
)

// JWTClaims represents the claims in a JWT
type JWTClaims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	// Standard claims
	ExpiresAt int64 `json:"exp"`
	IssuedAt  int64 `json:"iat"`
}

// JWTManager handles JWT operations
type JWTManager struct {
	secretKey []byte
	expiry    time.Duration
}

// NewJWTManager creates a new JWTManager
func NewJWTManager(secretKey string, expiry time.Duration) *JWTManager {
	return &JWTManager{
		secretKey: []byte(secretKey),
		expiry:    expiry,
	}
}

// Generate creates a new JWT token for a user
func (m *JWTManager) Generate(userID, username string) (string, error) {
	now := time.Now()
	claims := JWTClaims{
		UserID:    userID,
		Username:  username,
		ExpiresAt: now.Add(m.expiry).Unix(),
		IssuedAt:  now.Unix(),
	}

	// Create the JWT header
	header := map[string]string{
		"alg": "HS256",
		"typ": "JWT",
	}

	// Encode header
	headerJSON, err := json.Marshal(header)
	if err != nil {
		return "", err
	}
	headerEncoded := base64.RawURLEncoding.EncodeToString(headerJSON)

	// Encode payload (claims)
	payloadJSON, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}
	payloadEncoded := base64.RawURLEncoding.EncodeToString(payloadJSON)

	// Create signature
	signatureInput := headerEncoded + "." + payloadEncoded
	h := hmac.New(sha256.New, m.secretKey)
	h.Write([]byte(signatureInput))
	signature := h.Sum(nil)
	signatureEncoded := base64.RawURLEncoding.EncodeToString(signature)

	// Combine to form the complete JWT
	token := headerEncoded + "." + payloadEncoded + "." + signatureEncoded
	return token, nil
}

// Verify checks if the token is valid and returns the claims
func (m *JWTManager) Verify(token string) (*JWTClaims, error) {
	// Split the token into parts
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, ErrInvalidToken
	}

	// Verify signature
	signatureInput := parts[0] + "." + parts[1]
	h := hmac.New(sha256.New, m.secretKey)
	h.Write([]byte(signatureInput))
	signature := h.Sum(nil)
	expectedSignature := base64.RawURLEncoding.EncodeToString(signature)

	if parts[2] != expectedSignature {
		return nil, ErrInvalidToken
	}

	// Decode payload
	payloadJSON, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, fmt.Errorf("failed to decode payload: %w", err)
	}

	var claims JWTClaims
	if err := json.Unmarshal(payloadJSON, &claims); err != nil {
		return nil, fmt.Errorf("failed to parse claims: %w", err)
	}

	// Check expiration
	if time.Now().Unix() > claims.ExpiresAt {
		return nil, ErrExpiredToken
	}

	return &claims, nil
}
