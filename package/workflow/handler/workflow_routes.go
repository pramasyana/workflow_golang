package handler

import (
	"github.com/gofiber/fiber/v2"
)

// SetupWorkflowRoutes configures all workflow routes
// This function separates route definitions from router.go for better organization
func SetupWorkflowRoutes(app *fiber.App, workflowHandler *WorkflowHandler, jwtMiddleware fiber.Handler) {
	// Create workflow group with JWT middleware
	workflows := app.Group("/api/workflows", jwtMiddleware)

	// POST /api/workflows - Create a new workflow
	workflows.Post("", workflowHandler.Create)

	// GET /api/workflows - List all workflows with pagination
	workflows.Get("", workflowHandler.List)

	// GET /api/workflows/:id - Get a specific workflow by ID
	workflows.Get("/:id", workflowHandler.Get)

	// PUT /api/workflows/:id - Update a workflow
	workflows.Put("/:id", workflowHandler.Update)

	// DELETE /api/workflows/:id - Delete a workflow by ID
	workflows.Delete("/:id", workflowHandler.Delete)
}
