package middleware

import (
	"errors"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

// JWTClaims represents the JWT claims structure
type JWTClaims struct {
	UserID  string `json:"user_id"`
	Email   string `json:"email"`
	IsAdmin bool   `json:"is_admin"`
	ActorID string `json:"actor_id"`
	jwt.RegisteredClaims
}

// JWTConfig holds the JWT middleware configuration
type JWTConfig struct {
	Secret string
}

// NewJWTMiddleware creates a new JWT middleware
func NewJWTMiddleware(secret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get the Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"data":    nil,
				"error":   "Authorization header is required",
				"code":    "MISSING_AUTH_HEADER",
			})
		}

		// Check for Bearer token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"data":    nil,
				"error":   "Invalid authorization header format. Expected: Bearer <token>",
				"code":    "INVALID_AUTH_FORMAT",
			})
		}

		tokenString := parts[1]

		// Check if token is empty
		if tokenString == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"data":    nil,
				"error":   "Token is empty",
				"code":    "EMPTY_TOKEN",
			})
		}

		// Parse and validate the token
		token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			// Validate the signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(secret), nil
		})

		// Handle different token validation errors
		if err != nil {
			return handleTokenError(c, err)
		}

		// Extract claims
		claims, ok := token.Claims.(*JWTClaims)
		if !ok || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"data":    nil,
				"error":   "Invalid token claims",
				"code":    "INVALID_CLAIMS",
			})
		}

		// Check token expiration manually (additional check)
		if claims.ExpiresAt != nil && claims.ExpiresAt.Before(time.Now()) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"data":    nil,
				"error":   "Token has expired",
				"code":    "TOKEN_EXPIRED",
			})
		}

		// Store user info in context locals
		c.Locals("user_id", claims.UserID)
		c.Locals("user_email", claims.Email)
		c.Locals("is_admin", claims.IsAdmin)
		c.Locals("actor_id", claims.ActorID)

		return c.Next()
	}
}

// handleTokenError provides detailed error messages for different token validation failures
func handleTokenError(c *fiber.Ctx, err error) error {
	var jwtErr *jwt.ValidationError

	// Check if it's a validation error
	if errors.As(err, &jwtErr) {
		// Token expired
		if jwtErr.Errors&jwt.ValidationErrorExpired != 0 {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"data":    nil,
				"error":   "Token has expired. Please refresh your token or login again.",
				"code":    "TOKEN_EXPIRED",
			})
		}

		// Token not yet valid
		if jwtErr.Errors&jwt.ValidationErrorNotValidYet != 0 {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"data":    nil,
				"error":   "Token is not yet valid",
				"code":    "TOKEN_NOT_YET_VALID",
			})
		}

		// Invalid signature
		if jwtErr.Errors&jwt.ValidationErrorSignatureInvalid != 0 {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"data":    nil,
				"error":   "Invalid token signature",
				"code":    "INVALID_SIGNATURE",
			})
		}

		// Malformed token
		if jwtErr.Errors&jwt.ValidationErrorMalformed != 0 {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"data":    nil,
				"error":   "Malformed token. Please provide a valid JWT token.",
				"code":    "MALFORMED_TOKEN",
			})
		}

		// General validation error
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"data":    nil,
			"error":   "Token validation failed: " + err.Error(),
			"code":    "VALIDATION_ERROR",
		})
	}

	// Handle other errors (e.g., unmarshaling errors)
	if err.Error() == "unexpected signing method" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"data":    nil,
			"error":   "Invalid signing method used for token",
			"code":    "INVALID_SIGNING_METHOD",
		})
	}

	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"success": false,
		"data":    nil,
		"error":   "Invalid or expired token",
		"code":    "INVALID_TOKEN",
	})
}

// ExtractToken extracts the JWT token from the Authorization header
// Returns the token string and nil error if successful
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

// GetUserIDFromContext retrieves the user ID from Fiber context locals
func GetUserIDFromContext(c *fiber.Ctx) string {
	if userID := c.Locals("user_id"); userID != nil {
		return userID.(string)
	}
	return ""
}

// GetUserEmailFromContext retrieves the user email from Fiber context locals
func GetUserEmailFromContext(c *fiber.Ctx) string {
	if email := c.Locals("user_email"); email != nil {
		return email.(string)
	}
	return ""
}

// GetIsAdminFromContext retrieves the is_admin flag from Fiber context locals
func GetIsAdminFromContext(c *fiber.Ctx) bool {
	if isAdmin := c.Locals("is_admin"); isAdmin != nil {
		return isAdmin.(bool)
	}
	return false
}

// GetActorIDFromContext retrieves the actor_id from Fiber context locals
func GetActorIDFromContext(c *fiber.Ctx) string {
	if actorID := c.Locals("actor_id"); actorID != nil {
		return actorID.(string)
	}
	return ""
}
