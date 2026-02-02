package security

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"time"

	"golang.org/x/crypto/ssh"
)

// SSHKeyManager SSH 密钥管理器
type SSHKeyManager struct{}

// NewSSHKeyManager 创建 SSH 密钥管理器
func NewSSHKeyManager() *SSHKeyManager {
	return &SSHKeyManager{}
}

// GenerateKeyPair 生成 SSH 密钥对
func (m *SSHKeyManager) GenerateKeyPair(keyType SSHKeyType, comment string) (*SSHKeyPair, error) {
	switch keyType {
	case SSHKeyTypeRSA:
		return m.generateRSAKeyPair(comment)
	case SSHKeyTypeED25519:
		return m.generateED25519KeyPair(comment)
	case SSHKeyTypeECDSA:
		return m.generateECDSAKeyPair(comment)
	default:
		return nil, fmt.Errorf("不支持的密钥类型: %s", keyType)
	}
}

// generateRSAKeyPair 生成 RSA 密钥对
func (m *SSHKeyManager) generateRSAKeyPair(comment string) (*SSHKeyPair, error) {
	// 生成 RSA 私钥 (4096 位)
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, fmt.Errorf("生成 RSA 私钥失败: %w", err)
	}

	// 编码私钥为 PEM 格式
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})

	// 生成公钥
	publicKey, err := ssh.NewPublicKey(&privateKey.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("生成公钥失败: %w", err)
	}

	publicKeyStr := string(ssh.MarshalAuthorizedKey(publicKey))
	if comment != "" {
		publicKeyStr = fmt.Sprintf("%s %s", publicKeyStr[:len(publicKeyStr)-1], comment)
	}

	return &SSHKeyPair{
		Type:       SSHKeyTypeRSA,
		PublicKey:  publicKeyStr,
		PrivateKey: string(privateKeyPEM),
		Comment:    comment,
		CreatedAt:  time.Now(),
	}, nil
}

// generateED25519KeyPair 生成 ED25519 密钥对
func (m *SSHKeyManager) generateED25519KeyPair(comment string) (*SSHKeyPair, error) {
	// 生成 ED25519 密钥对
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("生成 ED25519 密钥失败: %w", err)
	}

	// 编码私钥
	privateKeyBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		return nil, fmt.Errorf("编码私钥失败: %w", err)
	}

	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: privateKeyBytes,
	})

	// 生成公钥
	sshPublicKey, err := ssh.NewPublicKey(publicKey)
	if err != nil {
		return nil, fmt.Errorf("生成公钥失败: %w", err)
	}

	publicKeyStr := string(ssh.MarshalAuthorizedKey(sshPublicKey))
	if comment != "" {
		publicKeyStr = fmt.Sprintf("%s %s", publicKeyStr[:len(publicKeyStr)-1], comment)
	}

	return &SSHKeyPair{
		Type:       SSHKeyTypeED25519,
		PublicKey:  publicKeyStr,
		PrivateKey: string(privateKeyPEM),
		Comment:    comment,
		CreatedAt:  time.Now(),
	}, nil
}

// generateECDSAKeyPair 生成 ECDSA 密钥对
func (m *SSHKeyManager) generateECDSAKeyPair(comment string) (*SSHKeyPair, error) {
	// 生成 ECDSA 私钥 (P-256 曲线)
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("生成 ECDSA 私钥失败: %w", err)
	}

	// 编码私钥
	privateKeyBytes, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		return nil, fmt.Errorf("编码私钥失败: %w", err)
	}

	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: privateKeyBytes,
	})

	// 生成公钥
	publicKey, err := ssh.NewPublicKey(&privateKey.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("生成公钥失败: %w", err)
	}

	publicKeyStr := string(ssh.MarshalAuthorizedKey(publicKey))
	if comment != "" {
		publicKeyStr = fmt.Sprintf("%s %s", publicKeyStr[:len(publicKeyStr)-1], comment)
	}

	return &SSHKeyPair{
		Type:       SSHKeyTypeECDSA,
		PublicKey:  publicKeyStr,
		PrivateKey: string(privateKeyPEM),
		Comment:    comment,
		CreatedAt:  time.Now(),
	}, nil
}

// ValidatePublicKey 验证公钥格式
func (m *SSHKeyManager) ValidatePublicKey(publicKey string) error {
	_, _, _, _, err := ssh.ParseAuthorizedKey([]byte(publicKey))
	if err != nil {
		return fmt.Errorf("无效的公钥格式: %w", err)
	}
	return nil
}
