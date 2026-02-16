#!/bin/bash
# AllFi Docker 一键部署脚本
# 用法: curl -sSL https://raw.githubusercontent.com/your-finance/allfi/main/deploy/docker-deploy.sh | bash
# 或在本地: bash deploy/docker-deploy.sh

set -e

# 颜色定义
CYAN='\033[36m'
GREEN='\033[32m'
YELLOW='\033[33m'
RED='\033[31m'
RESET='\033[0m'

echo ""
echo -e "${CYAN}╔══════════════════════════════════════════════╗${RESET}"
echo -e "${CYAN}║       AllFi — 全资产聚合平台 一键部署        ║${RESET}"
echo -e "${CYAN}╚══════════════════════════════════════════════╝${RESET}"
echo ""

# ==================== 前置检查 ====================

echo -e "${CYAN}>>> 检查前置条件...${RESET}"

# 检查 Docker
if ! command -v docker >/dev/null 2>&1; then
    echo -e "${RED}  ✗ 未检测到 Docker，请先安装 Docker${RESET}"
    echo "    安装指南: https://docs.docker.com/engine/install/"
    exit 1
fi
echo -e "${GREEN}  ✓ Docker $(docker --version | grep -oP '\d+\.\d+\.\d+' || echo 'detected')${RESET}"

# 检查 Docker Compose (v2 插件优先，回退 docker-compose)
COMPOSE_CMD=""
if docker compose version >/dev/null 2>&1; then
    COMPOSE_CMD="docker compose"
    echo -e "${GREEN}  ✓ Docker Compose (plugin) detected${RESET}"
elif command -v docker-compose >/dev/null 2>&1; then
    COMPOSE_CMD="docker-compose"
    echo -e "${GREEN}  ✓ docker-compose (standalone) detected${RESET}"
else
    echo -e "${RED}  ✗ 未检测到 Docker Compose，请先安装${RESET}"
    echo "    安装指南: https://docs.docker.com/compose/install/"
    exit 1
fi

# 检查 openssl
if ! command -v openssl >/dev/null 2>&1; then
    echo -e "${RED}  ✗ 未检测到 openssl，请先安装${RESET}"
    exit 1
fi
echo -e "${GREEN}  ✓ openssl detected${RESET}"

echo ""

# ==================== 判断运行模式 ====================

# 如果已在 allfi 项目根目录（有 docker-compose.yml 和 core/ webapp/），直接使用
if [ -f "docker-compose.yml" ] && [ -d "core" ] && [ -d "webapp" ]; then
    echo -e "${CYAN}>>> 检测到已在 AllFi 项目目录，使用本地文件...${RESET}"
    PROJECT_DIR="."
else
    # 远程模式：克隆仓库
    echo -e "${CYAN}>>> 克隆 AllFi 仓库...${RESET}"
    if [ -d "allfi" ]; then
        echo -e "${YELLOW}  ⊘ allfi 目录已存在，跳过克隆${RESET}"
    else
        git clone https://github.com/your-finance/allfi.git
        echo -e "${GREEN}  ✓ 仓库克隆完成${RESET}"
    fi
    PROJECT_DIR="allfi"
fi

cd "$PROJECT_DIR"

echo ""

# ==================== 生成 .env ====================

echo -e "${CYAN}>>> 配置环境变量...${RESET}"

if [ -f ".env" ]; then
    echo -e "${YELLOW}  ⊘ .env 已存在，跳过（如需重新生成，请先删除 .env）${RESET}"
else
    if [ ! -f ".env.example" ]; then
        echo -e "${RED}  ✗ 缺少 .env.example 文件${RESET}"
        exit 1
    fi

    cp .env.example .env

    # 自动生成 MASTER_KEY
    MASTER_KEY=$(openssl rand -base64 32)
    if [ "$(uname)" = "Darwin" ]; then
        sed -i '' "s|CHANGE_ME_USE_openssl_rand_base64_32|${MASTER_KEY}|" .env
    else
        sed -i "s|CHANGE_ME_USE_openssl_rand_base64_32|${MASTER_KEY}|" .env
    fi

    echo -e "${GREEN}  ✓ 已生成 .env${RESET}"
    echo -e "${GREEN}  ✓ ALLFI_MASTER_KEY 已自动生成${RESET}"
fi

echo ""

# ==================== 启动服务 ====================

echo -e "${CYAN}>>> 启动 Docker 服务...${RESET}"
echo -e "  ${YELLOW}首次启动需要构建镜像，可能需要几分钟...${RESET}"
echo ""

# 获取版本号（VERSION 文件由 CI 维护，是权威版本源）
if [ -f VERSION ]; then
    export ALLFI_VERSION=$(cat VERSION)
elif git describe --tags --abbrev=0 >/dev/null 2>&1; then
    export ALLFI_VERSION=$(git describe --tags --abbrev=0 | sed 's/^v//')
else
    export ALLFI_VERSION="dev"
fi
export GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo 'unknown')
export BUILD_TIME=$(date +%Y-%m-%dT%H:%M:%S)

echo -e "  版本: ${GREEN}${ALLFI_VERSION}${RESET}"
echo -e "  提交: ${GREEN}${GIT_COMMIT}${RESET}"
echo ""

$COMPOSE_CMD up -d --build

echo ""
echo -e "${GREEN}╔══════════════════════════════════════════════╗${RESET}"
echo -e "${GREEN}║           AllFi 部署完成！                   ║${RESET}"
echo -e "${GREEN}╚══════════════════════════════════════════════╝${RESET}"
echo ""
echo -e "  前端: ${CYAN}http://localhost:${FRONTEND_PORT:-3174}${RESET}"
echo -e "  后端: ${CYAN}http://localhost:${SERVER_PORT:-8080}${RESET}"
echo -e "  API 文档: ${CYAN}http://localhost:${SERVER_PORT:-8080}/api/v1/docs${RESET}"
echo ""
echo -e "  首次访问需设置 ${YELLOW}PIN 码${RESET}（4-8 位数字）"
echo ""
echo -e "${CYAN}常用命令:${RESET}"
echo -e "  查看日志:   ${GREEN}$COMPOSE_CMD logs -f${RESET}"
echo -e "  停止服务:   ${GREEN}$COMPOSE_CMD down${RESET}"
echo -e "  重启服务:   ${GREEN}$COMPOSE_CMD restart${RESET}"
echo -e "  重新构建:   ${GREEN}$COMPOSE_CMD up -d --build${RESET}"
echo ""
