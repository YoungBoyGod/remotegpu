package security

import "time"

// PasswordStrength 密码强度
type PasswordStrength string

const (
	// PasswordStrengthWeak 弱密码 (8位,仅字母数字)
	PasswordStrengthWeak PasswordStrength = "weak"

	// PasswordStrengthMedium 中等密码 (12位,字母数字+特殊字符)
	PasswordStrengthMedium PasswordStrength = "medium"

	// PasswordStrengthStrong 强密码 (16位,字母数字+特殊字符)
	PasswordStrengthStrong PasswordStrength = "strong"
)

// SSHKeyType SSH 密钥类型
type SSHKeyType string

const (
	// SSHKeyTypeRSA RSA 密钥
	SSHKeyTypeRSA SSHKeyType = "rsa"

	// SSHKeyTypeED25519 ED25519 密钥
	SSHKeyTypeED25519 SSHKeyType = "ed25519"

	// SSHKeyTypeECDSA ECDSA 密钥
	SSHKeyTypeECDSA SSHKeyType = "ecdsa"
)

// TokenType Token 类型
type TokenType string

const (
	// TokenTypeJWT JWT Token
	TokenTypeJWT TokenType = "jwt"

	// TokenTypeAPIKey API Key
	TokenTypeAPIKey TokenType = "api_key"

	// TokenTypeSession Session Token
	TokenTypeSession TokenType = "session"
)

// SSHKeyPair SSH 密钥对
type SSHKeyPair struct {
	Type       SSHKeyType `json:"type"`
	PublicKey  string     `json:"public_key"`
	PrivateKey string     `json:"private_key"`
	Comment    string     `json:"comment"`
	CreatedAt  time.Time  `json:"created_at"`
}

// TLSCertificate TLS 证书
type TLSCertificate struct {
	CommonName   string    `json:"common_name"`
	Certificate  string    `json:"certificate"`
	PrivateKey   string    `json:"private_key"`
	CACert       string    `json:"ca_cert"`
	ValidFrom    time.Time `json:"valid_from"`
	ValidUntil   time.Time `json:"valid_until"`
	IsCA         bool      `json:"is_ca"`
	Organization string    `json:"organization"`
}

// Token Token 信息
type Token struct {
	Type      TokenType         `json:"type"`
	Value     string            `json:"value"`
	ExpiresAt time.Time         `json:"expires_at"`
	Metadata  map[string]string `json:"metadata"`
	CreatedAt time.Time         `json:"created_at"`
}
