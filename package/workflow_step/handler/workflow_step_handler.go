package handler

import (
	"github.com/gofiber/fiber/v2"

	"workflow-approval/package/workflow_step/domain/dto"
	"workflow-approval/package/workflow_step/ports"
)

// WorkflowStepHandler handles HTTP requests for workflow step operations
type WorkflowStepHandler struct {
	stepService ports.WorkflowStepService
}

// NewWorkflowStepHandler creates a new WorkflowStepHandler instance
func NewWorkflowStepHandler(stepService ports.WorkflowStepService) *WorkflowStepHandler {
	return &WorkflowStepHandler{
		stepService: stepService,
	}
}

// Routes defines all routes for workflow step module
// Mounts routes under /api/workflows/:id/steps
func (h *WorkflowStepHandler) Routes(group fiber.Router) {
	// POST /api/workflows/:id/steps - Create a new step
	group.Post("", h.Create)

	// GET /api/workflows/:id/steps - List all steps
	group.Get("", h.GetAll)

	// GET /api/workflows/:id/steps/:stepId - Get a specific step
	group.Get("/:stepId", h.Get)

	// PUT /api/workflows/:id/steps/:stepId - Update a step
	group.Put("/:stepId", h.Update)

	// DELETE /api/workflows/:id/steps/:stepId - Delete a step
	group.Delete("/:stepId", h.Delete)
}

// Create creates a new workflow step
// POST /workflows/:id/steps
func (h *WorkflowStepHandler) Create(c *fiber.Ctx) error {
	workflowID := c.Params("id")
	if workflowID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"data":    nil,
			"error":   "Workflow ID is required",
		})
	}

	var req dto.CreateStepRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"data":    nil,
			"error":   "Invalid request body",
		})
	}

	step, err := h.stepService.CreateStep(c.Context(), workflowID, req.Level, req.ActorID, req.Conditions)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"data":    nil,
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    dto.ToStepResponse(step),
		"error":   nil,
	})
}

// GetAll retrieves all steps for a workflow
// GET /workflows/:id/steps
func (h *WorkflowStepHandler) GetAll(c *fiber.Ctx) error {
	workflowID := c.Params("id")
	if workflowID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"data":    nil,
			"error":   "Workflow ID is required",
		})
	}

	steps, err := h.stepService.GetSteps(c.Context(), workflowID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"data":    nil,
			"error":   "Failed to get steps",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    dto.ToStepResponseList(steps),
		"error":   nil,
	})
}

// Get retrieves a step by ID
// GET /workflows/:id/steps/:stepId
func (h *WorkflowStepHandler) Get(c *fiber.Ctx) error {
	stepID := c.Params("stepId")
	if stepID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"data":    nil,
			"error":   "Step ID is required",
		})
	}

	step, err := h.stepService.GetStepByID(c.Context(), stepID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"data":    nil,
			"error":   "Step not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    dto.ToStepResponse(step),
		"error":   nil,
	})
}

// Delete deletes a step by ID
// DELETE /workflows/:id/steps/:stepId
func (h *WorkflowStepHandler) Delete(c *fiber.Ctx) error {
	stepID := c.Params("stepId")
	if stepID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"data":    nil,
			"error":   "Step ID is required",
		})
	}

	err := h.stepService.DeleteStep(c.Context(), stepID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"data":    nil,
			"error":   "Step not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    fiber.Map{"message": "Step deleted successfully"},
		"error":   nil,
	})
}

// Update updates a workflow step
// PUT /workflows/:id/steps/:stepId
func (h *WorkflowStepHandler) Update(c *fiber.Ctx) error {
	stepID := c.Params("stepId")
	if stepID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"data":    nil,
			"error":   "Step ID is required",
		})
	}

	var req dto.UpdateStepRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"data":    nil,
			"error":   "Invalid request body",
		})
	}

	step, err := h.stepService.UpdateStep(c.Context(), stepID, req.Level, req.ActorID, req.Conditions)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"data":    nil,
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    dto.ToStepResponse(step),
		"error":   nil,
	})
}
