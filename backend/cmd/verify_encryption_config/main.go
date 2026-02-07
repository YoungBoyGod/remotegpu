package main

import (
	"fmt"
	"log"

	"github.com/YoungBoyGod/remotegpu/config"
	"github.com/YoungBoyGod/remotegpu/pkg/crypto"
)

func main() {
	// 加载配置
	if err := config.LoadConfig("config/config.yaml"); err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	fmt.Println("=== 加密密钥配置验证 ===")
	fmt.Printf("配置文件中的密钥: %s\n", config.GlobalConfig.Encryption.Key)
	fmt.Println()

	// 测试加密
	testPassword := "test123"
	encrypted, err := crypto.EncryptAES256GCM(testPassword)
	if err != nil {
		log.Fatalf("加密失败: %v", err)
	}

	fmt.Printf("原始密码: %s\n", testPassword)
	fmt.Printf("加密后长度: %d 字节\n", len(encrypted))
	fmt.Println()

	// 测试解密
	decrypted, err := crypto.DecryptAES256GCM(encrypted)
	if err != nil {
		log.Fatalf("解密失败: %v", err)
	}

	fmt.Printf("解密后密码: %s\n", decrypted)

	if decrypted == testPassword {
		fmt.Println("\n✅ 配置文件中的加密密钥工作正常！")
	} else {
		fmt.Println("\n❌ 加密/解密失败")
	}
}
