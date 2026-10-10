#!/usr/bin/env bash
set -e

BUILD_VERSION=${1:-"v1.0.0"}
GIT_COMMIT=${2:-$(git rev-parse --short HEAD 2>/dev/null || echo "manual")}
BUILD_TIME=$(date '+%Y-%m-%d %H:%M:%S')

echo "=== 使用 Docker Compose 自动化部署 ==="
echo "Version:    $BUILD_VERSION"
echo "Git Commit: $GIT_COMMIT"
echo "Build Time: $BUILD_TIME"

BUILD_VERSION="$BUILD_VERSION" \
GIT_COMMIT="$GIT_COMMIT" \
BUILD_TIME="$BUILD_TIME" \
docker compose up -d --build --remove-orphans

echo "=== 部署完成，服务已在后台运行 ==="
