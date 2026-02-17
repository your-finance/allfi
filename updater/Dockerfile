# AllFi Updater Sidecar
# 轻量级更新执行器，负责 Docker Compose 模式下的版本更新和回滚
FROM golang:1.24-alpine AS builder
WORKDIR /build
COPY go.mod go.sum* ./
RUN go mod download 2>/dev/null; true
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o allfi-updater .

FROM alpine:3.21
RUN apk add --no-cache ca-certificates git docker-cli docker-cli-compose curl
# 配置 Git 强制使用 HTTPS 协议，避免 SSH 相关问题
RUN git config --global url."https://github.com/".insteadOf git@github.com:
RUN git config --global url."https://".insteadOf git://
WORKDIR /app
COPY --from=builder /build/allfi-updater .
EXPOSE 8081
ENTRYPOINT ["/app/allfi-updater"]
