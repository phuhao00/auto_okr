@echo off
chcp 65001 >nul
echo ========================================
echo    Git 报告生成器 - Docker 一键部署
echo    (使用国内镜像源解决网络问题)
echo ========================================
echo.

REM 检查 Docker 是否安装
docker --version >nul 2>&1
if %errorlevel% neq 0 (
    echo [错误] Docker 未安装或未启动
    echo 请先安装 Docker Desktop: https://www.docker.com/products/docker-desktop
    pause
    exit /b 1
)

REM 检查 Docker Compose 是否可用
docker compose version >nul 2>&1
if %errorlevel% neq 0 (
    echo [错误] Docker Compose 不可用
    echo 请确保 Docker Desktop 正常运行
    pause
    exit /b 1
)

echo [信息] Docker 环境检查通过
echo.

REM 停止并删除现有容器
echo [步骤 1/5] 清理现有容器...
docker compose -f docker-compose.china.yml down --remove-orphans
echo.

REM 构建 Docker 镜像
echo [步骤 2/5] 构建 Docker 镜像 (使用国内镜像源)...
docker compose -f docker-compose.china.yml build --no-cache
if %errorlevel% neq 0 (
    echo [错误] Docker 镜像构建失败
    echo 请检查网络连接或查看 docker-troubleshoot.md 获取解决方案
    pause
    exit /b 1
)
echo.

REM 启动服务
echo [步骤 3/5] 启动服务...
docker compose -f docker-compose.china.yml up -d
if %errorlevel% neq 0 (
    echo [错误] 服务启动失败
    pause
    exit /b 1
)
echo.

REM 等待服务启动
echo [步骤 4/5] 等待服务启动...
timeout /t 10 /nobreak >nul
echo.

REM 检查服务状态
echo [步骤 5/5] 检查服务状态...
docker compose -f docker-compose.china.yml ps
echo.

echo ========================================
echo           部署完成！
echo ========================================
echo.
echo 🌐 前端界面: http://localhost:3000
echo 🔧 后端 API: http://localhost:8080
echo 💚 健康检查: http://localhost:8080/health
echo.
echo ========================================
echo           使用说明
echo ========================================
echo.
echo 1. 在浏览器中打开 http://localhost:3000
echo 2. 选择要分析的 Git 仓库目录
echo 3. 设置报告参数（日期范围、作者等）
echo 4. 生成并查看报告
echo.
echo ========================================
echo           管理命令
echo ========================================
echo.
echo 停止服务: docker compose -f docker-compose.china.yml stop
echo 重启服务: docker compose -f docker-compose.china.yml restart
echo 查看日志: docker compose -f docker-compose.china.yml logs
echo 删除服务: docker compose -f docker-compose.china.yml down
echo.
echo 如遇到问题，请查看 docker-troubleshoot.md 文件
echo.
pause