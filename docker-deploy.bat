@echo off
chcp 65001 >nul
setlocal enabledelayedexpansion

echo 🚀 开始部署Git报告生成器...
echo.

REM 检查Docker是否安装
docker --version >nul 2>&1
if errorlevel 1 (
    echo ❌ 错误: Docker未安装，请先安装Docker Desktop
    echo 下载地址: https://www.docker.com/products/docker-desktop
    pause
    exit /b 1
)

REM 检查Docker Compose是否可用
docker compose version >nul 2>&1
if errorlevel 1 (
    docker-compose --version >nul 2>&1
    if errorlevel 1 (
        echo ❌ 错误: Docker Compose未安装或不可用
        pause
        exit /b 1
    )
    set COMPOSE_CMD=docker-compose
) else (
    set COMPOSE_CMD=docker compose
)

echo ✅ Docker环境检查通过
echo.

REM 停止并删除现有容器
echo 🛑 停止现有容器...
%COMPOSE_CMD% down >nul 2>&1

REM 构建镜像
echo 🔨 构建Docker镜像...
%COMPOSE_CMD% build --no-cache
if errorlevel 1 (
    echo ❌ 构建失败，请检查错误信息
    pause
    exit /b 1
)

REM 启动服务
echo 🚀 启动服务...
%COMPOSE_CMD% up -d
if errorlevel 1 (
    echo ❌ 启动失败，请检查错误信息
    pause
    exit /b 1
)

REM 等待服务启动
echo ⏳ 等待服务启动...
timeout /t 10 /nobreak >nul

REM 检查服务状态
echo 🔍 检查服务状态...
%COMPOSE_CMD% ps

echo.
echo ✅ 部署完成！
echo.
echo 📱 前端界面: http://localhost:3000
echo 🔧 后端API: http://localhost:8080
echo 📊 健康检查: http://localhost:8080/api/health
echo.
echo 💡 使用说明:
echo    - 在前端界面中输入Git仓库路径
echo    - 选择报告类型和日期范围
echo    - 点击生成报告按钮
echo.
echo 🛠️ 管理命令:
echo    停止服务: %COMPOSE_CMD% down
echo    查看日志: %COMPOSE_CMD% logs -f
echo    重启服务: %COMPOSE_CMD% restart
echo.
echo 按任意键打开前端界面...
pause >nul
start http://localhost:3000