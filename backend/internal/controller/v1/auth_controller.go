package v1

import (
	"github.com/YoungBoyGod/remotegpu/internal/service"
	"github.com/gin-gonic/gin"
)

type AuthController struct {
	BaseController
	authService *service.AuthService
}

func NewAuthController(authService *service.AuthService) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (c *AuthController) Login(ctx *gin.Context) {
	var req LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.Error(ctx, 400, "Invalid request parameters")
		return
	}

	accessToken, refreshToken, expiresIn, err := c.authService.Login(ctx, req.Username, req.Password)
	if err != nil {
		c.Error(ctx, 401, "Authentication failed")
		return
	}

	c.Success(ctx, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"expires_in":    expiresIn,
	})
}

func (c *AuthController) Refresh(ctx *gin.Context) {
	// TODO: Implement refresh token logic
	c.Success(ctx, gin.H{"message": "Token refreshed"})
}

func (c *AuthController) Logout(ctx *gin.Context) {
	c.Success(ctx, gin.H{"message": "Logged out successfully"})
}

func (c *AuthController) GetProfile(ctx *gin.Context) {
	// User ID comes from middleware
	userID := ctx.GetUint("userID")
	profile, err := c.authService.GetProfile(ctx, userID)
	if err != nil {
		c.Error(ctx, 500, "Failed to get profile")
		return
	}
	c.Success(ctx, profile)
}
