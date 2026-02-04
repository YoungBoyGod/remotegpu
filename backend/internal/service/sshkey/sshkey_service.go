package sshkey

import (
	"context"
	"errors"
	"strings"

	"github.com/YoungBoyGod/remotegpu/internal/dao"
	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"golang.org/x/crypto/ssh"
	"gorm.io/gorm"
)

var (
	ErrInvalidPublicKey  = errors.New("无效的 SSH 公钥格式")
	ErrKeyAlreadyExists  = errors.New("该 SSH 密钥已存在")
	ErrKeyNotFound       = errors.New("SSH 密钥不存在")
	ErrKeyNotOwnedByUser = errors.New("无权操作此 SSH 密钥")
)

type SSHKeyService struct {
	sshKeyDao *dao.SSHKeyDao
}

func NewSSHKeyService(db *gorm.DB) *SSHKeyService {
	return &SSHKeyService{
		sshKeyDao: dao.NewSSHKeyDao(db),
	}
}

// CreateKey 创建 SSH 密钥
func (s *SSHKeyService) CreateKey(ctx context.Context, customerID uint, name, publicKey string) (*entity.SSHKey, error) {
	// 解析并验证公钥
	fingerprint, err := parseSSHPublicKey(publicKey)
	if err != nil {
		return nil, ErrInvalidPublicKey
	}

	// 检查是否已存在相同指纹的密钥
	exists, err := s.sshKeyDao.ExistsByFingerprint(ctx, customerID, fingerprint)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrKeyAlreadyExists
	}

	key := &entity.SSHKey{
		CustomerID:  customerID,
		Name:        name,
		PublicKey:   strings.TrimSpace(publicKey),
		Fingerprint: fingerprint,
	}

	if err := s.sshKeyDao.Create(ctx, key); err != nil {
		return nil, err
	}

	return key, nil
}

// ListKeys 列出用户的所有 SSH 密钥
func (s *SSHKeyService) ListKeys(ctx context.Context, customerID uint) ([]entity.SSHKey, error) {
	return s.sshKeyDao.ListByCustomerID(ctx, customerID)
}

// DeleteKey 删除 SSH 密钥
func (s *SSHKeyService) DeleteKey(ctx context.Context, customerID, keyID uint) error {
	key, err := s.sshKeyDao.FindByID(ctx, keyID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrKeyNotFound
		}
		return err
	}

	if key.CustomerID != customerID {
		return ErrKeyNotOwnedByUser
	}

	return s.sshKeyDao.Delete(ctx, keyID)
}

// GetKey 获取单个 SSH 密钥
func (s *SSHKeyService) GetKey(ctx context.Context, customerID, keyID uint) (*entity.SSHKey, error) {
	key, err := s.sshKeyDao.FindByID(ctx, keyID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrKeyNotFound
		}
		return nil, err
	}

	if key.CustomerID != customerID {
		return nil, ErrKeyNotOwnedByUser
	}

	return key, nil
}

// parseSSHPublicKey 解析 SSH 公钥并返回指纹
func parseSSHPublicKey(publicKey string) (string, error) {
	// CodeX 2026-02-04: use ssh.ParseAuthorizedKey for robust parsing.
	key, _, _, _, err := ssh.ParseAuthorizedKey([]byte(strings.TrimSpace(publicKey)))
	if err != nil {
		return "", ErrInvalidPublicKey
	}

	return ssh.FingerprintSHA256(key), nil
}

// GetMD5Fingerprint 获取 MD5 格式的指纹（兼容旧格式）
func GetMD5Fingerprint(publicKey string) (string, error) {
	key, _, _, _, err := ssh.ParseAuthorizedKey([]byte(strings.TrimSpace(publicKey)))
	if err != nil {
		return "", ErrInvalidPublicKey
	}

	fingerprint := ssh.FingerprintLegacyMD5(key)
	return strings.TrimPrefix(fingerprint, "MD5:"), nil
}
