#!/usr/bin/env bash
set -e

echo "=== [1/2] 构建 Go 后端 ==="
cd backend
go test -v ./...
CGO_ENABLED=0 go build -ldflags="-s -w -X 'main.Version=local-test'" -o server .
cd ..

echo "=== [2/2] 构建 Vue 前端 ==="
cd frontend
npm install --registry=https://registry.npmmirror.com
npm run build
cd ..

echo "=== 本地直接构建完成 ==="
