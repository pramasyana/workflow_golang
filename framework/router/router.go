package router

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"workflow-approval/framework/middleware"
	actorHandler "workflow-approval/package/actor/handler"
	authHandler "workflow-approval/package/auth/handler"
	requestHandler "workflow-approval/package/request/handler"
	userHandler "workflow-approval/package/user/handler"
	workflowHandler "workflow-approval/package/workflow/handler"
	workflowStepHandler "workflow-approval/package/workflow_step/handler"
)

// Config holds the router configuration
type Config struct {
	JWTSecret           string
	AuthHandler         *authHandler.AuthHandler
	UserHandler         *userHandler.UserHandler
	WorkflowHandler     *workflowHandler.WorkflowHandler
	WorkflowStepHandler *workflowStepHandler.WorkflowStepHandler
	RequestHandler      *requestHandler.RequestHandler
	ActorHandler        *actorHandler.ActorHandler
}

// Setup configures the Fiber application with all routes
func Setup(cfg Config) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:      "Workflow Approval System",
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		ErrorHandler: customErrorHandler,
	})

	// =========================================
	// Global Middleware
	// =========================================
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${latency} ${method} ${path}\n",
	}))
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	}))

	// =========================================
	// Health Check (Public)
	// =========================================
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"success": true,
			"data":    fiber.Map{"status": "healthy"},
			"error":   nil,
		})
	})

	// =========================================
	// Authentication Routes (Public - No JWT Required)
	// =========================================
	auth := app.Group("/auth")
	cfg.AuthHandler.Routes(auth)

	// =========================================
	// Protected Routes (JWT Authentication Required)
	// =========================================
	jwtMiddleware := middleware.NewJWTMiddleware(cfg.JWTSecret)
	api := app.Group("/api", jwtMiddleware)

	// =========================================
	// Actor Routes
	// =========================================
	actors := api.Group("/actors")
	cfg.ActorHandler.Routes(actors)

	// =========================================
	// Workflow Routes
	// =========================================
	workflows := api.Group("/workflows")
	cfg.WorkflowHandler.Routes(workflows)

	// Nested Workflow Steps
	workflowSteps := workflows.Group("/:id/steps")
	cfg.WorkflowStepHandler.Routes(workflowSteps)

	// =========================================
	// Request Routes
	// =========================================
	requests := api.Group("/requests")
	cfg.RequestHandler.Routes(requests)

	// =========================================
	// User Routes
	// =========================================
	cfg.UserHandler.Routes(api)

	return app
}

// customErrorHandler handles errors globally
func customErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	message := "Internal server error"

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		message = e.Message
	}

	return c.Status(code).JSON(fiber.Map{
		"success": false,
		"data":    nil,
		"error":   message,
	})
}
