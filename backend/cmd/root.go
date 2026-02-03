package main

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	configPath string
	mode       string
)

var rootCmd = &cobra.Command{
	Use:   "remotegpu",
	Short: "RemoteGPU Backend Service CLI",
	Long:  `RemoteGPU 后端服务与管理工具命令行接口`,
}

// Execute 执行根命令
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	// 全局标志
	rootCmd.PersistentFlags().StringVar(&configPath, "config", "./config/config.yaml", "配置文件路径")
	rootCmd.PersistentFlags().StringVar(&mode, "mode", "debug", "运行模式: debug, release, test")
}