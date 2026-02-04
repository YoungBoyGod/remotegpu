package auth

import (
	"strings"

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
	var req apiV1.RefreshTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.Error(ctx, 400, "Invalid request parameters")
		return
	}

	accessToken, refreshToken, expiresIn, err := c.authService.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		c.Error(ctx, 401, "Invalid refresh token")
		return
	}

	c.Success(ctx, apiV1.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    expiresIn,
	})
}

func (c *AuthController) Logout(ctx *gin.Context) {
	// 从 Header 获取 token
	authHeader := ctx.GetHeader("Authorization")
	if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
		token := strings.TrimPrefix(authHeader, "Bearer ")
		c.authService.Logout(ctx, token)
	}
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

// AdminLogin Admin 专用登录接口
func (c *AuthController) AdminLogin(ctx *gin.Context) {
	var req apiV1.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.Error(ctx, 400, "Invalid request parameters")
		return
	}

	accessToken, refreshToken, expiresIn, err := c.authService.AdminLogin(ctx, req.Username, req.Password)
	if err != nil {
		if err.Error() == "permission denied: admin role required" {
			c.Error(ctx, 403, "Permission denied")
			return
		}
		c.Error(ctx, 401, "Authentication failed")
		return
	}

	c.Success(ctx, apiV1.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    expiresIn,
	})
}
