package handler

import (
	"github.com/gofiber/fiber/v2"

	"workflow-approval/package/actor/domain/dto"
	"workflow-approval/package/actor/ports"
)

// ActorHandler handles HTTP requests for actor operations
type ActorHandler struct {
	actorService ports.ActorService
}

// NewActorHandler creates a new ActorHandler instance
func NewActorHandler(actorService ports.ActorService) *ActorHandler {
	return &ActorHandler{
		actorService: actorService,
	}
}

// Routes defines all routes for actor module
// Mounts routes under /api/actors
func (h *ActorHandler) Routes(group fiber.Router) {
	// GET /api/actors - Get all actors
	group.Get("/", h.GetAllActors)

	// GET /api/actors/:id - Get actor by ID
	group.Get("/:id", h.GetActorByID)

	// POST /api/actors - Create a new actor
	group.Post("/", h.CreateActor)

	// PUT /api/actors/:id - Update actor
	group.Put("/:id", h.UpdateActor)

	// DELETE /api/actors/:id - Delete actor
	group.Delete("/:id", h.DeleteActor)
}

// GetAllActors handles getting all actors
// GET /api/actors
func (h *ActorHandler) GetAllActors(c *fiber.Ctx) error {
	actors, err := h.actorService.GetAllActors(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"data":    nil,
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    dto.ToActorResponseList(actors),
		"error":   nil,
	})
}

// GetActorByID handles getting an actor by ID
// GET /api/actors/:id
func (h *ActorHandler) GetActorByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"data":    nil,
			"error":   "Actor ID is required",
		})
	}

	actor, err := h.actorService.GetActorByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"data":    nil,
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    dto.ToActorResponse(actor),
		"error":   nil,
	})
}

// CreateActor handles creating a new actor
// POST /api/actors
func (h *ActorHandler) CreateActor(c *fiber.Ctx) error {
	var req dto.CreateActorRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"data":    nil,
			"error":   "Invalid request body",
		})
	}

	if req.Name == "" || req.Code == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"data":    nil,
			"error":   "Name and code are required",
		})
	}

	actor, err := h.actorService.CreateActor(c.Context(), req.Name, req.Code)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"data":    nil,
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    dto.ToActorResponse(actor),
		"error":   nil,
	})
}

// UpdateActor handles updating an actor
// PUT /api/actors/:id
func (h *ActorHandler) UpdateActor(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"data":    nil,
			"error":   "Actor ID is required",
		})
	}

	var req dto.UpdateActorRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"data":    nil,
			"error":   "Invalid request body",
		})
	}

	actor, err := h.actorService.UpdateActor(c.Context(), id, req.Name)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"data":    nil,
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    dto.ToActorResponse(actor),
		"error":   nil,
	})
}

// DeleteActor handles deleting an actor
// DELETE /api/actors/:id
func (h *ActorHandler) DeleteActor(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"data":    nil,
			"error":   "Actor ID is required",
		})
	}

	err := h.actorService.DeleteActor(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"data":    nil,
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"message": "Actor deleted successfully",
		},
		"error": nil,
	})
}
