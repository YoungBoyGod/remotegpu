package auth

import (
	"strings"

	apiV1 "github.com/YoungBoyGod/remotegpu/api/v1"
	"github.com/YoungBoyGod/remotegpu/internal/controller/v1/common"
	serviceAuth "github.com/YoungBoyGod/remotegpu/internal/service/auth"
	"github.com/YoungBoyGod/remotegpu/pkg/errors"
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

// Login 用户登录
// @Summary 用户登录
// @Description 使用用户名和密码获取访问令牌
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body v1.LoginRequest true "登录请求"
// @Success 200 {object} v1.LoginResponse
// @Failure 400 {object} common.ErrorResponse
// @Failure 401 {object} common.ErrorResponse
// @Router /auth/login [post]
func (c *AuthController) Login(ctx *gin.Context) {
	var req apiV1.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.Error(ctx, 400, "Invalid request parameters")
		return
	}

	accessToken, refreshToken, expiresIn, err := c.authService.Login(ctx, req.Username, req.Password)
	if err != nil {
		if appErr := errors.GetAppError(err); appErr != nil {
			c.Error(ctx, appErr.Code, appErr.Message)
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

// Refresh 刷新令牌
// @Summary 刷新令牌
// @Description 使用刷新令牌获取新的访问令牌
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body v1.RefreshTokenRequest true "刷新令牌请求"
// @Success 200 {object} v1.LoginResponse
// @Failure 400 {object} common.ErrorResponse
// @Failure 401 {object} common.ErrorResponse
// @Router /auth/refresh [post]
func (c *AuthController) Refresh(ctx *gin.Context) {
	var req apiV1.RefreshTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.Error(ctx, 400, "Invalid request parameters")
		return
	}

	accessToken, refreshToken, expiresIn, err := c.authService.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		if appErr := errors.GetAppError(err); appErr != nil {
			c.Error(ctx, appErr.Code, appErr.Message)
			return
		}
		c.Error(ctx, 401, "Invalid refresh token")
		return
	}

	c.Success(ctx, apiV1.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    expiresIn,
	})
}

// Logout 用户登出
// @Summary 用户登出
// @Description 使当前访问令牌失效
// @Tags Auth
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} map[string]string
// @Router /auth/logout [post]
func (c *AuthController) Logout(ctx *gin.Context) {
	// 从 Header 获取 token
	authHeader := ctx.GetHeader("Authorization")
	if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
		token := strings.TrimPrefix(authHeader, "Bearer ")
		c.authService.Logout(ctx, token)
	}
	c.Success(ctx, gin.H{"message": "Logged out successfully"})
}

// GetProfile 获取个人资料
// @Summary 获取个人资料
// @Description 获取当前登录用户的详细信息
// @Tags Auth
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} entity.Customer
// @Failure 401 {object} common.ErrorResponse
// @Router /auth/profile [get]
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

// AdminLogin 管理员登录
// @Summary 管理员登录
// @Description 管理员专用登录接口，校验管理员角色
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body v1.LoginRequest true "登录请求"
// @Success 200 {object} v1.LoginResponse
// @Failure 400 {object} common.ErrorResponse
// @Failure 401 {object} common.ErrorResponse
// @Failure 403 {object} common.ErrorResponse
// @Router /auth/admin/login [post]
func (c *AuthController) AdminLogin(ctx *gin.Context) {
	var req apiV1.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.Error(ctx, 400, "Invalid request parameters")
		return
	}

	accessToken, refreshToken, expiresIn, err := c.authService.AdminLogin(ctx, req.Username, req.Password)
	if err != nil {
		if appErr := errors.GetAppError(err); appErr != nil {
			c.Error(ctx, appErr.Code, appErr.Message)
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
