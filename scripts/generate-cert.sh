#!/bin/bash
# AllFi 自签名 TLS 证书生成脚本（开发/自托管用）
# 用法: ./scripts/generate-cert.sh

set -euo pipefail

CERT_DIR="${CERT_DIR:-./certs}"
mkdir -p "$CERT_DIR"

if [ -f "$CERT_DIR/server.crt" ]; then
  echo "证书已存在: $CERT_DIR/server.crt"
  echo "删除后重新运行以重新生成"
  exit 0
fi

echo "生成自签名 TLS 证书..."

openssl req -x509 -newkey rsa:4096 \
  -keyout "$CERT_DIR/server.key" \
  -out "$CERT_DIR/server.crt" \
  -days 365 -nodes \
  -subj "/CN=localhost/O=AllFi/C=US" \
  -addext "subjectAltName=DNS:localhost,IP:127.0.0.1"

chmod 600 "$CERT_DIR/server.key"
chmod 644 "$CERT_DIR/server.crt"

echo "证书已生成:"
echo "  证书: $CERT_DIR/server.crt"
echo "  密钥: $CERT_DIR/server.key"
echo "有效期: 365 天"
echo ""
echo "启用 HTTPS: 取消注释 webapp/nginx.conf 中的 HTTPS 配置"
