package usecase

import (
	"context"
	"errors"
	"regexp"

	"golang.org/x/crypto/bcrypt"

	"workflow-approval/package/user/domain"
	"workflow-approval/package/user/ports"
	"workflow-approval/package/user/repository"
)

var (
	ErrInvalidEmail       = errors.New("invalid email format")
	ErrWeakPassword       = errors.New("password must be at least 8 characters")
	ErrInvalidLogin       = errors.New("invalid email or password")
	ErrEmailRequired      = errors.New("email is required")
	ErrPasswordRequired   = errors.New("password is required")
	ErrNameRequired       = errors.New("name is required")
	ErrActorIDRequired    = errors.New("actor_id is required for non-admin users")
	ErrActorIDNotRequired = errors.New("actor_id should be empty for admin users")
	ErrActorIDNotFound    = errors.New("actor_id not found")
)

// UserServiceImpl implements UserService interface
type UserServiceImpl struct {
	userRepo  ports.UserRepository
	actorRepo ports.ActorRepository
}

// NewUserService creates a new UserServiceImpl instance
func NewUserService(userRepo ports.UserRepository, actorRepo ports.ActorRepository) ports.UserService {
	return &UserServiceImpl{
		userRepo:  userRepo,
		actorRepo: actorRepo,
	}
}

// Register creates a new user account
func (s *UserServiceImpl) Register(ctx context.Context, email, password, name string, isAdmin bool, actorID *string) (*domain.User, error) {
	// Validation
	if email == "" {
		return nil, ErrEmailRequired
	}
	if password == "" {
		return nil, ErrPasswordRequired
	}
	if name == "" {
		return nil, ErrNameRequired
	}

	// Validate actor_id based on is_admin
	if isAdmin {
		// Admin users should not have actor_id
		if actorID != nil && *actorID != "" {
			return nil, ErrActorIDNotRequired
		}
		actorID = nil
	} else {
		// Non-admin users must have actor_id
		if actorID == nil || *actorID == "" {
			return nil, ErrActorIDRequired
		}

		// Verify actor_id exists in database
		_, err := s.actorRepo.GetByID(ctx, *actorID)
		if err != nil {
			return nil, ErrActorIDNotFound
		}
	}

	// Validate email format
	if !isValidEmail(email) {
		return nil, ErrInvalidEmail
	}

	// Validate password strength
	if len(password) < 8 {
		return nil, ErrWeakPassword
	}

	// Check if user already exists
	_, err := s.userRepo.GetByEmail(ctx, email)
	if err == nil {
		return nil, ErrUserAlreadyExist
	}
	if !errors.Is(err, repository.ErrUserNotFound) {
		return nil, err
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create user
	user := domain.NewUser(email, string(hashedPassword), name, isAdmin, actorID)
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// Login authenticates a user and returns a JWT token
func (s *UserServiceImpl) Login(ctx context.Context, email, password string) (*domain.User, string, error) {
	// Validation
	if email == "" || password == "" {
		return nil, "", ErrInvalidLogin
	}

	// Get user by email
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, "", ErrInvalidLogin
		}
		return nil, "", err
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, "", ErrInvalidLogin
	}

	// Generate JWT token (token generation will be handled by the handler)
	return user, "", nil
}

// GetUserByID retrieves a user by ID
func (s *UserServiceImpl) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	return s.userRepo.GetByID(ctx, id)
}

// GetUserByEmail retrieves a user by email
func (s *UserServiceImpl) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	return s.userRepo.GetByEmail(ctx, email)
}

// UpdateUser updates user's profile (name)
func (s *UserServiceImpl) UpdateUser(ctx context.Context, id string, name string) (*domain.User, error) {
	if name == "" {
		return nil, ErrNameRequired
	}

	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	user.Name = name
	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// isValidEmail validates email format using regex
func isValidEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(pattern, email)
	return matched
}

// ErrUserAlreadyExist is returned when a user already exists
var ErrUserAlreadyExist = errors.New("user already exists")
