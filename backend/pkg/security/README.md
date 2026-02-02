# pkg/security 使用示例

## 1. 密码生成

```go
package main

import (
    "github.com/YoungBoyGod/remotegpu/pkg/security"
)

func ExamplePasswordGenerator() {
    // 创建密码生成器
    passwordGen := security.NewPasswordGenerator()

    // 根据强度生成密码
    weakPassword, _ := passwordGen.GenerateByStrength(security.PasswordStrengthWeak)
    // weakPassword: 8位,包含大小写字母和数字

    mediumPassword, _ := passwordGen.GenerateByStrength(security.PasswordStrengthMedium)
    // mediumPassword: 12位,包含大小写字母、数字和特殊字符

    strongPassword, _ := passwordGen.GenerateByStrength(security.PasswordStrengthStrong)
    // strongPassword: 16位,包含大小写字母、数字和特殊字符

    // 自定义配置生成密码
    config := &security.PasswordConfig{
        Length:         20,
        IncludeLower:   true,
        IncludeUpper:   true,
        IncludeDigits:  true,
        IncludeSpecial: true,
    }
    customPassword, _ := passwordGen.Generate(config)

    // 验证密码强度
    err := passwordGen.ValidateStrength(customPassword, security.PasswordStrengthMedium)
    if err != nil {
        // 密码强度不足
    }
}
```

## 2. SSH 密钥管理

```go
func ExampleSSHKeyManager() {
    // 创建 SSH 密钥管理器
    sshKeyMgr := security.NewSSHKeyManager()

    // 生成 RSA 密钥对
    rsaKeyPair, _ := sshKeyMgr.GenerateKeyPair(
        security.SSHKeyTypeRSA,
        "user@example.com",
    )
    // rsaKeyPair.PublicKey: SSH 公钥
    // rsaKeyPair.PrivateKey: SSH 私钥 (PEM 格式)

    // 生成 ED25519 密钥对 (推荐)
    ed25519KeyPair, _ := sshKeyMgr.GenerateKeyPair(
        security.SSHKeyTypeED25519,
        "user@example.com",
    )

    // 生成 ECDSA 密钥对
    ecdsaKeyPair, _ := sshKeyMgr.GenerateKeyPair(
        security.SSHKeyTypeECDSA,
        "user@example.com",
    )

    // 验证公钥格式
    err := sshKeyMgr.ValidatePublicKey(rsaKeyPair.PublicKey)
    if err != nil {
        // 公钥格式无效
    }
}
```

## 3. Token 管理

```go
func ExampleTokenManager() {
    // 创建 Token 管理器
    tokenMgr := security.NewTokenManager("your-jwt-secret-key")

    // 生成 JWT Token
    claims := map[string]interface{
        "user_id": "123",
        "role":    "admin",
    }
    jwtToken, _ := tokenMgr.GenerateJWT(claims, 24*time.Hour)
    // jwtToken.Value: JWT Token 字符串
    // jwtToken.ExpiresAt: 过期时间

    // 验证 JWT Token
    validatedClaims, err := tokenMgr.ValidateJWT(jwtToken.Value)
    if err != nil {
        // Token 无效或已过期
    }

    // 生成 API Key
    apiKey, _ := tokenMgr.GenerateAPIKey("env", 32)
    // apiKey.Value: "env_<random-base64-string>"

    // 生成 Session Token
    sessionToken, _ := tokenMgr.GenerateSessionToken(64, 7*24*time.Hour)
    // sessionToken.Value: 随机 Session Token
    // sessionToken.ExpiresAt: 7天后过期

    // 检查 Token 是否过期
    if tokenMgr.IsExpired(jwtToken) {
        // Token 已过期
    }
}
```

## 4. TLS 证书管理

```go
func ExampleTLSManager() {
    // 创建 TLS 证书管理器
    tlsMgr := security.NewTLSManager()

    // 生成自签名证书
    config := &security.TLSCertConfig{
        CommonName:   "example.com",
        Organization: "Example Corp",
        Country:      "CN",
        Province:     "Beijing",
        Locality:     "Beijing",
        ValidFor:     365 * 24 * time.Hour, // 1年
        IsCA:         false,
        KeySize:      2048,
    }

    cert, _ := tlsMgr.GenerateSelfSignedCert(config)
    // cert.Certificate: PEM 格式的证书
    // cert.PrivateKey: PEM 格式的私钥
    // cert.ValidFrom: 证书生效时间
    // cert.ValidUntil: 证书过期时间

    // 验证证书
    err := tlsMgr.ValidateCert(cert.Certificate)
    if err != nil {
        // 证书无效或已过期
    }
}
```

## 5. 安全管理器 (推荐使用)

```go
func ExampleSecurityManager() {
    // 创建安全管理器
    securityMgr := security.NewSecurityManager("your-jwt-secret-key")

    // 生成密码
    password, _ := securityMgr.GeneratePassword(security.PasswordStrengthStrong)

    // 生成 SSH 密钥对
    sshKey, _ := securityMgr.GenerateSSHKeyPair(
        security.SSHKeyTypeED25519,
        "env-123@example.com",
    )

    // 生成 API Key
    apiKey, _ := securityMgr.GenerateAPIKey("env-123", 32)

    // 生成 JWT Token
    claims := map[string]interface{}{
        "env_id": "env-123",
        "user_id": "user-456",
    }
    jwtToken, _ := securityMgr.GenerateJWTToken(claims, 24*time.Hour)

    // 生成 TLS 证书
    tlsConfig := &security.TLSCertConfig{
        CommonName:   "env-123.example.com",
        Organization: "Example Corp",
        ValidFor:     365 * 24 * time.Hour,
        KeySize:      2048,
    }
    tlsCert, _ := securityMgr.GenerateTLSCert(tlsConfig)

    // 存储凭证
    credentials := map[string]interface{}{
        "password": password,
        "ssh_key":  sshKey,
        "api_key":  apiKey,
        "jwt":      jwtToken,
        "tls":      tlsCert,
    }
    securityMgr.StoreCredential("env-123", credentials)

    // 获取凭证
    storedCreds, _ := securityMgr.GetCredential("env-123")

    // 删除凭证
    securityMgr.DeleteCredential("env-123")
}
```

## 6. 在 Service 层使用

```go
// internal/service/environment.go

import (
    "github.com/YoungBoyGod/remotegpu/pkg/security"
)

type EnvironmentService struct {
    securityMgr *security.SecurityManager
    // ...
}

func (s *EnvironmentService) CreateEnvironment(req *CreateEnvironmentRequest) (*entity.Environment, error) {
    // ... 创建环境逻辑 ...

    // 生成 SSH 密码
    sshPassword, err := s.securityMgr.GeneratePassword(security.PasswordStrengthStrong)
    if err != nil {
        return nil, fmt.Errorf("生成 SSH 密码失败: %w", err)
    }

    // 生成 SSH 密钥对
    sshKey, err := s.securityMgr.GenerateSSHKeyPair(
        security.SSHKeyTypeED25519,
        fmt.Sprintf("%s@%s", env.ID, "remotegpu.com"),
    )
    if err != nil {
        return nil, fmt.Errorf("生成 SSH 密钥失败: %w", err)
    }

    // 生成 API Key
    apiKey, err := s.securityMgr.GenerateAPIKey(env.ID, 32)
    if err != nil {
        return nil, fmt.Errorf("生成 API Key 失败: %w", err)
    }

    // 存储凭证
    credentials := map[string]interface{}{
        "ssh_password": sshPassword,
        "ssh_key":      sshKey,
        "api_key":      apiKey,
    }
    s.securityMgr.StoreCredential(env.ID, credentials)

    // 保存到数据库
    // ...

    return env, nil
}

func (s *EnvironmentService) DeleteEnvironment(envID string) error {
    // ... 删除环境逻辑 ...

    // 清理凭证
    s.securityMgr.DeleteCredential(envID)

    return nil
}
```

