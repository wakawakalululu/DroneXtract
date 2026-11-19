#!/bin/sh
# Linux / macOS 一键安装脚本。
#
# 做的事：
#   1. 检测 ExifTool，缺失时按发行版用包管理器安装（自动处理 sudo）
#   2. 在程序目录生成 .env 配置（已存在则保留）
#   3. 把二进制安装到 ~/.local/bin 并提示 PATH
#
# 用法：
#   ./install.sh              安装到 ~/.local/bin
#   ./install.sh --prefix /opt/dronextract
#   ./install.sh --no-deps    跳过依赖安装
set -eu
cd "$(dirname "$0")"

PREFIX="${HOME}/.local"
SKIP_DEPS=0

while [ $# -gt 0 ]; do
  case "$1" in
    --prefix) PREFIX="$2"; shift 2 ;;
    --no-deps) SKIP_DEPS=1; shift ;;
    -h|--help)
      sed -n '2,8p' "$0" | sed 's/^#\{0,1\} \{0,1\}//'
      exit 0 ;;
    *) echo "未知参数: $1" >&2; exit 1 ;;
  esac
done

say() { printf '\033[36m%s\033[0m\n' "$*"; }
ok()  { printf '\033[32m%s\033[0m\n' "$*"; }
warn(){ printf '\033[33m%s\033[0m\n' "$*"; }
err() { printf '\033[31m%s\033[0m\n' "$*" >&2; }

# ---------- 1. 依赖 ----------
install_exiftool() {
  if command -v exiftool >/dev/null 2>&1; then
    ok "ExifTool 已安装：$(exiftool -ver 2>/dev/null || echo '?')"
    return 0
  fi

  warn "未检测到 ExifTool，尝试安装..."
  SUDO=""
  if [ "$(id -u)" -ne 0 ]; then
    if command -v sudo >/dev/null 2>&1; then
      SUDO="sudo"
    else
      err "需要 root 权限且未安装 sudo，请手动执行："
      err "  apt-get install -y libimage-exiftool-perl   # Debian/Ubuntu"
      err "  dnf install -y perl-Image-ExifTool         # Fedora/RHEL"
      return 1
    fi
  fi

  if command -v apt-get >/dev/null 2>&1; then
    $SUDO apt-get update -qq || true
    $SUDO env DEBIAN_FRONTEND=noninteractive apt-get install -y libimage-exiftool-perl
  elif command -v dnf >/dev/null 2>&1; then
    $SUDO dnf install -y perl-Image-ExifTool
  elif command -v yum >/dev/null 2>&1; then
    $SUDO yum install -y perl-Image-ExifTool
  elif command -v pacman >/dev/null 2>&1; then
    $SUDO pacman -S --noconfirm perl-image-exiftool
  elif command -v zypper >/dev/null 2>&1; then
    $SUDO zypper install -y exiftool
  elif command -v apk >/dev/null 2>&1; then
    $SUDO apk add exiftool
  elif command -v brew >/dev/null 2>&1; then
    brew install exiftool
  else
    err "未找到可用的包管理器，请参考 https://exiftool.org/install.html"
    return 1
  fi

  hash -r 2>/dev/null || true
  if command -v exiftool >/dev/null 2>&1; then
    ok "ExifTool 安装成功"
  else
    warn "安装完成但 exiftool 不在 PATH，请重开终端或执行 hash -r"
    return 1
  fi
}

# ---------- 2. 配置 ----------
write_config() {
  # 配置放在程序所在目录，便于整目录搬迁
  cfg_dir="$PREFIX/bin"
  cfg="$cfg_dir/.env"
  if [ -f "$cfg" ]; then
    ok "配置已存在，保留：$cfg"
    return 0
  fi
  mkdir -p "$cfg_dir"
  cat > "$cfg" <<'ENVEOF'
TELEMETRY_VIS_DOWNSAMPLE = 40
FLIGHT_MAP_DOWNSAMPLE = 10
ANALYSIS_DOWNSAMPLE = 40
ANALYSIS_MAX_VARIANCE = "300 90 100 1000 120 100 30 30 360 40 25 1684 1684 1684 1684 1684 1684 1684 1684 360 180 180 100 20 5 500 50"
ENVEOF
  ok "已生成配置：$cfg"
}

# ---------- 3. 安装二进制 ----------
install_binary() {
  arch=$(uname -m)
  case "$arch" in
    x86_64|amd64) arch=amd64 ;;
    aarch64|arm64) arch=arm64 ;;
    *) warn "未识别的架构 $arch，尝试 amd64"; arch=amd64 ;;
  esac

  os=$(uname -s | tr '[:upper:]' '[:lower:]')
  case "$os" in
    linux|darwin) ;;
    *) err "本脚本仅支持 Linux / macOS"; exit 1 ;;
  esac

  src="dist/DroneXtract-${os}-${arch}"
  if [ ! -f "$src" ]; then
    say "未找到预编译产物，正在本地构建 $os/$arch ..."
    command -v go >/dev/null 2>&1 || { err "需要 Go 才能从源码构建"; exit 1; }
    mkdir -p dist
    CGO_ENABLED=0 GOOS="$os" GOARCH="$arch" go build -trimpath -ldflags "-s -w" -o "$src" .
  fi

  mkdir -p "$PREFIX/bin"
  install -m 0755 "$src" "$PREFIX/bin/dronextract"
  ok "已安装：$PREFIX/bin/dronextract"

  case ":$PATH:" in
    *":$PREFIX/bin:"*) ;;
    *) warn "请把 $PREFIX/bin 加入 PATH，例如："
       warn "  echo 'export PATH=\"$PREFIX/bin:\$PATH\"' >> ~/.profile && . ~/.profile" ;;
  esac
}

say "DroneXtract 安装程序"
[ "$SKIP_DEPS" -eq 0 ] && install_exiftool || true
install_binary   # 先建目录并放置二进制
write_config     # 配置写在程序所在目录

echo
say "安装完成，运行："
echo "  dronextract            # 启动界面"
echo "  dronextract --check    # 复查依赖"
