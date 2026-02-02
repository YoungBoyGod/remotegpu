package security

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"time"
)

// TLSManager TLS 证书管理器
type TLSManager struct{}

// NewTLSManager 创建 TLS 证书管理器
func NewTLSManager() *TLSManager {
	return &TLSManager{}
}

// TLSCertConfig TLS 证书配置
type TLSCertConfig struct {
	CommonName   string
	Organization string
	Country      string
	Province     string
	Locality     string
	ValidFor     time.Duration
	IsCA         bool
	KeySize      int
}

// GenerateSelfSignedCert 生成自签名证书
func (m *TLSManager) GenerateSelfSignedCert(config *TLSCertConfig) (*TLSCertificate, error) {
	// 生成私钥
	privateKey, err := rsa.GenerateKey(rand.Reader, config.KeySize)
	if err != nil {
		return nil, fmt.Errorf("生成私钥失败: %w", err)
	}

	// 创建证书模板
	template := m.createCertTemplate(config)

	// 自签名
	certBytes, err := x509.CreateCertificate(rand.Reader, template, template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return nil, fmt.Errorf("创建证书失败: %w", err)
	}

	// 编码证书
	certPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})

	// 编码私钥
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})

	return &TLSCertificate{
		CommonName:   config.CommonName,
		Certificate:  string(certPEM),
		PrivateKey:   string(privateKeyPEM),
		ValidFrom:    template.NotBefore,
		ValidUntil:   template.NotAfter,
		IsCA:         config.IsCA,
		Organization: config.Organization,
	}, nil
}

// createCertTemplate 创建证书模板
func (m *TLSManager) createCertTemplate(config *TLSCertConfig) *x509.Certificate {
	notBefore := time.Now()
	notAfter := notBefore.Add(config.ValidFor)

	serialNumber, _ := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))

	template := &x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			CommonName:   config.CommonName,
			Organization: []string{config.Organization},
			Country:      []string{config.Country},
			Province:     []string{config.Province},
			Locality:     []string{config.Locality},
		},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
	}

	if config.IsCA {
		template.IsCA = true
		template.KeyUsage |= x509.KeyUsageCertSign
	}

	return template
}

// ValidateCert 验证证书
func (m *TLSManager) ValidateCert(certPEM string) error {
	block, _ := pem.Decode([]byte(certPEM))
	if block == nil {
		return fmt.Errorf("无法解析证书 PEM")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return fmt.Errorf("解析证书失败: %w", err)
	}

	// 检查证书是否过期
	now := time.Now()
	if now.Before(cert.NotBefore) {
		return fmt.Errorf("证书尚未生效")
	}
	if now.After(cert.NotAfter) {
		return fmt.Errorf("证书已过期")
	}

	return nil
}
