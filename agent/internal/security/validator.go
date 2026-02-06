package security

import (
	"fmt"
	"strings"
)

// Validator 命令校验器
type Validator struct {
	allowedCommands []string
	blockedPatterns []string
}

// NewValidator 创建校验器
func NewValidator(allowed []string, blocked []string) *Validator {
	return &Validator{
		allowedCommands: allowed,
		blockedPatterns: blocked,
	}
}

// Validate 校验命令是否允许执行
func (v *Validator) Validate(command string, args []string) error {
	// 检查白名单（如果配置了）
	if len(v.allowedCommands) > 0 {
		allowed := false
		for _, cmd := range v.allowedCommands {
			if command == cmd {
				allowed = true
				break
			}
		}
		if !allowed {
			return fmt.Errorf("command %q not in allowed list", command)
		}
	}

	// 检查黑名单
	fullCmd := command + " " + strings.Join(args, " ")
	for _, pattern := range v.blockedPatterns {
		if strings.Contains(fullCmd, pattern) {
			return fmt.Errorf("command matches blocked pattern %q", pattern)
		}
	}

	return nil
}

// Enabled 检查校验器是否启用
func (v *Validator) Enabled() bool {
	return len(v.allowedCommands) > 0 || len(v.blockedPatterns) > 0
}
