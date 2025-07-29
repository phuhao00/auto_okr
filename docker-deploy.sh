#!/bin/bash

# Git报告生成器 - Docker一键部署脚本

set -e

echo "🚀 开始部署Git报告生成器..."

# 检查Docker和Docker Compose是否安装
if ! command -v docker &> /dev/null; then
    echo "❌ 错误: Docker未安装，请先安装Docker"
    exit 1
fi

if ! command -v docker-compose &> /dev/null && ! docker compose version &> /dev/null; then
    echo "❌ 错误: Docker Compose未安装，请先安装Docker Compose"
    exit 1
fi

# 停止并删除现有容器（如果存在）
echo "🛑 停止现有容器..."
docker-compose down 2>/dev/null || docker compose down 2>/dev/null || true

# 构建并启动服务
echo "🔨 构建Docker镜像..."
if command -v docker-compose &> /dev/null; then
    docker-compose build --no-cache
else
    docker compose build --no-cache
fi

echo "🚀 启动服务..."
if command -v docker-compose &> /dev/null; then
    docker-compose up -d
else
    docker compose up -d
fi

# 等待服务启动
echo "⏳ 等待服务启动..."
sleep 10

# 检查服务状态
echo "🔍 检查服务状态..."
if command -v docker-compose &> /dev/null; then
    docker-compose ps
else
    docker compose ps
fi

# 显示访问信息
echo ""
echo "✅ 部署完成！"
echo "📱 前端界面: http://localhost:3000"
echo "🔧 后端API: http://localhost:8080"
echo "📊 健康检查: http://localhost:8080/api/health"
echo ""
echo "💡 使用说明:"
echo "   - 在前端界面中输入Git仓库路径"
echo "   - 选择报告类型和日期范围"
echo "   - 点击生成报告按钮"
echo ""
echo "🛠️ 管理命令:"
echo "   停止服务: docker-compose down 或 docker compose down"
echo "   查看日志: docker-compose logs -f 或 docker compose logs -f"
echo "   重启服务: docker-compose restart 或 docker compose restart"