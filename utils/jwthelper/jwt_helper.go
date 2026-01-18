package jwthelper

import (
	"encoding/base64"
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"

	"workflow-approval/package/user/domain"
)

// JWTClaims represents the JWT claims structure
type JWTClaims struct {
	UserID  string `json:"user_id"`
	Email   string `json:"email"`
	IsAdmin bool   `json:"is_admin"`
	ActorID string `json:"actor_id"`
	jwt.RegisteredClaims
}

// JWTHelper provides JWT utility functions
type JWTHelper struct {
	secret     string
	expiration time.Duration
}

// NewJWTHelper creates a new JWTHelper instance
func NewJWTHelper(secret string, expiration time.Duration) *JWTHelper {
	return &JWTHelper{
		secret:     secret,
		expiration: expiration,
	}
}

// GetSecret returns the JWT secret
func (h *JWTHelper) GetSecret() string {
	return h.secret
}

// GetExpiration returns the JWT expiration duration
func (h *JWTHelper) GetExpiration() time.Duration {
	return h.expiration
}

// ParseBasicAuth parses the Authorization header for Basic authentication
// Returns email, password, and error
func ParseBasicAuth(authHeader string) (string, string, error) {
	if authHeader == "" {
		return "", "", errors.New("authorization header is required")
	}

	// Check for Basic prefix
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "basic" {
		return "", "", errors.New("invalid authorization header format")
	}

	// Decode base64 credentials
	decoded, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return "", "", errors.New("invalid base64 encoding")
	}

	// Split email:password
	credentialParts := strings.SplitN(string(decoded), ":", 2)
	if len(credentialParts) != 2 {
		return "", "", errors.New("invalid credentials format")
	}

	email := credentialParts[0]
	password := credentialParts[1]

	if email == "" || password == "" {
		return "", "", errors.New("email and password are required")
	}

	return email, password, nil
}

// GenerateJWT generates a JWT token for the user
func (h *JWTHelper) GenerateJWT(user *domain.User) (string, error) {
	actorID := ""
	if user.ActorID != nil {
		actorID = *user.ActorID
	}

	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"email":    user.Email,
		"is_admin": user.IsAdmin,
		"actor_id": actorID,
		"exp":      time.Now().Add(h.expiration).Unix(),
		"iat":      time.Now().Unix(),
		"iss":      "workflow-approval-system",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(h.secret))
}

// ValidateToken validates a JWT token and returns the claims
func (h *JWTHelper) ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(h.secret), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}

// ValidateTokenMapClaims validates a JWT token and returns map claims
func (h *JWTHelper) ValidateTokenMapClaims(tokenString string) (map[string]interface{}, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(h.secret), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}

// ExtractToken extracts the JWT token from the Authorization header
func ExtractToken(authHeader string) (string, error) {
	if authHeader == "" {
		return "", errors.New("authorization header is required")
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return "", errors.New("invalid authorization header format. Expected: Bearer <token>")
	}

	return parts[1], nil
}

// GlobalJWTHelper instance for application-wide use
var globalJWTHelper *JWTHelper

// InitJWTHelper initializes the global JWT helper
func InitJWTHelper(secret string, expiration time.Duration) {
	globalJWTHelper = NewJWTHelper(secret, expiration)
}

// GetJWTHelper returns the global JWT helper instance
func GetJWTHelper() *JWTHelper {
	return globalJWTHelper
}

// Global helper functions that use the global instance

// GenerateToken generates a JWT token using the global helper
func GenerateToken(user *domain.User) (string, error) {
	if globalJWTHelper == nil {
		return "", errors.New("JWT helper not initialized")
	}
	return globalJWTHelper.GenerateJWT(user)
}

// ValidateToken validates a token using the global helper
func ValidateToken(tokenString string) (*JWTClaims, error) {
	if globalJWTHelper == nil {
		return nil, errors.New("JWT helper not initialized")
	}
	return globalJWTHelper.ValidateToken(tokenString)
}

// ParseBasicAuthHeader parses Basic Auth using the global helper
func ParseBasicAuthHeader(authHeader string) (string, string, error) {
	return ParseBasicAuth(authHeader)
}
