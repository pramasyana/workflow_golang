package handler

import (
	"github.com/gofiber/fiber/v2"
)

// SetupWorkflowStepRoutes configures all workflow step routes
// This function separates route definitions from router.go for better organization
func SetupWorkflowStepRoutes(app *fiber.App, workflowStepHandler *WorkflowStepHandler, jwtMiddleware fiber.Handler) {
	// Create workflow steps group nested under /api/workflows/:workflow_id/steps
	workflowSteps := app.Group("/api/workflows/:id/steps", jwtMiddleware)

	// POST /api/workflows/:id/steps - Create a new step for a workflow
	workflowSteps.Post("", workflowStepHandler.Create)

	// GET /api/workflows/:id/steps - List all steps for a workflow
	workflowSteps.Get("", workflowStepHandler.GetAll)

	// GET /api/workflows/:id/steps/:stepId - Get a specific step by ID
	workflowSteps.Get("/:stepId", workflowStepHandler.Get)

	// PUT /api/workflows/:id/steps/:stepId - Update a step
	workflowSteps.Put("/:stepId", workflowStepHandler.Update)

	// DELETE /api/workflows/:id/steps/:stepId - Delete a step by ID
	workflowSteps.Delete("/:stepId", workflowStepHandler.Delete)
}
