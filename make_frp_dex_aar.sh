#!/usr/bin/env bash
set -e

##
## 将 gomobile bind 生成的 frp_raw.aar 转成“可动态加载”的 AAR（包含 classes.dex）
## 使用方式：
##   1. 在目标目录执行 gomobile bind：
##        gomobile bind -target=android -javapkg com.android.http -o ./frp_raw.aar github.com/fatedier/frp/cmd/frp/frps github.com/fatedier/frp/cmd/frp/frpc
##   2. 然后在同一目录运行：
##        ./make_frp_dex_aar.sh              # 当前目录有 frp_raw.aar
##      或：
##        ./make_frp_dex_aar.sh ./release    # 指定目录下有 frp_raw.aar
##

# ===== 根据你本机环境修改 =====
# 优先使用环境变量 ANDROID_SDK_ROOT / ANDROID_HOME
SDK_ROOT="${ANDROID_SDK_ROOT:-$ANDROID_HOME}"

# 如果环境变量没设，尝试使用 Windows 默认安装路径（根据你当前用户名 admin 写死）
if [ -z "$SDK_ROOT" ]; then
  if [ -d "/c/Users/admin/AppData/Local/Android/Sdk" ]; then
    SDK_ROOT="/c/Users/admin/AppData/Local/Android/Sdk"
  fi
fi

BUILD_TOOLS_VERSION="34.0.0"   # 改成你本机安装的 build-tools 版本，如 33.0.2

# 参数：
#   $1: frp_raw.aar 所在目录（默认当前目录）
#   $2: 原始 AAR 文件名（默认 frp_raw.aar）
#   $3: 输出 AAR 文件名（默认 frp.aar）
INPUT_DIR="${1:-.}"
INPUT_DIR_ABS="$(cd "$INPUT_DIR" && pwd)"
INPUT_AAR="${2:-frp_raw.aar}"   # 目录下的输入 AAR 文件名
OUTPUT_AAR="${3:-frp.aar}"      # 目录下的输出 AAR 文件名

if [ -z "$SDK_ROOT" ]; then
  echo "请先设置 ANDROID_SDK_ROOT 或 ANDROID_HOME 环境变量，或在 make_frp_dex_aar.sh 中修改 SDK_ROOT 默认路径"
  exit 1
fi

if [ ! -f "$INPUT_DIR_ABS/$INPUT_AAR" ]; then
  echo "输入 AAR 不存在: $INPUT_DIR_ABS/$INPUT_AAR"
  exit 1
fi

D8="$SDK_ROOT/build-tools/$BUILD_TOOLS_VERSION/d8"

# 检查是否安装 zip（用于重新打包 AAR）
if ! command -v zip >/dev/null 2>&1; then
  echo "未找到 zip 命令，请先安装 zip（例如："
  echo "  Linux: sudo apt-get install zip unzip  或  sudo yum install zip unzip"
  echo "  macOS: brew install zip unzip"
  echo "  Windows Git Bash（有 Chocolatey）: choco install zip unzip -y"
  echo "）"
  exit 1
fi

# 在 Linux/macOS 下通常是无扩展名的 d8，可执行；
# 在 Windows + Git Bash 下通常是 d8.bat。
if [ ! -x "$D8" ]; then
  if [ -f "${D8}.bat" ]; then
    D8="${D8}.bat"
  else
    echo "找不到 d8 可执行文件: $D8 或 ${D8}.bat"
    exit 1
  fi
fi

WORK_DIR="$INPUT_DIR_ABS/tmp_frp_aar"
rm -rf "$WORK_DIR"
mkdir -p "$WORK_DIR"

echo "[0/4] 输入目录: $INPUT_DIR_ABS"
echo "[1/4] 解压 AAR: $INPUT_DIR_ABS/$INPUT_AAR -> $WORK_DIR"
unzip -q "$INPUT_DIR_ABS/$INPUT_AAR" -d "$WORK_DIR"

if [ ! -f "$WORK_DIR/classes.jar" ]; then
  echo "AAR 中没有 classes.jar，结构不对？"
  exit 1
fi

echo "[2/4] 使用 d8 将 classes.jar 转成 classes.dex ..."
"$D8" \
  --release \
  --min-api 19 \
  --output "$WORK_DIR" \
  "$WORK_DIR/classes.jar"

if [ ! -f "$WORK_DIR/classes.dex" ]; then
  echo "d8 没有生成 classes.dex，请检查 d8 输出"
  exit 1
fi

# 如需避免混淆，可以删除 classes.jar，只保留 classes.dex 和 jni/ 等
# rm "$WORK_DIR/classes.jar"

echo "[3/4] 重新打包 AAR（包含 classes.dex 和原来的 jni/、res/ 等）..."
cd "$WORK_DIR"
zip -qr "../$OUTPUT_AAR" ./*
cd - >/dev/null

rm -rf "$WORK_DIR"

echo "[4/4] 完成！输出文件: $INPUT_DIR_ABS/$OUTPUT_AAR"

# 清理中间产物 frp_raw.aar
rm -f "$INPUT_DIR_ABS/$INPUT_AAR"