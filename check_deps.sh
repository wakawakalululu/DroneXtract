#!/bin/sh
# 依赖自检：缺少必需依赖时退出码为 2，便于 CI / 部署脚本判定。
set -e
cd "$(dirname "$0")"

if [ ! -x ./DroneXtract.exe ] && [ ! -x ./DroneXtract ]; then
  echo "未找到可执行文件，请先构建： go build -o DroneXtract ." >&2
  exit 1
fi

BIN=./DroneXtract
[ -x ./DroneXtract.exe ] && BIN=./DroneXtract.exe

# 行尾自检：CRLF 会让 Linux 报 bad interpreter
if grep -qU $'' "$0" 2>/dev/null; then
  echo "警告：脚本包含 CR 字符，请执行 sed -i 's/$//' $0" >&2
fi

exec "$BIN" --check
