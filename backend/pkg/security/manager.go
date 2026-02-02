package security

import (
	"fmt"
	"sync"
	"time"
)

// SecurityManager 安全管理器
type SecurityManager struct {
	passwordGen *PasswordGenerator
	sshKeyMgr   *SSHKeyManager
	tokenMgr    *TokenManager
	tlsMgr      *TLSManager

	// 存储生成的凭证
	credentials map[string]interface{}
	mu          sync.RWMutex
}

// NewSecurityManager 创建安全管理器
func NewSecurityManager(jwtSecret string) *SecurityManager {
	return &SecurityManager{
		passwordGen: NewPasswordGenerator(),
		sshKeyMgr:   NewSSHKeyManager(),
		tokenMgr:    NewTokenManager(jwtSecret),
		tlsMgr:      NewTLSManager(),
		credentials: make(map[string]interface{}),
	}
}

// GeneratePassword 生成密码
func (m *SecurityManager) GeneratePassword(strength PasswordStrength) (string, error) {
	return m.passwordGen.GenerateByStrength(strength)
}

// GenerateCustomPassword 生成自定义配置的密码
func (m *SecurityManager) GenerateCustomPassword(config *PasswordConfig) (string, error) {
	return m.passwordGen.Generate(config)
}

// ValidatePasswordStrength 验证密码强度
func (m *SecurityManager) ValidatePasswordStrength(password string, minStrength PasswordStrength) error {
	return m.passwordGen.ValidateStrength(password, minStrength)
}

// GenerateSSHKeyPair 生成 SSH 密钥对
func (m *SecurityManager) GenerateSSHKeyPair(keyType SSHKeyType, comment string) (*SSHKeyPair, error) {
	return m.sshKeyMgr.GenerateKeyPair(keyType, comment)
}

// ValidateSSHPublicKey 验证 SSH 公钥
func (m *SecurityManager) ValidateSSHPublicKey(publicKey string) error {
	return m.sshKeyMgr.ValidatePublicKey(publicKey)
}

// GenerateJWTToken 生成 JWT Token
func (m *SecurityManager) GenerateJWTToken(claims map[string]interface{}, expiresIn time.Duration) (*Token, error) {
	return m.tokenMgr.GenerateJWT(claims, expiresIn)
}

// ValidateJWTToken 验证 JWT Token
func (m *SecurityManager) ValidateJWTToken(tokenString string) (map[string]interface{}, error) {
	return m.tokenMgr.ValidateJWT(tokenString)
}

// GenerateAPIKey 生成 API Key
func (m *SecurityManager) GenerateAPIKey(prefix string, length int) (*Token, error) {
	return m.tokenMgr.GenerateAPIKey(prefix, length)
}

// GenerateSessionToken 生成 Session Token
func (m *SecurityManager) GenerateSessionToken(length int, expiresIn time.Duration) (*Token, error) {
	return m.tokenMgr.GenerateSessionToken(length, expiresIn)
}

// IsTokenExpired 检查 Token 是否过期
func (m *SecurityManager) IsTokenExpired(token *Token) bool {
	return m.tokenMgr.IsExpired(token)
}

// GenerateTLSCert 生成 TLS 证书
func (m *SecurityManager) GenerateTLSCert(config *TLSCertConfig) (*TLSCertificate, error) {
	return m.tlsMgr.GenerateSelfSignedCert(config)
}

// ValidateTLSCert 验证 TLS 证书
func (m *SecurityManager) ValidateTLSCert(certPEM string) error {
	return m.tlsMgr.ValidateCert(certPEM)
}

// StoreCredential 存储凭证
func (m *SecurityManager) StoreCredential(envID string, credential interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.credentials[envID] = credential
}

// GetCredential 获取凭证
func (m *SecurityManager) GetCredential(envID string) (interface{}, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	credential, ok := m.credentials[envID]
	if !ok {
		return nil, fmt.Errorf("凭证不存在: %s", envID)
	}
	return credential, nil
}

// DeleteCredential 删除凭证
func (m *SecurityManager) DeleteCredential(envID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.credentials, envID)
}

