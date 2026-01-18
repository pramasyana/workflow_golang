package handler

import (
	"github.com/gofiber/fiber/v2"

	"workflow-approval/package/user/domain/dto"
	"workflow-approval/package/user/ports"
)

// UserHandler handles HTTP requests for user operations
type UserHandler struct {
	userService ports.UserService
}

// NewUserHandler creates a new UserHandler instance
func NewUserHandler(userService ports.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// Create creates a new user
// POST /api/users
func (h *UserHandler) Create(c *fiber.Ctx) error {
	var req dto.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"data":    nil,
			"error":   "Invalid request body",
		})
	}

	user, err := h.userService.Register(c.Context(), req.Email, req.Password, req.Name, req.IsAdmin, req.ActorID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"data":    nil,
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    dto.ToUserResponse(user),
		"error":   nil,
	})
}

// Routes defines all routes for user module
// Mounts routes under /api
func (h *UserHandler) Routes(group fiber.Router) {
	// POST /api/users - Create a new user (Admin only in production)
	group.Post("/users", h.Create)

	// GET /api/profile - Get current user's profile
	group.Get("/profile", h.GetProfile)

	// PUT /api/profile - Update current user's profile
	group.Put("/profile", h.UpdateProfile)
}

// GetProfile handles getting the current user's profile
// GET /api/profile
func (h *UserHandler) GetProfile(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	user, err := h.userService.GetUserByID(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"data":    nil,
			"error":   "User not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    dto.ToUserResponse(user),
		"error":   nil,
	})
}

// UpdateProfile handles updating the current user's profile
// PUT /api/profile
func (h *UserHandler) UpdateProfile(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	var req dto.UpdateProfileRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"data":    nil,
			"error":   "Invalid request body",
		})
	}

	if req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"data":    nil,
			"error":   "Name is required",
		})
	}

	user, err := h.userService.UpdateUser(c.Context(), userID, req.Name)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"data":    nil,
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    dto.ToUserResponse(user),
		"error":   nil,
	})
}
