package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"workflow-approval/config"
	"workflow-approval/framework/router"
	actorHandler "workflow-approval/package/actor/handler"
	actorRepo "workflow-approval/package/actor/repository"
	actorUsecase "workflow-approval/package/actor/usecase"
	approvalHistoryRepo "workflow-approval/package/approval_history/repository"
	approvalHistoryUsecase "workflow-approval/package/approval_history/usecase"
	authHandler "workflow-approval/package/auth/handler"
	authRepo "workflow-approval/package/auth/repository"
	authUsecase "workflow-approval/package/auth/usecase"
	reqHandler "workflow-approval/package/request/handler"
	reqRepo "workflow-approval/package/request/repository"
	reqUsecase "workflow-approval/package/request/usecase"
	userHandler "workflow-approval/package/user/handler"
	userRepo "workflow-approval/package/user/repository"
	userUsecase "workflow-approval/package/user/usecase"
	wfHandler "workflow-approval/package/workflow/handler"
	wfRepo "workflow-approval/package/workflow/repository"
	wfUsecase "workflow-approval/package/workflow/usecase"
	stepHandler "workflow-approval/package/workflow_step/handler"
	stepRepo "workflow-approval/package/workflow_step/repository"
	stepUsecase "workflow-approval/package/workflow_step/usecase"
	"workflow-approval/utils/jwthelper"
)

func main() {
	// Load configuration
	cfg, err := config.Load("config/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database connection (MySQL)
	db, err := initDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Run migrations using raw SQL for MariaDB compatibility
	if err := runMigrations(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize repositories
	userRepository := userRepo.NewUserRepository(db)
	workflowRepository := wfRepo.NewWorkflowRepository(db)
	workflowStepRepository := stepRepo.NewWorkflowStepRepository(db)
	requestRepository := reqRepo.NewRequestRepository(db)
	actorRepository := actorRepo.NewActorRepository(db)
	approvalHistoryRepository := approvalHistoryRepo.NewApprovalHistoryRepository(db)

	// Initialize services
	userService := userUsecase.NewUserService(userRepository, actorRepository)
	workflowService := wfUsecase.NewWorkflowService(workflowRepository)
	workflowStepService := stepUsecase.NewWorkflowStepService(workflowStepRepository, actorRepository, workflowRepository)
	approvalHistoryService := approvalHistoryUsecase.NewApprovalHistoryService(approvalHistoryRepository)
	requestService := reqUsecase.NewRequestService(requestRepository, workflowRepository, workflowStepRepository, approvalHistoryRepository)
	actorService := actorUsecase.NewActorService(actorRepository)

	// Initialize auth services
	jwtExpiry := time.Duration(cfg.JWT.Expiration) * time.Hour
	jwtHelper := jwthelper.NewJWTHelper(cfg.JWT.Secret, jwtExpiry)
	authRepository := authRepo.NewAuthRepository(userRepository)
	authService := authUsecase.NewAuthService(authRepository, jwtHelper)

	// Initialize handlers
	authHTTPHandler := authHandler.NewAuthHandler(authService, jwtHelper)
	userHTTPHandler := userHandler.NewUserHandler(userService)
	workflowHTTPHandler := wfHandler.NewWorkflowHandler(workflowService)
	workflowStepHTTPHandler := stepHandler.NewWorkflowStepHandler(workflowStepService)
	requestHTTPHandler := reqHandler.NewRequestHandler(requestService, approvalHistoryService)
	actorHTTPHandler := actorHandler.NewActorHandler(actorService)

	// Setup router
	app := router.Setup(router.Config{
		JWTSecret:           cfg.JWT.Secret,
		AuthHandler:         authHTTPHandler,
		UserHandler:         userHTTPHandler,
		WorkflowHandler:     workflowHTTPHandler,
		WorkflowStepHandler: workflowStepHTTPHandler,
		RequestHandler:      requestHTTPHandler,
		ActorHandler:        actorHTTPHandler,
	})

	// Start server in a goroutine
	go func() {
		if err := app.Listen(cfg.App.Address()); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	log.Printf("Server started on %s", cfg.App.Address())

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	if err := app.Shutdown(); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}

	log.Println("Server exited")
}

// initDatabase initializes MySQL database connection
func initDatabase(cfg *config.Config) (*gorm.DB, error) {
	dsn := cfg.Database.DSN()

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MySQL database: %w", err)
	}

	// Get underlying SQL DB to configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}

	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.Database.GetConnMaxLifetime())

	return db, nil
}

// runMigrations runs database migrations using raw SQL for MariaDB compatibility
func runMigrations(db *gorm.DB) error {
	// Get database name from config - we need to create DB first if it doesn't exist
	dbName := "workflow_approval"

	// Create database if it doesn't exist
	createDBsql := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci", dbName)
	if err := db.Exec(createDBsql).Error; err != nil {
		return fmt.Errorf("failed to create database: %w", err)
	}

	// Switch to the target database
	useDBsql := fmt.Sprintf("USE %s", dbName)
	if err := db.Exec(useDBsql).Error; err != nil {
		return fmt.Errorf("failed to use database: %w", err)
	}

	log.Printf("Using database: %s", dbName)

	// Disable foreign key checks for migration
	db.Exec("SET FOREIGN_KEY_CHECKS=0")

	// Create actors table
	createActorsSQL := `
	CREATE TABLE IF NOT EXISTS actors (
		id VARCHAR(36) PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		code VARCHAR(50) NOT NULL UNIQUE,
		created_at DATETIME,
		updated_at DATETIME
	) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_general_ci
	`
	if err := db.Exec(createActorsSQL).Error; err != nil {
		return fmt.Errorf("failed to create actors table: %w", err)
	}

	// Create users table
	createUsersSQL := `
	CREATE TABLE IF NOT EXISTS users (
		id VARCHAR(36) PRIMARY KEY,
		email VARCHAR(255) NOT NULL UNIQUE,
		password VARCHAR(255) NOT NULL,
		name VARCHAR(255) NOT NULL,
		is_admin BOOLEAN DEFAULT FALSE,
		actor_id VARCHAR(36),
		created_at DATETIME,
		updated_at DATETIME,
		INDEX idx_actor_id (actor_id)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_general_ci
	`
	if err := db.Exec(createUsersSQL).Error; err != nil {
		return fmt.Errorf("failed to create users table: %w", err)
	}

	// Insert default admin user if not exists
	insertAdminSQL := `
	INSERT IGNORE INTO users (id, email, password, name, is_admin, actor_id, created_at, updated_at)
	VALUES ('b693fdee-f2bb-11f0-8cc1-7a447c8d071a', 'administrator@gmail.com', '$2a$10$qToIHURdilgF4kZCTnvR5.8AVRuljFLW1vHoRpAWjDIQjpPZyZ.ie', 'Administrator', 1, NULL, NOW(), NOW())
	`
	if err := db.Exec(insertAdminSQL).Error; err != nil {
		return fmt.Errorf("failed to insert admin user: %w", err)
	}

	// Create workflows table
	createWorkflowsSQL := `
	CREATE TABLE IF NOT EXISTS workflows (
		id VARCHAR(36) PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		created_at DATETIME,
		updated_at DATETIME
	) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_general_ci
	`
	if err := db.Exec(createWorkflowsSQL).Error; err != nil {
		return fmt.Errorf("failed to create workflows table: %w", err)
	}

	// Create workflow_steps table
	createStepsSQL := `
	CREATE TABLE IF NOT EXISTS workflow_steps (
		id VARCHAR(36) PRIMARY KEY,
		workflow_id VARCHAR(36) NOT NULL,
		level INT NOT NULL,
		actor_id VARCHAR(36) NOT NULL,
		conditions TEXT,
		description VARCHAR(500),
		created_at DATETIME,
		updated_at DATETIME,
		INDEX idx_workflow_id (workflow_id),
		INDEX idx_actor_id (actor_id),
		UNIQUE INDEX idx_workflow_level (workflow_id, level)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_general_ci
	`
	if err := db.Exec(createStepsSQL).Error; err != nil {
		return fmt.Errorf("failed to create workflow_steps table: %w", err)
	}

	// Create requests table
	createRequestsSQL := `
	CREATE TABLE IF NOT EXISTS requests (
		id VARCHAR(36) PRIMARY KEY,
		workflow_id VARCHAR(36) NOT NULL,
		requester_id VARCHAR(36) NOT NULL,
		current_step INT DEFAULT 1,
		status VARCHAR(20) DEFAULT 'PENDING',
		amount DECIMAL(15,2) NOT NULL,
		title VARCHAR(255),
		description TEXT,
		version INT DEFAULT 1,
		created_at DATETIME,
		updated_at DATETIME,
		INDEX idx_workflow_id (workflow_id),
		INDEX idx_requester_id (requester_id)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_general_ci
	`
	if err := db.Exec(createRequestsSQL).Error; err != nil {
		return fmt.Errorf("failed to create requests table: %w", err)
	}

	// Create approval_history table (with user_id column via ALTER for safety)
	createApprovalHistorySQL := `
	CREATE TABLE IF NOT EXISTS approval_history (
		id VARCHAR(36) PRIMARY KEY,
		request_id VARCHAR(36) NOT NULL,
		workflow_id VARCHAR(36) NOT NULL,
		step_level INT NOT NULL,
		actor_id VARCHAR(36) NOT NULL,
		user_id VARCHAR(36) NOT NULL,
		action VARCHAR(20) NOT NULL,
		comment TEXT,
		created_at DATETIME,
		INDEX idx_request_id (request_id),
		INDEX idx_actor_id (actor_id),
		INDEX idx_user_id (user_id)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_general_ci
	`
	if err := db.Exec(createApprovalHistorySQL).Error; err != nil {
		return fmt.Errorf("failed to create approval_history table: %w", err)
	}

	// Add user_id column if it doesn't exist (for existing tables)
	alterApprovalHistorySQL := `
	ALTER TABLE approval_history
	ADD COLUMN IF NOT EXISTS user_id VARCHAR(36) NOT NULL AFTER actor_id,
	ADD INDEX IF NOT EXISTS idx_user_id (user_id)
	`
	if err := db.Exec(alterApprovalHistorySQL).Error; err != nil {
		// Log warning but don't fail - column might already exist
		log.Printf("Warning: failed to add user_id column to approval_history: %v", err)
	}

	// Re-enable foreign key checks
	db.Exec("SET FOREIGN_KEY_CHECKS=1")

	log.Println("Database migrations completed successfully")
	return nil
}
