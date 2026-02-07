# RemoteGPU 项目统一 Makefile
# 用法: make help

.PHONY: help
help: ## 显示帮助信息
	@echo "RemoteGPU 项目构建命令"
	@echo "========================"
	@echo ""
	@echo "全局命令:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

# ==================== 全局 ====================

.PHONY: all build test lint clean
all: build ## 构建所有模块

build: backend-build agent-build frontend-build ## 构建所有模块

test: backend-test agent-test frontend-type-check ## 运行所有测试

lint: backend-lint frontend-lint ## 运行所有 lint 检查

clean: backend-clean agent-clean frontend-clean ## 清理所有构建产物

# ==================== 后端 ====================

.PHONY: backend-build backend-run backend-dev backend-test backend-lint backend-clean

backend-build: ## 编译后端
	@echo ">>> 编译后端..."
	@cd backend && go build -o bin/remotegpu ./cmd

backend-run: ## 运行后端
	@cd backend && go run ./cmd server

backend-dev: ## 后端开发模式（热更新）
	@cd backend && air

backend-test: ## 运行后端测试
	@echo ">>> 运行后端测试..."
	@cd backend && go test ./...

backend-lint: ## 后端代码检查
	@echo ">>> 后端 lint..."
	@cd backend && go vet ./...

backend-clean: ## 清理后端构建产物
	@rm -rf backend/bin/ backend/tmp/ backend/logs/

# ==================== Agent ====================

.PHONY: agent-build agent-test agent-clean

agent-build: ## 编译 Agent
	@echo ">>> 编译 Agent..."
	@cd agent && go build -o remotegpu-agent ./cmd

agent-test: ## 运行 Agent 测试
	@echo ">>> 运行 Agent 测试..."
	@cd agent && go test ./...

agent-clean: ## 清理 Agent 构建产物
	@rm -f agent/remotegpu-agent

# ==================== 前端 ====================

.PHONY: frontend-install frontend-build frontend-dev frontend-type-check frontend-lint frontend-clean

frontend-install: ## 安装前端依赖
	@echo ">>> 安装前端依赖..."
	@cd frontend && npm install

frontend-build: ## 构建前端
	@echo ">>> 构建前端..."
	@cd frontend && npm run build

frontend-dev: ## 前端开发服务器
	@cd frontend && npm run dev

frontend-type-check: ## 前端类型检查
	@echo ">>> 前端类型检查..."
	@cd frontend && npm run type-check

frontend-lint: ## 前端 lint 检查
	@echo ">>> 前端 lint..."
	@cd frontend && npm run lint

frontend-clean: ## 清理前端构建产物
	@rm -rf frontend/dist/ frontend/node_modules/.cache/

# ==================== 数据库 ====================

.PHONY: db-migrate db-migrate-dry

DB_HOST ?= 192.168.10.210
DB_PORT ?= 5432
DB_USER ?= remotegpu_user
DB_NAME ?= remotegpu

db-migrate: ## 执行数据库迁移（按编号顺序）
	@echo ">>> 执行数据库迁移..."
	@for f in $$(ls backend/sql/*.sql | sort -V); do \
		echo "  执行: $$f"; \
		PGPASSWORD=$(DB_PASSWORD) psql -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER) -d $(DB_NAME) -f "$$f"; \
	done
	@echo ">>> 迁移完成"

db-migrate-dry: ## 预览迁移脚本（不执行）
	@echo ">>> 迁移脚本列表:"
	@ls -1 backend/sql/*.sql | sort -V | while read f; do \
		echo "  $$f"; \
	done

# ==================== Docker ====================

.PHONY: infra-up infra-down infra-status

infra-up: ## 启动基础设施服务
	@bash docker-compose/start-all.sh

infra-down: ## 停止基础设施服务
	@bash docker-compose/stop-all.sh

infra-status: ## 检查基础设施状态
	@bash docker-compose/check-status.sh

# ==================== 测试环境 ====================

.PHONY: test-env-up test-env-down test-env-build

test-env-build: ## 构建测试机器镜像
	@echo ">>> 构建测试环境镜像..."
	@cd docker-compose/test-env && docker compose build

test-env-up: ## 启动测试机器（模拟 GPU 节点）
	@echo ">>> 启动测试环境..."
	@cd docker-compose/test-env && docker compose up -d --build

test-env-down: ## 停止测试机器
	@echo ">>> 停止测试环境..."
	@cd docker-compose/test-env && docker compose down

# ==================== 开发快捷命令 ====================

.PHONY: dev setup

dev: ## 同时启动后端和前端开发服务器
	@echo ">>> 启动后端..."
	@cd backend && go run ./cmd server &
	@echo ">>> 启动前端..."
	@cd frontend && npm run dev

setup: frontend-install ## 初始化开发环境
	@echo ">>> 开发环境初始化完成"
