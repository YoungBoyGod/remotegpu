package auth

import (
	apiV1 "github.com/YoungBoyGod/remotegpu/api/v1"
	"github.com/YoungBoyGod/remotegpu/internal/controller/v1/common"
	serviceAuth "github.com/YoungBoyGod/remotegpu/internal/service/auth"
	"github.com/gin-gonic/gin"
)

type AuthController struct {
	common.BaseController
	authService *serviceAuth.AuthService
}

func NewAuthController(authService *serviceAuth.AuthService) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

func (c *AuthController) Login(ctx *gin.Context) {
	var req apiV1.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.Error(ctx, 400, "Invalid request parameters")
		return
	}

	accessToken, refreshToken, expiresIn, err := c.authService.Login(ctx, req.Username, req.Password)
	if err != nil {
		c.Error(ctx, 401, "Authentication failed")
		return
	}

	c.Success(ctx, apiV1.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    expiresIn,
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
