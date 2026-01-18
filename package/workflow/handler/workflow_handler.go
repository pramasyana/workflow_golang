package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"workflow-approval/package/workflow/domain/dto"
	"workflow-approval/package/workflow/ports"
)

// WorkflowHandler handles HTTP requests for workflow operations
type WorkflowHandler struct {
	workflowService ports.WorkflowService
}

// NewWorkflowHandler creates a new WorkflowHandler instance
func NewWorkflowHandler(workflowService ports.WorkflowService) *WorkflowHandler {
	return &WorkflowHandler{
		workflowService: workflowService,
	}
}

// Routes defines all routes for workflow module
// Mounts routes under /api/workflows
func (h *WorkflowHandler) Routes(group fiber.Router) {
	// POST /api/workflows - Create a new workflow
	// Request body: { "name": "Workflow Name" }
	group.Post("", h.Create)

	// GET /api/workflows - List all workflows with pagination
	// Query params: page (default: 1), limit (default: 10)
	group.Get("", h.List)

	// GET /api/workflows/:id - Get a specific workflow by ID
	group.Get("/:id", h.Get)

	// PUT /api/workflows/:id - Update a workflow by ID
	group.Put("/:id", h.Update)

	// DELETE /api/workflows/:id - Delete a workflow by ID
	group.Delete("/:id", h.Delete)
}

// Create creates a new workflow
// POST /workflows
func (h *WorkflowHandler) Create(c *fiber.Ctx) error {
	var req dto.CreateWorkflowRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"data":    nil,
			"error":   "Invalid request body",
		})
	}

	workflow, err := h.workflowService.CreateWorkflow(c.Context(), req.Name)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"data":    nil,
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    dto.ToWorkflowResponse(workflow),
		"error":   nil,
	})
}

// Get retrieves a workflow by ID
// GET /workflows/:id
func (h *WorkflowHandler) Get(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"data":    nil,
			"error":   "Workflow ID is required",
		})
	}

	workflow, err := h.workflowService.GetWorkflow(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"data":    nil,
			"error":   "Workflow not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    dto.ToWorkflowResponse(workflow),
		"error":   nil,
	})
}

// List retrieves a list of workflows
// GET /workflows
func (h *WorkflowHandler) List(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	workflows, total, err := h.workflowService.ListWorkflows(c.Context(), page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"data":    nil,
			"error":   "Failed to list workflows",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"workflows": dto.ToWorkflowResponseList(workflows),
			"total":     total,
			"page":      page,
			"limit":     limit,
		},
		"error": nil,
	})
}

// Delete deletes a workflow by ID
// DELETE /workflows/:id
func (h *WorkflowHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"data":    nil,
			"error":   "Workflow ID is required",
		})
	}

	err := h.workflowService.DeleteWorkflow(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"data":    nil,
			"error":   "Workflow not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    fiber.Map{"message": "Workflow deleted successfully"},
		"error":   nil,
	})
}

// Update updates a workflow by ID
// PUT /workflows/:id
func (h *WorkflowHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"data":    nil,
			"error":   "Workflow ID is required",
		})
	}

	var req dto.UpdateWorkflowRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"data":    nil,
			"error":   "Invalid request body",
		})
	}

	workflow, err := h.workflowService.UpdateWorkflow(c.Context(), id, req.Name)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"data":    nil,
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    dto.ToWorkflowResponse(workflow),
		"error":   nil,
	})
}
