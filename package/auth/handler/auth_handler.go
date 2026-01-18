package handler

import (
	"github.com/gofiber/fiber/v2"

	"workflow-approval/package/auth/domain/dto"
	"workflow-approval/package/auth/ports"
	"workflow-approval/utils/jwthelper"
)

// AuthHandler handles HTTP requests for authentication operations
type AuthHandler struct {
	authService ports.AuthService
	jwtHelper   *jwthelper.JWTHelper
}

// NewAuthHandler creates a new AuthHandler instance
func NewAuthHandler(authService ports.AuthService, jwtHelper *jwthelper.JWTHelper) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		jwtHelper:   jwtHelper,
	}
}

// Routes defines all routes for authentication module
// Mounts routes under /auth
func (h *AuthHandler) Routes(group fiber.Router) {
	// POST /auth/login - Login user and get JWT token (Request Body)
	group.Post("/login", h.Login)

	// POST /auth/refresh - Refresh JWT token
	group.Post("/refresh", h.Refresh)

	// POST /auth/logout - Logout user (invalidate token)
	group.Post("/logout", h.Logout)
}

// Login handles user login using request body
// POST /auth/login
// Request Body: {"email": "user@example.com", "password": "password"}
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	// Parse request body
	var req dto.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"data":    nil,
			"error":   "Invalid request body",
		})
	}

	// Validate input
	if req.Email == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"data":    nil,
			"error":   "Email and password are required",
		})
	}

	user, _, err := h.authService.Login(c.Context(), req.Email, req.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"data":    nil,
			"error":   "Invalid email or password",
		})
	}

	// Generate JWT token using helper
	token, err := h.jwtHelper.GenerateJWT(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"data":    nil,
			"error":   "Failed to generate token",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data": dto.AuthResponse{
			User:  dto.ToUserResponse(user),
			Token: token,
		},
		"error": nil,
	})
}

// Refresh handles token refresh
// POST /auth/refresh
// Authorization: Bearer <token>
func (h *AuthHandler) Refresh(c *fiber.Ctx) error {
	// Get the current token from Authorization header
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"data":    nil,
			"error":   "Authorization header is required",
		})
	}

	// Extract Bearer token
	tokenString, err := jwthelper.ExtractToken(authHeader)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"data":    nil,
			"error":   err.Error(),
		})
	}

	// Validate token and get claims
	claims, err := h.jwtHelper.ValidateTokenMapClaims(tokenString)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"data":    nil,
			"error":   err.Error(),
		})
	}

	// Get user from claims
	userID, ok := claims["user_id"].(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"data":    nil,
			"error":   "Invalid token claims",
		})
	}

	// Refresh token using auth service
	_, newToken, err := h.authService.Refresh(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"data":    nil,
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"token": newToken,
		},
		"error": nil,
	})
}

// Logout handles user logout
// POST /auth/logout
// Note: This is a client-side logout. For server-side token invalidation,
// you would need to implement a token blacklist (e.g., using Redis)
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	// Get the token from Authorization header
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"data":    nil,
			"error":   "Authorization header is required",
		})
	}

	// Extract Bearer token
	tokenString, err := jwthelper.ExtractToken(authHeader)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"data":    nil,
			"error":   err.Error(),
		})
	}

	// Validate token and logout using auth service
	err = h.authService.Logout(c.Context(), tokenString)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"data":    nil,
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"message": "Successfully logged out",
		},
		"error": nil,
	})
}
