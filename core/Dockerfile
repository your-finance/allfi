# AllFi 全栈 Docker 构建
# 注意：构建上下文必须是项目根目录 (.)
# 参考 docker-compose.yml: services.backend.build.context: .

# ====================
# 阶段 1: 前端构建
# ====================
FROM node:20-alpine AS frontend
WORKDIR /build

# 复制 pnpm 依赖定义
COPY webapp/package.json webapp/pnpm-lock.yaml ./

# 安装 pnpm
RUN corepack enable && corepack prepare pnpm@latest --activate

# 安装依赖
RUN pnpm install --frozen-lockfile

# 复制源码
COPY webapp/ .

# 构建
ARG APP_VERSION=dev
RUN echo "${APP_VERSION}" > /VERSION
RUN pnpm build

# ====================
# 阶段 2: 后端构建 (含前端嵌入)
# ====================
FROM golang:1.24-alpine AS backend
WORKDIR /build

# 复制 go.mod 依赖
COPY core/go.mod core/go.sum ./
RUN go mod download

# 复制后端源码
COPY core/ .

# 复制前端产物到正确位置 (internal/statics/dist)
# 这使得 embed.FS 可以打包前端资源
COPY --from=frontend /build/dist ./internal/statics/dist

# 构建二进制
ARG APP_VERSION=dev
ARG BUILD_TIME=unknown
ARG GIT_COMMIT=unknown

# 增加共享内存，禁用 CGO，限制并行编译以降低内存消耗
RUN mount -t tmpfs -o size=2G tmpfs /tmp 2>/dev/null || true
RUN CGO_ENABLED=0 GOMAXPROCS=2 GOGC=20 go build \
  -p=1 \
  -ldflags="-s -w \
  -X your-finance/allfi/internal/version.Version=${APP_VERSION} \
  -X your-finance/allfi/internal/version.BuildTime=${BUILD_TIME} \
  -X your-finance/allfi/internal/version.GitCommit=${GIT_COMMIT}" \
  -o allfi cmd/server/main.go

# ====================
# 阶段 3: 运行时环境
# ====================
FROM alpine:3.21

# 安装基础工具 (证书、时区)
RUN apk add --no-cache ca-certificates tzdata

# 创建非特权用户
RUN addgroup -S allfi && adduser -S allfi -G allfi

WORKDIR /app

# 从构建阶段复制二进制
COPY --from=backend /build/allfi .
# 复制配置模板
COPY --from=backend /build/manifest/config/config.yaml manifest/config/config.yaml

# 设置权限
RUN mkdir -p /app/data /app/logs && chown -R allfi:allfi /app

# 切换用户
USER allfi

# 暴露端口 (只需暴露后端端口，前端通过同一端口访问)
EXPOSE 8080

# 数据持久化
VOLUME ["/app/data"]

ENTRYPOINT ["./allfi"]
