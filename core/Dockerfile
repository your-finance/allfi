# AllFi 后端 Docker 多阶段构建
# 构建阶段：编译 Go 二进制
FROM golang:1.24-alpine AS builder
RUN apk add --no-cache gcc musl-dev sqlite-dev
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ARG APP_VERSION=dev
ARG BUILD_TIME=unknown
ARG GIT_COMMIT=unknown
RUN CGO_ENABLED=1 go build \
  -ldflags="-s -w -X your-finance/allfi/internal/version.Version=${APP_VERSION} -X your-finance/allfi/internal/version.BuildTime=${BUILD_TIME} -X your-finance/allfi/internal/version.GitCommit=${GIT_COMMIT}" \
  -o allfi cmd/server/main.go

# 运行阶段：最小化镜像
FROM alpine:3.21
RUN apk add --no-cache ca-certificates sqlite-libs tzdata
RUN addgroup -S allfi && adduser -S allfi -G allfi
WORKDIR /app
COPY --from=builder /build/allfi .
COPY --from=builder /build/manifest/config/config.yaml manifest/config/config.yaml
RUN mkdir -p /app/data /app/logs && chown -R allfi:allfi /app
USER allfi
EXPOSE 8080
VOLUME ["/app/data"]
ENTRYPOINT ["./allfi"]
