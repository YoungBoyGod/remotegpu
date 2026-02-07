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

// KeySyncer 密钥同步接口（避免循环依赖）
type KeySyncer interface {
	EnqueueSyncKeys(ctx context.Context, hostID string, publicKeys []string, username string) error
}

type SSHKeyService struct {
	sshKeyDao     *dao.SSHKeyDao
	allocationDao *dao.AllocationDao
	keySyncer     KeySyncer
}

func NewSSHKeyService(db *gorm.DB) *SSHKeyService {
	return &SSHKeyService{
		sshKeyDao:     dao.NewSSHKeyDao(db),
		allocationDao: dao.NewAllocationDao(db),
	}
}

// SetKeySyncer 设置密钥同步器（在 router 初始化时注入，避免循环依赖）
func (s *SSHKeyService) SetKeySyncer(syncer KeySyncer) {
	s.keySyncer = syncer
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

	// 异步同步密钥到客户已分配的所有机器
	go s.syncKeysToAllocatedMachines(context.Background(), customerID)

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

	if err := s.sshKeyDao.Delete(ctx, keyID); err != nil {
		return err
	}

	// 异步同步密钥到客户已分配的所有机器
	go s.syncKeysToAllocatedMachines(context.Background(), customerID)

	return nil
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

// syncKeysToAllocatedMachines 将客户的所有公钥同步到其已分配的机器
func (s *SSHKeyService) syncKeysToAllocatedMachines(ctx context.Context, customerID uint) {
	if s.keySyncer == nil {
		return
	}

	// 获取客户所有公钥
	keys, err := s.sshKeyDao.ListByCustomerID(ctx, customerID)
	if err != nil {
		return
	}

	publicKeys := make([]string, 0, len(keys))
	for _, k := range keys {
		publicKeys = append(publicKeys, k.PublicKey)
	}

	// 获取客户所有活跃分配的机器
	allocations, err := s.allocationDao.FindAllActiveByCustomerID(ctx, customerID)
	if err != nil {
		return
	}

	// 逐台机器入队同步任务
	for _, alloc := range allocations {
		_ = s.keySyncer.EnqueueSyncKeys(ctx, alloc.HostID, publicKeys, "")
	}
}
