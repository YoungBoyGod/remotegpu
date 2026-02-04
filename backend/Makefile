.PHONY: help build run dev clean test

help:
	@echo "可用命令:"
	@echo "  make build   - 编译项目"
	@echo "  make run     - 运行项目"
	@echo "  make dev     - 开发模式（热更新）"
	@echo "  make clean   - 清理编译文件"
	@echo "  make test    - 运行测试"

build:
	@echo "编译项目..."
	@go build -o bin/remotegpu ./cmd

run:
	@echo "运行项目..."
	@go run ./cmd server

dev:
	@echo "启动开发模式（热更新）..."
	@air

clean:
	@echo "清理编译文件..."
	@rm -rf bin/ tmp/ logs/

test:
	@echo "运行测试..."
	@go test -v ./...
