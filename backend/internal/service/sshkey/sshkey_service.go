package sshkey

import (
	"context"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"github.com/YoungBoyGod/remotegpu/internal/dao"
	"github.com/YoungBoyGod/remotegpu/internal/model/entity"
	"gorm.io/gorm"
)

var (
	ErrInvalidPublicKey   = errors.New("无效的 SSH 公钥格式")
	ErrKeyAlreadyExists   = errors.New("该 SSH 密钥已存在")
	ErrKeyNotFound        = errors.New("SSH 密钥不存在")
	ErrKeyNotOwnedByUser  = errors.New("无权操作此 SSH 密钥")
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
	parts := strings.Fields(publicKey)
	if len(parts) < 2 {
		return "", ErrInvalidPublicKey
	}

	keyType := parts[0]
	if keyType != "ssh-rsa" && keyType != "ssh-ed25519" &&
	   keyType != "ecdsa-sha2-nistp256" && keyType != "ecdsa-sha2-nistp384" &&
	   keyType != "ecdsa-sha2-nistp521" {
		return "", ErrInvalidPublicKey
	}

	keyData, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return "", ErrInvalidPublicKey
	}

	// 生成 SHA256 指纹
	hash := sha256.Sum256(keyData)
	fingerprint := "SHA256:" + base64.RawStdEncoding.EncodeToString(hash[:])

	return fingerprint, nil
}

// GetMD5Fingerprint 获取 MD5 格式的指纹（兼容旧格式）
func GetMD5Fingerprint(publicKey string) (string, error) {
	parts := strings.Fields(publicKey)
	if len(parts) < 2 {
		return "", ErrInvalidPublicKey
	}

	keyData, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return "", ErrInvalidPublicKey
	}

	hash := md5.Sum(keyData)
	var fp []string
	for _, b := range hash {
		fp = append(fp, fmt.Sprintf("%02x", b))
	}
	return strings.Join(fp, ":"), nil
}
