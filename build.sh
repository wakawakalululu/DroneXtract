#!/bin/sh
# 跨平台构建脚本：一次产出 Linux / macOS / Windows 三套二进制。
#
# 用法：
#   ./build.sh                构建当前平台
#   ./build.sh all            构建 linux/amd64、linux/arm64、darwin/arm64、windows/amd64
#   ./build.sh linux arm64    构建指定平台
#
# 产物放在 dist/ 下，Linux/macOS 二进制会自动带上可执行权限。
set -eu
cd "$(dirname "$0")"

mkdir -p dist

build() {
  os=$1
  arch=$2
  out="dist/DroneXtract-${os}-${arch}"
  [ "$os" = "windows" ] && out="${out}.exe"

  echo "==> ${os}/${arch} -> ${out}"
  CGO_ENABLED=0 GOOS="$os" GOARCH="$arch" go build -trimpath -ldflags "-s -w" -o "$out" .

  # Windows 上 chmod 无意义，仅对类 Unix 平台补执行位
  if [ "$os" != "windows" ]; then
    chmod +x "$out"
  fi
}

case "${1:-}" in
  all)
    build linux amd64
    build linux arm64
    build darwin amd64
    build darwin arm64
    build windows amd64
    ;;
  linux|darwin|windows)
    build "$1" "${2:-amd64}"
    ;;
  "")
    goos=$(go env GOOS)
    goarch=$(go env GOARCH)
    build "$goos" "$goarch"
    ;;
  *)
    echo "用法: $0 [all|linux|darwin|windows] [amd64|arm64]" >&2
    exit 1
    ;;
esac

echo
echo "完成，产物位于 dist/"
ls -lh dist/ 2>/dev/null | tail -n +2
