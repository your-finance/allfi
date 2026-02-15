# AllFi Makefile — 统一开发命令入口
# 用法: make help

.PHONY: help setup dev dev-mock dev-backend dev-frontend build build-backend build-frontend \
        docker docker-build docker-down docker-logs clean health swagger

# 默认目标
.DEFAULT_GOAL := help

# 版本信息
VERSION := $(shell cat VERSION 2>/dev/null || echo "dev")
BUILD_TIME := $(shell date +%Y-%m-%dT%H:%M:%S)
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
LDFLAGS := -s -w \
  -X your-finance/allfi/internal/version.Version=$(VERSION) \
  -X your-finance/allfi/internal/version.BuildTime=$(BUILD_TIME) \
  -X your-finance/allfi/internal/version.GitCommit=$(GIT_COMMIT)

# 颜色定义
CYAN  := \033[36m
GREEN := \033[32m
YELLOW := \033[33m
RESET := \033[0m

help: ## 显示所有可用命令
	@echo ""
	@echo "$(CYAN)AllFi 开发命令$(RESET)"
	@echo "============================================"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  $(GREEN)%-18s$(RESET) %s\n", $$1, $$2}'
	@echo ""

# ==================== 初始化 ====================

setup: ## 一键初始化（生成 .env + 安装前后端依赖）
	@echo "$(CYAN)>>> 初始化 AllFi 开发环境...$(RESET)"
	@# 生成 .env
	@if [ ! -f .env ]; then \
		cp .env.example .env; \
		MASTER_KEY=$$(openssl rand -base64 32); \
		if [ "$$(uname)" = "Darwin" ]; then \
			sed -i '' "s|CHANGE_ME_USE_openssl_rand_base64_32|$$MASTER_KEY|" .env; \
		else \
			sed -i "s|CHANGE_ME_USE_openssl_rand_base64_32|$$MASTER_KEY|" .env; \
		fi; \
		echo "$(GREEN)  ✓ 已生成 .env（MASTER_KEY 已自动填入）$(RESET)"; \
	else \
		echo "$(YELLOW)  ⊘ .env 已存在，跳过$(RESET)"; \
	fi
	@# 安装后端依赖
	@echo "$(CYAN)  → 安装 Go 依赖...$(RESET)"
	@cd core && go mod download && go mod verify
	@echo "$(GREEN)  ✓ Go 依赖安装完成$(RESET)"
	@# 安装前端依赖
	@echo "$(CYAN)  → 安装前端依赖...$(RESET)"
	@cd webapp && pnpm install
	@echo "$(GREEN)  ✓ 前端依赖安装完成$(RESET)"
	@echo ""
	@echo "$(GREEN)>>> 初始化完成！$(RESET)"
	@echo "  运行 $(CYAN)make dev$(RESET) 启动开发服务器"
	@echo "  运行 $(CYAN)make dev-mock$(RESET) 纯前端 Mock 模式（无需后端）"

# ==================== 开发 ====================

dev: ## 同时启动前后端开发服务器
	@echo "$(CYAN)>>> 启动 AllFi 开发环境...$(RESET)"
	@echo "  后端: http://localhost:8080"
	@echo "  前端: http://localhost:3174"
	@echo "  Swagger: http://localhost:8080/swagger/"
	@echo "  按 Ctrl+C 停止所有服务"
	@echo ""
	@trap 'kill 0; exit 0' INT TERM; \
		cd core && go run -ldflags="$(LDFLAGS)" cmd/server/main.go & \
		cd webapp && pnpm dev; \
		wait

dev-mock: ## 纯前端 Mock 模式（无需后端，快速体验 UI）
	@echo "$(CYAN)>>> 启动 Mock 模式...$(RESET)"
	@echo "  前端: http://localhost:3174（Mock 数据）"
	@echo "  无需启动后端，所有数据为模拟数据"
	@echo ""
	@cd webapp && pnpm dev:mock

dev-backend: ## 仅启动后端开发服务器
	@echo "$(CYAN)>>> 启动后端...$(RESET)"
	@echo "  API: http://localhost:8080"
	@echo "  Swagger: http://localhost:8080/swagger/"
	@echo ""
	@cd core && go run -ldflags="$(LDFLAGS)" cmd/server/main.go

dev-frontend: ## 仅启动前端开发服务器
	@echo "$(CYAN)>>> 启动前端...$(RESET)"
	@echo "  前端: http://localhost:3174"
	@echo ""
	@cd webapp && pnpm dev

# ==================== 构建 ====================

build: build-backend build-frontend ## 构建前后端

build-backend: ## 构建后端二进制文件
	@echo "$(CYAN)>>> 构建后端 v$(VERSION)...$(RESET)"
	@cd core && CGO_ENABLED=1 go build -ldflags="$(LDFLAGS)" -o allfi cmd/server/main.go
	@echo "$(GREEN)  ✓ 后端构建完成: core/allfi$(RESET)"

build-frontend: ## 构建前端生产版本
	@echo "$(CYAN)>>> 构建前端...$(RESET)"
	@cd webapp && pnpm build
	@echo "$(GREEN)  ✓ 前端构建完成: webapp/dist/$(RESET)"

# ==================== Docker ====================

docker: ## Docker Compose 启动（后台运行）
	@echo "$(CYAN)>>> Docker Compose 启动...$(RESET)"
	@docker-compose up -d
	@echo "$(GREEN)  ✓ 服务已启动$(RESET)"
	@echo "  前端: http://localhost:$${FRONTEND_PORT:-5173}"
	@echo "  后端: http://localhost:$${SERVER_PORT:-8080}"

docker-build: ## Docker Compose 重新构建并启动
	@echo "$(CYAN)>>> Docker Compose 重新构建...$(RESET)"
	@docker-compose up -d --build

docker-down: ## Docker Compose 停止
	@docker-compose down
	@echo "$(GREEN)  ✓ 服务已停止$(RESET)"

docker-logs: ## 查看 Docker 日志
	@docker-compose logs -f

# ==================== 工具 ====================

health: ## 健康检查
	@echo "$(CYAN)>>> 健康检查...$(RESET)"
	@curl -s http://localhost:8080/api/v1/health | python3 -m json.tool 2>/dev/null || \
		echo "$(YELLOW)  ✗ 后端未运行或无法连接$(RESET)"

swagger: ## 打开 Swagger UI
	@echo "$(CYAN)>>> 打开 Swagger UI...$(RESET)"
	@open http://localhost:8080/swagger/ 2>/dev/null || \
		xdg-open http://localhost:8080/swagger/ 2>/dev/null || \
		echo "  请手动访问: http://localhost:8080/swagger/"

clean: ## 清理构建产物
	@echo "$(CYAN)>>> 清理构建产物...$(RESET)"
	@rm -f core/allfi
	@rm -rf webapp/dist
	@echo "$(GREEN)  ✓ 清理完成$(RESET)"
