package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"os"

	"github.com/YoungBoyGod/remotegpu/config"
)

// EncryptAES256GCM 使用 AES-256-GCM 加密数据
// @author Claude
// @description 修复 P0 安全问题：SSH 凭据加密存储
// @param plaintext 明文
// @return 加密后的 base64 字符串、错误
// @modified 2026-02-06
func EncryptAES256GCM(plaintext string) (string, error) {
	if plaintext == "" {
		return "", nil
	}

	key, err := getEncryptionKey()
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecryptAES256GCM 使用 AES-256-GCM 解密数据
// @author Claude
// @description 修复 P0 安全问题：SSH 凭据解密
// @param ciphertext 加密后的 base64 字符串
// @return 明文、错误
// @modified 2026-02-06
func DecryptAES256GCM(ciphertext string) (string, error) {
	if ciphertext == "" {
		return "", nil
	}

	key, err := getEncryptionKey()
	if err != nil {
		return "", err
	}

	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	nonce, encryptedData := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, encryptedData, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

// getEncryptionKey 获取加密密钥
// 优先级：配置文件 > 环境变量 > 默认密钥
func getEncryptionKey() ([]byte, error) {
	var keyStr string

	// 1. 优先从配置文件读取
	if config.GlobalConfig != nil && config.GlobalConfig.Encryption.Key != "" {
		keyStr = config.GlobalConfig.Encryption.Key
	} else {
		// 2. 从环境变量读取
		keyStr = os.Getenv("ENCRYPTION_KEY")
		if keyStr == "" {
			// 3. 使用默认密钥（仅用于开发测试）
			keyStr = "default-32-byte-encryption-key!!"
		}
	}

	key := []byte(keyStr)
	if len(key) != 32 {
		return nil, errors.New("encryption key must be 32 bytes for AES-256")
	}

	return key, nil
}
