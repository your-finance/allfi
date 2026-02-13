#!/bin/bash
# AllFi 数据恢复脚本
# 用法: ./scripts/restore.sh <备份文件路径>

set -euo pipefail

BACKUP_FILE="${1:-}"
DB_PATH="${ALLFI_DB_PATH:-./core/data/allfi.db}"

if [ -z "$BACKUP_FILE" ]; then
  echo "用法: $0 <备份文件路径>"
  echo ""
  echo "可用备份:"
  ls -lt ./backups/allfi_backup_* 2>/dev/null | head -10 || echo "  无备份文件"
  exit 1
fi

if [ ! -f "$BACKUP_FILE" ]; then
  echo "错误: 备份文件不存在: $BACKUP_FILE"
  exit 1
fi

echo "警告: 这将覆盖当前数据库！"
read -p "确认恢复? (y/N) " confirm
[ "$confirm" = "y" ] || exit 0

# 解密（如果是加密备份）
TEMP_TAR=""
if [[ "$BACKUP_FILE" == *.enc ]]; then
  TEMP_TAR=$(mktemp)
  openssl enc -aes-256-cbc -pbkdf2 -d -salt \
    -in "$BACKUP_FILE" -out "$TEMP_TAR"
  BACKUP_FILE="$TEMP_TAR"
fi

# 停止服务
echo "停止 AllFi 服务..."
docker-compose stop backend 2>/dev/null || true

# 备份当前数据库
if [ -f "$DB_PATH" ]; then
  cp "$DB_PATH" "${DB_PATH}.before_restore"
  echo "已备份当前数据库到: ${DB_PATH}.before_restore"
fi

# 恢复：解压到临时目录，再拷贝数据库
TEMP_DIR=$(mktemp -d)
tar -xzf "$BACKUP_FILE" -C "$TEMP_DIR"
DB_FILE=$(find "$TEMP_DIR" -name "*.db" -type f | head -1)

if [ -z "$DB_FILE" ]; then
  echo "错误: 备份中未找到数据库文件"
  rm -rf "$TEMP_DIR"
  [ -n "$TEMP_TAR" ] && rm -f "$TEMP_TAR"
  exit 1
fi

mkdir -p "$(dirname "$DB_PATH")"
cp "$DB_FILE" "$DB_PATH"

# 清理
rm -rf "$TEMP_DIR"
[ -n "$TEMP_TAR" ] && rm -f "$TEMP_TAR"

echo "恢复完成！请重启服务: docker-compose up -d"
