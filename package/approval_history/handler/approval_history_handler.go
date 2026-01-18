package handler

import (
	"github.com/gofiber/fiber/v2"

	"workflow-approval/package/approval_history/ports"
)

// ApprovalHistoryHandler handles HTTP requests for approval history operations
type ApprovalHistoryHandler struct {
	historyService ports.ApprovalHistoryService
}

// NewApprovalHistoryHandler creates a new ApprovalHistoryHandler instance
func NewApprovalHistoryHandler(historyService ports.ApprovalHistoryService) *ApprovalHistoryHandler {
	return &ApprovalHistoryHandler{
		historyService: historyService,
	}
}

// GetRequestHistory retrieves the approval history for a specific request
// GET /api/requests/:id/history
func (h *ApprovalHistoryHandler) GetRequestHistory(c *fiber.Ctx) error {
	requestID := c.Params("id")
	if requestID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"data":    nil,
			"error":   "Request ID is required",
		})
	}

	history, err := h.historyService.GetHistoryByRequestID(c.Context(), requestID)
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
			"request_id": requestID,
			"history":    history,
		},
		"error": nil,
	})
}
