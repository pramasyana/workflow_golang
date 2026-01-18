package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	approvalHistoryPorts "workflow-approval/package/approval_history/ports"
	"workflow-approval/package/request/domain"
	"workflow-approval/package/request/domain/dto"
	reqPorts "workflow-approval/package/request/ports"
)

// RequestHandler handles HTTP requests for request operations
type RequestHandler struct {
	requestService reqPorts.RequestService
	historyService approvalHistoryPorts.ApprovalHistoryService
}

// NewRequestHandler creates a new RequestHandler instance
func NewRequestHandler(requestService reqPorts.RequestService, historyService approvalHistoryPorts.ApprovalHistoryService) *RequestHandler {
	return &RequestHandler{
		requestService: requestService,
		historyService: historyService,
	}
}

// Routes defines all routes for request module
// Mounts routes under /api/requests
func (h *RequestHandler) Routes(group fiber.Router) {
	// POST /api/requests - Create a new request
	group.Post("", h.Create)

	// GET /api/requests - List all requests with pagination
	group.Get("", h.List)

	// GET /api/requests/:id - Get a specific request
	group.Get("/:id", h.Get)

	// PUT /api/requests/:id - Update a request (only if PENDING)
	group.Put("/:id", h.Update)

	// POST /api/requests/:id/approve - Approve a request
	group.Post("/:id/approve", h.Approve)

	// POST /api/requests/:id/reject - Reject a request
	group.Post("/:id/reject", h.Reject)

	// GET /api/requests/:id/history - Get approval history for a request
	group.Get("/:id/history", h.GetHistory)

	// DELETE /api/requests/:id - Delete a request
	group.Delete("/:id", h.Delete)
}

// Create creates a new request
// POST /requests
func (h *RequestHandler) Create(c *fiber.Ctx) error {
	requesterID := c.Locals("user_id").(string)

	var req dto.CreateRequestRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"data":    nil,
			"error":   "Invalid request body",
		})
	}

	request, err := h.requestService.CreateRequest(c.Context(), req.WorkflowID, requesterID, req.Amount, req.Title, req.Description)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"data":    nil,
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    dto.ToRequestResponse(request),
		"error":   nil,
	})
}

// Get retrieves a request by ID
// GET /requests/:id
func (h *RequestHandler) Get(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"data":    nil,
			"error":   "Request ID is required",
		})
	}

	request, err := h.requestService.GetRequest(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"data":    nil,
			"error":   "Request not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    dto.ToRequestResponse(request),
		"error":   nil,
	})
}

// List retrieves a list of requests
// GET /requests
func (h *RequestHandler) List(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	statusStr := c.Query("status")
	var status *domain.RequestStatus
	if statusStr != "" {
		s := domain.RequestStatus(statusStr)
		status = &s
	}

	requests, total, err := h.requestService.ListRequests(c.Context(), page, limit, status)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"data":    nil,
			"error":   "Failed to list requests",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"requests": dto.ToRequestResponseList(requests),
			"total":    total,
			"page":     page,
			"limit":    limit,
		},
		"error": nil,
	})
}

// Update updates a request
// PUT /requests/:id
func (h *RequestHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"data":    nil,
			"error":   "Request ID is required",
		})
	}

	var req dto.UpdateRequestRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"data":    nil,
			"error":   "Invalid request body",
		})
	}

	request, err := h.requestService.UpdateRequest(c.Context(), id, req.Amount, req.Title, req.Description)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"data":    nil,
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    dto.ToRequestResponse(request),
		"error":   nil,
	})
}

// Approve approves a request
// POST /requests/:id/approve
func (h *RequestHandler) Approve(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"data":    nil,
			"error":   "Request ID is required",
		})
	}

	userID := c.Locals("user_id").(string)
	actorID := c.Locals("actor_id").(string)
	isAdmin := c.Locals("is_admin").(bool)

	request, err := h.requestService.Approve(c.Context(), id, userID, actorID, isAdmin)
	if err != nil {
		// Return 403 for unauthorized actor errors
		if err.Error() == "unauthorized actor: you are not the assigned approver for this step" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"success": false,
				"data":    nil,
				"error":   err.Error(),
			})
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"data":    nil,
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    dto.ToRequestResponse(request),
		"error":   nil,
	})
}

// Reject rejects a request
// POST /requests/:id/reject
func (h *RequestHandler) Reject(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"data":    nil,
			"error":   "Request ID is required",
		})
	}

	var req dto.RejectRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"data":    nil,
			"error":   "Invalid request body",
		})
	}

	userID := c.Locals("user_id").(string)
	actorID := c.Locals("actor_id").(string)
	isAdmin := c.Locals("is_admin").(bool)

	request, err := h.requestService.Reject(c.Context(), id, userID, actorID, isAdmin, req.Reason)
	if err != nil {
		// Return 403 for unauthorized actor errors
		if err.Error() == "unauthorized actor: you are not the assigned approver for this step" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"success": false,
				"data":    nil,
				"error":   err.Error(),
			})
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"data":    nil,
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    dto.ToRequestResponse(request),
		"error":   nil,
	})
}

// Delete deletes a request by ID
// DELETE /requests/:id
func (h *RequestHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"data":    nil,
			"error":   "Request ID is required",
		})
	}

	err := h.requestService.DeleteRequest(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"data":    nil,
			"error":   "Request not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    fiber.Map{"message": "Request deleted successfully"},
		"error":   nil,
	})
}

// GetHistory retrieves the approval history for a specific request
// GET /requests/:id/history
func (h *RequestHandler) GetHistory(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"data":    nil,
			"error":   "Request ID is required",
		})
	}

	history, err := h.historyService.GetHistoryByRequestID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"data":    nil,
			"error":   "Failed to get approval history",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"request_id": id,
			"history":    history,
		},
		"error": nil,
	})
}
