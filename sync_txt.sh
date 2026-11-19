#!/bin/sh
# 把 txt/ 下的界面文本同步到 assets/txt/，供 go:embed 内嵌。
# 修改菜单文案后运行本脚本，再 go build 即可。
set -e
cd "$(dirname "$0")"
mkdir -p assets/txt/steganography
cp txt/*.txt assets/txt/
cp txt/steganography/*.txt assets/txt/steganography/
echo "synced txt/ -> assets/txt/"
