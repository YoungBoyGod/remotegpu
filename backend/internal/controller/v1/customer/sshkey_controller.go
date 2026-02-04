package customer

import (
	"errors"
	"strconv"

	apiV1 "github.com/YoungBoyGod/remotegpu/api/v1"
	"github.com/YoungBoyGod/remotegpu/internal/controller/v1/common"
	"github.com/YoungBoyGod/remotegpu/internal/service/sshkey"
	"github.com/gin-gonic/gin"
)

type SSHKeyController struct {
	common.BaseController
	sshKeyService *sshkey.SSHKeyService
}

func NewSSHKeyController(svc *sshkey.SSHKeyService) *SSHKeyController {
	return &SSHKeyController{
		sshKeyService: svc,
	}
}

func (c *SSHKeyController) getUserID(ctx *gin.Context) (uint, bool) {
	// CodeX 2026-02-04: avoid type-assert panic by using GetUint.
	userID := ctx.GetUint("userID")
	if userID == 0 {
		c.Error(ctx, 401, "未授权")
		return 0, false
	}
	return userID, true
}

// List 列出当前用户的所有 SSH 密钥
func (c *SSHKeyController) List(ctx *gin.Context) {
	userID, ok := c.getUserID(ctx)
	if !ok {
		return
	}

	keys, err := c.sshKeyService.ListKeys(ctx, userID)
	if err != nil {
		c.Error(ctx, 500, "获取 SSH 密钥列表失败")
		return
	}

	c.Success(ctx, gin.H{"list": keys})
}

// Create 创建新的 SSH 密钥
func (c *SSHKeyController) Create(ctx *gin.Context) {
	userID, ok := c.getUserID(ctx)
	if !ok {
		return
	}

	var req apiV1.CreateSSHKeyRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.Error(ctx, 400, "请求参数错误: "+err.Error())
		return
	}

	key, err := c.sshKeyService.CreateKey(ctx, userID, req.Name, req.PublicKey)
	if err != nil {
		if errors.Is(err, sshkey.ErrInvalidPublicKey) {
			c.Error(ctx, 400, "无效的 SSH 公钥格式")
			return
		}
		if errors.Is(err, sshkey.ErrKeyAlreadyExists) {
			c.Error(ctx, 409, "该 SSH 密钥已存在")
			return
		}
		c.Error(ctx, 500, "创建 SSH 密钥失败")
		return
	}

	c.Success(ctx, key)
}

// Delete 删除 SSH 密钥
func (c *SSHKeyController) Delete(ctx *gin.Context) {
	userID, ok := c.getUserID(ctx)
	if !ok {
		return
	}

	keyIDStr := ctx.Param("id")
	keyID, err := strconv.ParseUint(keyIDStr, 10, 64)
	if err != nil {
		c.Error(ctx, 400, "无效的密钥 ID")
		return
	}

	err = c.sshKeyService.DeleteKey(ctx, userID, uint(keyID))
	if err != nil {
		if errors.Is(err, sshkey.ErrKeyNotFound) {
			c.Error(ctx, 404, "SSH 密钥不存在")
			return
		}
		if errors.Is(err, sshkey.ErrKeyNotOwnedByUser) {
			c.Error(ctx, 403, "无权删除此密钥")
			return
		}
		c.Error(ctx, 500, "删除 SSH 密钥失败")
		return
	}

	c.Success(ctx, nil)
}
