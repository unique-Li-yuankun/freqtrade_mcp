#!/bin/bash

# Mac 多架构构建脚本
# 构建 Intel (x86_64) 和 Apple Silicon (arm64) 版本

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 项目信息
PROJECT_NAME="freqtrade_mcp"
VERSION="v1.0.0"
BUILD_DIR="build"
DIST_DIR="dist"

echo -e "${BLUE}🚀 开始构建 ${PROJECT_NAME} ${VERSION}${NC}"
echo "=================================="

# 创建构建和分发目录
echo -e "${YELLOW}📁 创建构建目录...${NC}"
mkdir -p $BUILD_DIR
mkdir -p $DIST_DIR

# 清理之前的构建
echo -e "${YELLOW}🧹 清理之前的构建文件...${NC}"
rm -rf $BUILD_DIR/*
rm -rf $DIST_DIR/*

# 构建 Intel (x86_64) 版本
echo -e "${GREEN}🔨 构建 Intel (x86_64) 版本...${NC}"
GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build \
    -ldflags="-s -w" \
    -o "$BUILD_DIR/${PROJECT_NAME}_darwin_amd64" \
    main.go

# 构建 Apple Silicon (arm64) 版本
echo -e "${GREEN}🔨 构建 Apple Silicon (arm64) 版本...${NC}"
GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build \
    -ldflags="-s -w" \
    -o "$BUILD_DIR/${PROJECT_NAME}_darwin_arm64" \
    main.go

# 创建通用二进制文件 (Universal Binary)
echo -e "${GREEN}🔗 创建通用二进制文件...${NC}"
lipo -create \
    "$BUILD_DIR/${PROJECT_NAME}_darwin_amd64" \
    "$BUILD_DIR/${PROJECT_NAME}_darwin_arm64" \
    -output "$BUILD_DIR/${PROJECT_NAME}_darwin_universal"

# 验证架构
echo -e "${YELLOW}🔍 验证二进制文件架构...${NC}"
echo "Intel (x86_64) 版本:"
file "$BUILD_DIR/${PROJECT_NAME}_darwin_amd64"
echo "Apple Silicon (arm64) 版本:"
file "$BUILD_DIR/${PROJECT_NAME}_darwin_arm64"
echo "通用版本:"
file "$BUILD_DIR/${PROJECT_NAME}_darwin_universal"

# 获取文件大小
echo -e "${YELLOW}📊 文件大小信息:${NC}"
echo "Intel (x86_64): $(du -h "$BUILD_DIR/${PROJECT_NAME}_darwin_amd64" | cut -f1)"
echo "Apple Silicon (arm64): $(du -h "$BUILD_DIR/${PROJECT_NAME}_darwin_arm64" | cut -f1)"
echo "通用版本: $(du -h "$BUILD_DIR/${PROJECT_NAME}_darwin_universal" | cut -f1)"

# 创建分发包
echo -e "${YELLOW}📦 创建分发包...${NC}"

# 创建 Intel 版本分发包
INTEL_PACKAGE="${DIST_DIR}/${PROJECT_NAME}_${VERSION}_darwin_amd64.tar.gz"
tar -czf "$INTEL_PACKAGE" \
    -C "$BUILD_DIR" \
    "${PROJECT_NAME}_darwin_amd64"

# 创建 Apple Silicon 版本分发包
ARM_PACKAGE="${DIST_DIR}/${PROJECT_NAME}_${VERSION}_darwin_arm64.tar.gz"
tar -czf "$ARM_PACKAGE" \
    -C "$BUILD_DIR" \
    "${PROJECT_NAME}_darwin_arm64"

# 创建通用版本分发包
UNIVERSAL_PACKAGE="${DIST_DIR}/${PROJECT_NAME}_${VERSION}_darwin_universal.tar.gz"
tar -czf "$UNIVERSAL_PACKAGE" \
    -C "$BUILD_DIR" \
    "${PROJECT_NAME}_darwin_universal"

# 创建完整分发包
FULL_PACKAGE="${DIST_DIR}/${PROJECT_NAME}_${VERSION}_darwin_all.tar.gz"
tar -czf "$FULL_PACKAGE" \
    -C "$BUILD_DIR" \
    "${PROJECT_NAME}_darwin_amd64" \
    "${PROJECT_NAME}_darwin_arm64" \
    "${PROJECT_NAME}_darwin_universal"

# 创建安装脚本
echo -e "${YELLOW}📝 创建安装脚本...${NC}"
cat > "$DIST_DIR/install.sh" << 'EOF'
#!/bin/bash

# Freqtrade MCP 安装脚本

set -e

PROJECT_NAME="freqtrade_mcp"
INSTALL_DIR="/usr/local/bin"

echo "🚀 安装 $PROJECT_NAME..."

# 检测系统架构
ARCH=$(uname -m)
echo "检测到系统架构: $ARCH"

# 选择对应的二进制文件
if [ "$ARCH" = "arm64" ]; then
    BINARY="${PROJECT_NAME}_darwin_arm64"
    echo "使用 Apple Silicon (arm64) 版本"
elif [ "$ARCH" = "x86_64" ]; then
    BINARY="${PROJECT_NAME}_darwin_amd64"
    echo "使用 Intel (x86_64) 版本"
else
    echo "❌ 不支持的架构: $ARCH"
    exit 1
fi

# 检查二进制文件是否存在
if [ ! -f "$BINARY" ]; then
    echo "❌ 找不到二进制文件: $BINARY"
    exit 1
fi

# 安装到系统目录
echo "📁 安装到: $INSTALL_DIR"
sudo cp "$BINARY" "$INSTALL_DIR/$PROJECT_NAME"
sudo chmod +x "$INSTALL_DIR/$PROJECT_NAME"

echo "✅ 安装完成！"
echo "使用方法: $PROJECT_NAME -dir /path/to/freqtrade -address localhost:8000"
EOF

chmod +x "$DIST_DIR/install.sh"

# 创建 README
echo -e "${YELLOW}📖 创建 README...${NC}"
cat > "$DIST_DIR/README.md" << EOF
# Freqtrade MCP ${VERSION}

## 概述

这是一个为 Cursor 提供 Freqtrade 功能的 MCP (Model Context Protocol) 服务器。

## 可用版本

- **Intel (x86_64)**: \`${PROJECT_NAME}_darwin_amd64\` - 适用于 Intel Mac
- **Apple Silicon (arm64)**: \`${PROJECT_NAME}_darwin_arm64\` - 适用于 M1/M2/M3/M4 Mac
- **通用版本**: \`${PROJECT_NAME}_darwin_universal\` - 同时支持 Intel 和 Apple Silicon

## 快速安装

### 自动安装
\`\`\`bash
chmod +x install.sh
./install.sh
\`\`\`

### 手动安装
1. 根据您的 Mac 架构选择对应的二进制文件
2. 将文件复制到 \`/usr/local/bin\` 目录
3. 添加执行权限

\`\`\`bash
# 对于 Intel Mac
sudo cp ${PROJECT_NAME}_darwin_amd64 /usr/local/bin/${PROJECT_NAME}
sudo chmod +x /usr/local/bin/${PROJECT_NAME}

# 对于 Apple Silicon Mac
sudo cp ${PROJECT_NAME}_darwin_arm64 /usr/local/bin/${PROJECT_NAME}
sudo chmod +x /usr/local/bin/${PROJECT_NAME}
\`\`\`

## 使用方法

\`\`\`bash
${PROJECT_NAME} -dir /path/to/your/freqtrade -address localhost:8000
\`\`\`

## 参数说明

- \`-dir\`: Freqtrade 目录路径（必需）
- \`-address\`: 监听地址（默认: localhost:8000）

## 功能特性

- **回测策略**: 执行策略回测
- **下载数据**: 从交易所下载历史数据

## 系统要求

- macOS 10.15 或更高版本
- 已安装 Freqtrade

## 许可证

请参考项目主目录的许可证文件。
EOF

# 显示构建结果
echo -e "${GREEN}✅ 构建完成！${NC}"
echo "=================================="
echo -e "${BLUE}📦 分发包位置:${NC}"
echo "Intel 版本: $INTEL_PACKAGE"
echo "Apple Silicon 版本: $ARM_PACKAGE"
echo "通用版本: $UNIVERSAL_PACKAGE"
echo "完整包: $FULL_PACKAGE"
echo "安装脚本: $DIST_DIR/install.sh"
echo "说明文档: $DIST_DIR/README.md"

echo -e "${BLUE}📁 构建目录内容:${NC}"
ls -la "$BUILD_DIR"

echo -e "${BLUE}📁 分发目录内容:${NC}"
ls -la "$DIST_DIR"

echo -e "${GREEN}🎉 所有文件已准备就绪！${NC}"
