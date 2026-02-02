package security

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"
)

const (
	// 字符集定义
	lowercaseLetters = "abcdefghijklmnopqrstuvwxyz"
	uppercaseLetters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digits           = "0123456789"
	specialChars     = "!@#$%^&*()_+-=[]{}|;:,.<>?"
)

// PasswordConfig 密码配置
type PasswordConfig struct {
	Length         int
	IncludeLower   bool
	IncludeUpper   bool
	IncludeDigits  bool
	IncludeSpecial bool
}

// PasswordGenerator 密码生成器
type PasswordGenerator struct{}

// NewPasswordGenerator 创建密码生成器
func NewPasswordGenerator() *PasswordGenerator {
	return &PasswordGenerator{}
}

// Generate 生成密码
func (g *PasswordGenerator) Generate(config *PasswordConfig) (string, error) {
	if config.Length < 8 {
		return "", fmt.Errorf("密码长度至少为 8 位")
	}

	// 构建字符集
	charset := ""
	if config.IncludeLower {
		charset += lowercaseLetters
	}
	if config.IncludeUpper {
		charset += uppercaseLetters
	}
	if config.IncludeDigits {
		charset += digits
	}
	if config.IncludeSpecial {
		charset += specialChars
	}

	if charset == "" {
		return "", fmt.Errorf("至少需要包含一种字符类型")
	}

	// 生成密码
	password := make([]byte, config.Length)
	for i := 0; i < config.Length; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", fmt.Errorf("生成随机数失败: %w", err)
		}
		password[i] = charset[num.Int64()]
	}

	return string(password), nil
}

// GenerateByStrength 根据强度生成密码
func (g *PasswordGenerator) GenerateByStrength(strength PasswordStrength) (string, error) {
	var config *PasswordConfig

	switch strength {
	case PasswordStrengthWeak:
		config = &PasswordConfig{
			Length:         8,
			IncludeLower:   true,
			IncludeUpper:   true,
			IncludeDigits:  true,
			IncludeSpecial: false,
		}
	case PasswordStrengthMedium:
		config = &PasswordConfig{
			Length:         12,
			IncludeLower:   true,
			IncludeUpper:   true,
			IncludeDigits:  true,
			IncludeSpecial: true,
		}
	case PasswordStrengthStrong:
		config = &PasswordConfig{
			Length:         16,
			IncludeLower:   true,
			IncludeUpper:   true,
			IncludeDigits:  true,
			IncludeSpecial: true,
		}
	default:
		return "", fmt.Errorf("不支持的密码强度: %s", strength)
	}

	return g.Generate(config)
}

// ValidateStrength 验证密码强度
func (g *PasswordGenerator) ValidateStrength(password string, minStrength PasswordStrength) error {
	if len(password) < 8 {
		return fmt.Errorf("密码长度至少为 8 位")
	}

	hasLower := strings.ContainsAny(password, lowercaseLetters)
	hasUpper := strings.ContainsAny(password, uppercaseLetters)
	hasDigit := strings.ContainsAny(password, digits)
	hasSpecial := strings.ContainsAny(password, specialChars)

	// 计算密码强度
	var actualStrength PasswordStrength
	if len(password) >= 16 && hasLower && hasUpper && hasDigit && hasSpecial {
		actualStrength = PasswordStrengthStrong
	} else if len(password) >= 12 && hasLower && hasUpper && hasDigit {
		actualStrength = PasswordStrengthMedium
	} else if len(password) >= 8 && hasLower && hasUpper && hasDigit {
		actualStrength = PasswordStrengthWeak
	} else {
		return fmt.Errorf("密码不符合最低强度要求")
	}

	// 检查是否满足最低强度要求
	strengthOrder := map[PasswordStrength]int{
		PasswordStrengthWeak:   1,
		PasswordStrengthMedium: 2,
		PasswordStrengthStrong: 3,
	}

	if strengthOrder[actualStrength] < strengthOrder[minStrength] {
		return fmt.Errorf("密码强度不足,要求至少为 %s", minStrength)
	}

	return nil
}
