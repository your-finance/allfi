#!/bin/bash
# AllFi 数据备份脚本（支持 AES 加密）
# 用法: ./scripts/backup.sh
# 加密备份: ENCRYPT_BACKUP=true ./scripts/backup.sh

set -euo pipefail

BACKUP_DIR="${BACKUP_DIR:-./backups}"
DB_PATH="${ALLFI_DB_PATH:-./core/data/allfi.db}"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
BACKUP_FILE="${BACKUP_DIR}/allfi_backup_${TIMESTAMP}.tar.gz"

mkdir -p "$BACKUP_DIR"

# 检查数据库文件是否存在
if [ ! -f "$DB_PATH" ]; then
  echo "错误: 数据库文件不存在: $DB_PATH"
  exit 1
fi

echo "正在备份 AllFi 数据..."

# 使用 SQLite .backup 命令确保一致性
TEMP_DB="${BACKUP_DIR}/allfi_temp_${TIMESTAMP}.db"
sqlite3 "$DB_PATH" ".backup '${TEMP_DB}'"

# 打包临时数据库
tar -czf "$BACKUP_FILE" -C "$BACKUP_DIR" "$(basename "$TEMP_DB")"

# 清理临时文件
rm -f "$TEMP_DB"

# 可选：使用 openssl 加密
if [ "${ENCRYPT_BACKUP:-false}" = "true" ]; then
  openssl enc -aes-256-cbc -pbkdf2 -salt \
    -in "$BACKUP_FILE" -out "${BACKUP_FILE}.enc"
  rm "$BACKUP_FILE"
  BACKUP_FILE="${BACKUP_FILE}.enc"
  echo "已加密备份"
fi

# 保留最近 30 个备份
ls -t "${BACKUP_DIR}"/allfi_backup_* 2>/dev/null | tail -n +31 | xargs rm -f 2>/dev/null || true

echo "备份完成: $BACKUP_FILE"
echo "大小: $(du -h "$BACKUP_FILE" | cut -f1)"
