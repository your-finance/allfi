#!/usr/bin/env bash
# AllFi 快速启动脚本
# 用法:
#   bash scripts/quickstart.sh          # 完整启动（前端 + 后端）
#   bash scripts/quickstart.sh --mock   # 仅前端 Mock 模式
#   bash scripts/quickstart.sh --check  # 仅检测依赖
#   bash scripts/quickstart.sh --docker # Docker 模式启动

set -euo pipefail

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
CYAN='\033[0;36m'
RESET='\033[0m'

# 项目根目录
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

# 最低版本要求
REQUIRED_GO_MAJOR=1
REQUIRED_GO_MINOR=24
REQUIRED_NODE_MAJOR=18

# ==================== 工具函数 ====================

info()    { echo -e "${CYAN}>>> $*${RESET}"; }
success() { echo -e "${GREEN}  ✓ $*${RESET}"; }
warn()    { echo -e "${YELLOW}  ⚠ $*${RESET}"; }
error()   { echo -e "${RED}  ✗ $*${RESET}"; }

# 版本号比较：检查 $1 >= $2.$3
check_version() {
    local version="$1"
    local required_major="$2"
    local required_minor="$3"
    local major minor
    major=$(echo "$version" | grep -oE '^[0-9]+')
    minor=$(echo "$version" | grep -oE '\.[0-9]+' | head -1 | tr -d '.')
    if [ "$major" -gt "$required_major" ] 2>/dev/null; then return 0; fi
    if [ "$major" -eq "$required_major" ] && [ "$minor" -ge "$required_minor" ] 2>/dev/null; then return 0; fi
    return 1
}

# ==================== 依赖检测 ====================

check_dependencies() {
    info "检测开发环境依赖..."
    local all_ok=true

    # Go
    if command -v go &>/dev/null; then
        local go_version
        go_version=$(go version | grep -oE '[0-9]+\.[0-9]+(\.[0-9]+)?')
        if check_version "$go_version" $REQUIRED_GO_MAJOR $REQUIRED_GO_MINOR; then
            success "Go $go_version"
        else
            error "Go $go_version（需要 >= $REQUIRED_GO_MAJOR.$REQUIRED_GO_MINOR）"
            all_ok=false
        fi
    else
        error "Go 未安装（需要 >= $REQUIRED_GO_MAJOR.$REQUIRED_GO_MINOR）"
        echo "        安装: https://golang.org/dl/"
        all_ok=false
    fi

    # Node.js
    if command -v node &>/dev/null; then
        local node_version
        node_version=$(node --version | tr -d 'v')
        if check_version "$node_version" $REQUIRED_NODE_MAJOR 0; then
            success "Node.js $node_version"
        else
            error "Node.js $node_version（需要 >= $REQUIRED_NODE_MAJOR）"
            all_ok=false
        fi
    else
        error "Node.js 未安装（需要 >= $REQUIRED_NODE_MAJOR）"
        echo "        安装: https://nodejs.org/"
        all_ok=false
    fi

    # pnpm
    if command -v pnpm &>/dev/null; then
        success "pnpm $(pnpm --version)"
    else
        error "pnpm 未安装"
        echo "        安装: npm install -g pnpm"
        all_ok=false
    fi

    # Docker（可选）
    if command -v docker &>/dev/null; then
        success "Docker $(docker --version | grep -oE '[0-9]+\.[0-9]+\.[0-9]+')"
    else
        warn "Docker 未安装（仅 Docker 部署需要）"
    fi

    # Docker Compose（可选）
    if command -v docker-compose &>/dev/null || docker compose version &>/dev/null 2>&1; then
        success "Docker Compose 可用"
    else
        warn "Docker Compose 未安装（仅 Docker 部署需要）"
    fi

    # GCC（后端编译需要）
    if command -v gcc &>/dev/null; then
        success "GCC $(gcc --version | head -1 | grep -oE '[0-9]+\.[0-9]+\.[0-9]+')"
    else
        warn "GCC 未安装（后端编译需要，macOS 可运行 xcode-select --install）"
    fi

    echo ""
    if $all_ok; then
        success "所有必需依赖已就绪！"
    else
        error "部分依赖缺失，请先安装后再试"
        return 1
    fi
}

# ==================== 代理检测 ====================

check_proxy() {
    info "检测网络环境..."

    # 检测是否已配置代理
    if [ -n "${HTTP_PROXY:-}" ] || [ -n "${HTTPS_PROXY:-}" ]; then
        success "已检测到代理配置: ${HTTP_PROXY:-${HTTPS_PROXY:-}}"
        return
    fi

    # 检测 Go 模块代理
    local goproxy="${GOPROXY:-}"
    if [ -n "$goproxy" ] && [ "$goproxy" != "https://proxy.golang.org,direct" ]; then
        success "Go 模块代理: $goproxy"
    else
        warn "未配置 Go 模块代理"
        echo "        如果下载 Go 依赖较慢，建议设置:"
        echo "        export GOPROXY=https://goproxy.cn,direct"
    fi

    # 检测 npm/pnpm 镜像源
    if command -v pnpm &>/dev/null; then
        local registry
        registry=$(pnpm config get registry 2>/dev/null || echo "")
        if echo "$registry" | grep -qi "npmmirror\|taobao\|cnpm"; then
            success "pnpm 镜像源: $registry"
        else
            warn "pnpm 使用默认源 ($registry)"
            echo "        如果安装依赖较慢，建议设置:"
            echo "        pnpm config set registry https://registry.npmmirror.com"
        fi
    fi

    echo ""
    echo "  详细代理配置指南: docs/guides/proxy-guide.md"
    echo ""
}

# ==================== 环境初始化 ====================

init_env() {
    info "初始化环境配置..."

    cd "$PROJECT_ROOT"

    if [ ! -f .env ]; then
        cp .env.example .env
        # 生成随机 MASTER_KEY
        local master_key
        master_key=$(openssl rand -base64 32)
        if [ "$(uname)" = "Darwin" ]; then
            sed -i '' "s|CHANGE_ME_USE_openssl_rand_base64_32|$master_key|" .env
        else
            sed -i "s|CHANGE_ME_USE_openssl_rand_base64_32|$master_key|" .env
        fi
        success "已生成 .env（MASTER_KEY 已自动填入）"
    else
        success ".env 已存在，跳过"
    fi
}

# ==================== 安装依赖 ====================

install_deps() {
    info "安装项目依赖..."

    # 后端
    echo -e "  ${CYAN}→ 安装 Go 依赖...${RESET}"
    cd "$PROJECT_ROOT/core"
    go mod download && go mod verify
    success "Go 依赖安装完成"

    # 前端
    echo -e "  ${CYAN}→ 安装前端依赖...${RESET}"
    cd "$PROJECT_ROOT/webapp"
    pnpm install
    success "前端依赖安装完成"
}

# ==================== 启动服务 ====================

start_full() {
    info "启动 AllFi（前端 + 后端）..."
    echo ""
    echo "  后端 API:  http://localhost:8080"
    echo "  前端页面:  http://localhost:3174"
    echo "  Swagger:   http://localhost:8080/swagger/"
    echo "  按 Ctrl+C 停止所有服务"
    echo ""

    cd "$PROJECT_ROOT"
    trap 'kill 0; exit 0' INT TERM
    (cd core && go run cmd/server/main.go) &
    sleep 2

    # 等待后端启动
    local retries=0
    while ! curl -s http://localhost:8080/api/v1/health &>/dev/null; do
        retries=$((retries + 1))
        if [ $retries -gt 15 ]; then
            warn "后端启动超时，但前端仍可使用"
            break
        fi
        sleep 1
    done

    if curl -s http://localhost:8080/api/v1/health &>/dev/null; then
        success "后端已就绪"
    fi

    cd "$PROJECT_ROOT/webapp"
    pnpm dev
    wait
}

start_mock() {
    info "启动 AllFi Mock 模式（纯前端）..."
    echo ""
    echo "  前端页面: http://localhost:3174"
    echo "  数据来源: 模拟数据（无需后端）"
    echo "  按 Ctrl+C 停止"
    echo ""

    cd "$PROJECT_ROOT/webapp"
    pnpm dev:mock
}

start_docker() {
    info "启动 AllFi Docker 模式..."

    cd "$PROJECT_ROOT"
    init_env
    docker-compose up -d --build

    echo ""
    success "Docker 服务已启动"
    echo "  前端: http://localhost:${FRONTEND_PORT:-5173}"
    echo "  后端: http://localhost:${SERVER_PORT:-8080}"
    echo ""
    echo "  查看日志: docker-compose logs -f"
    echo "  停止服务: docker-compose down"
}

# ==================== 主逻辑 ====================

main() {
    echo ""
    echo "╔══════════════════════════════════════╗"
    echo "║     AllFi 快速启动脚本 v1.0          ║"
    echo "╚══════════════════════════════════════╝"
    echo ""

    local mode="${1:-full}"

    case "$mode" in
        --check|-c)
            check_dependencies
            check_proxy
            ;;
        --mock|-m)
            check_dependencies || exit 1
            init_env
            cd "$PROJECT_ROOT/webapp"
            if [ ! -d "node_modules" ]; then
                pnpm install
            fi
            start_mock
            ;;
        --docker|-d)
            if ! command -v docker &>/dev/null; then
                error "Docker 未安装"
                exit 1
            fi
            init_env
            start_docker
            ;;
        --help|-h)
            echo "用法: bash scripts/quickstart.sh [选项]"
            echo ""
            echo "选项:"
            echo "  (无参数)     完整启动（前端 + 后端）"
            echo "  --mock, -m   纯前端 Mock 模式（无需后端）"
            echo "  --docker, -d Docker 模式启动"
            echo "  --check, -c  仅检测依赖环境"
            echo "  --help, -h   显示帮助信息"
            echo ""
            echo "示例:"
            echo "  bash scripts/quickstart.sh          # 首次完整启动"
            echo "  bash scripts/quickstart.sh --mock   # 快速体验 UI"
            echo "  bash scripts/quickstart.sh --check  # 检查环境"
            ;;
        *)
            check_dependencies || exit 1
            check_proxy
            init_env
            install_deps
            start_full
            ;;
    esac
}

main "$@"
