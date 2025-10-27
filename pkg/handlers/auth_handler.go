package handlers

import (
	"mini-trh-backend/internal/consts"
	"mini-trh-backend/internal/utils"
	"mini-trh-backend/pkg/api/dto"
	"mini-trh-backend/pkg/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthHandler struct {
	authService *services.AuthService
	userRepo    interface {
		FindById(uuid.UUID) (interface{}, error)
		FindAll(int, int) (interface{}, error)
	}
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Login handles user login
// @Summary User login
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login credentials"
// @Success 200 {object} dto.LoginResponse
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, err.Error())
		return
	}

	token, user, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		if err == services.ErrInvalidCredentials {
			utils.UnauthorizedResponse(c, "Invalid email or password")
			return
		}
		utils.InternalServerErrorResponse(c, "Failed to login")
		return
	}

	response := dto.LoginResponse{
		Token: token,
		User: dto.UserDTO{
			ID:        user.ID,
			Email:     user.Email,
			Role:      user.Role,
			CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z"),
		},
	}
	utils.SuccessResponse(c, response)
}

// Register handles user registration (Admin only)
// @Summary Register new user
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.RegisterRequest true "User registration data"
// @Success 201 {object} dto.UserDTO
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, err.Error())
		return
	}

	// Default role is user
	role := req.Role
	if role == "" {
		role = consts.RoleUser
	}

	// Create user
	user, err := h.authService.Register(req.Email, req.Password, role)
	if err != nil {
		if err == services.ErrUserAlreadyExists {
			utils.BadRequestResponse(c, "User with this email already exists")
			return
		}
		utils.InternalServerErrorResponse(c, "Failed to register user")
		return
	}

	response := dto.UserDTO{
		ID:        user.ID,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}

	utils.CreatedResponse(c, response)
}

// GetProfile returns the current user's profile
// @Summary Get user profile
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.UserDTO
// @Router /auth/profile [get]
func (h *AuthHandler) GetProfile(c *gin.Context) {
	// Get user info from context (set by JWT middleware)
	userID, _ := c.Get("user_id")
	email, _ := c.Get("user_email")
	role, _ := c.Get("user_role")

	response := dto.UserDTO{
		ID:    userID.(uuid.UUID),
		Email: email.(string),
		Role:  role.(string),
	}

	utils.SuccessResponse(c, response)
}
