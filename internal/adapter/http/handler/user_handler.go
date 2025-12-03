package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mashurimansur/goCMS/internal/domain/user"
	userusecase "github.com/mashurimansur/goCMS/internal/usecase/user"
)

type UserHandler struct {
	userUseCase userusecase.UseCase
}

func NewUserHandler(userUseCase userusecase.UseCase) *UserHandler {
	return &UserHandler{
		userUseCase: userUseCase,
	}
}

func (h *UserHandler) Register(router *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	public := router.Group("/auth")
	{
		public.POST("/register", h.register)
		public.POST("/login", h.login)
	}

	admin := router.Group("/admin/users")
	admin.Use(authMiddleware)
	{
		admin.GET("/:id", h.getProfile)
		admin.PUT("/:id", h.updateProfile)
		admin.GET("/", h.listUsers)
		admin.DELETE("/:id", h.deleteUser)
	}
}

type registerRequest struct {
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Username string `json:"username"`
	Phone    string `json:"phone"`
}

// @Summary      Register a new user
// @Description  Register a new user account
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body registerRequest true "Register Request"
// @Success      201  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /auth/register [post]
func (h *UserHandler) register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u := &user.User{
		FullName: req.FullName,
		Email:    req.Email,
		Username: req.Username,
		Phone:    req.Phone,
	}

	if err := h.userUseCase.Register(c.Request.Context(), u, req.Password); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "user created successfully"})
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type loginResponse struct {
	AccessToken string     `json:"access_token"`
	User        *user.User `json:"user"`
}

// @Summary      Login
// @Description  Authenticate user and get PASETO token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body loginRequest true "Login Request"
// @Success      200  {object}  loginResponse
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Router       /auth/login [post]
func (h *UserHandler) login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	accessToken, u, err := h.userUseCase.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, loginResponse{
		AccessToken: accessToken,
		User:        u,
	})
}

// @Summary      Get user profile
// @Description  Get authenticated user's profile
// @Tags         users
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  user.User
// @Failure      401  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /users/profile [get]
func (h *UserHandler) getProfile(c *gin.Context) {
	// Assuming auth middleware sets "user_id" in context
	// userID, exists := c.Get("user_id")
	// if !exists {
	// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
	// 	return
	// }
	id := c.Param("id")
	u, err := h.userUseCase.GetProfile(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, u)
}

// @Summary      Update user profile
// @Description  Update authenticated user's profile
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body user.User true "User Update Request"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /users/profile [put]
func (h *UserHandler) updateProfile(c *gin.Context) {
	var u user.User
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// userID, exists := c.Get("user_id")
	// if !exists {
	// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
	// 	return
	// }

	u.ID = c.Param("id")
	if err := h.userUseCase.UpdateProfile(c.Request.Context(), &u); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "profile updated successfully"})
}

// @Summary      List all users
// @Description  Get list of users with pagination
// @Tags         users
// @Produce      json
// @Security     BearerAuth
// @Param        limit   query     int  false  "Limit"  default(10)
// @Param        offset  query     int  false  "Offset" default(0)
// @Success      200  {array}   user.User
// @Failure      500  {object}  map[string]string
// @Router       /users [get]
func (h *UserHandler) listUsers(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, _ := strconv.Atoi(limitStr)
	offset, _ := strconv.Atoi(offsetStr)

	users, err := h.userUseCase.ListUsers(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}

// @Summary      Delete user
// @Description  Delete a user by ID
// @Tags         users
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "User ID"
// @Success      200  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /users/{id} [delete]
func (h *UserHandler) deleteUser(c *gin.Context) {
	id := c.Param("id")
	if err := h.userUseCase.DeleteUser(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user deleted successfully"})
}
