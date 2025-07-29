@echo off
chcp 65001 >nul
echo ========================================
echo    配置 Docker 镜像加速器
echo ========================================
echo.

echo [信息] 正在配置 Docker 镜像加速器...
echo.

REM 创建 Docker 配置目录
if not exist "%USERPROFILE%\.docker" mkdir "%USERPROFILE%\.docker"

REM 创建 daemon.json 配置文件
echo {
  "registry-mirrors": [
    "https://docker.mirrors.ustc.edu.cn",
    "https://hub-mirror.c.163.com",
    "https://mirror.baidubce.com"
  ],
  "insecure-registries": [],
  "debug": false,
  "experimental": false
} > "%USERPROFILE%\.docker\daemon.json"

echo [完成] Docker 镜像加速器配置已创建
echo 配置文件位置: %USERPROFILE%\.docker\daemon.json
echo.
echo [重要] 请重启 Docker Desktop 以使配置生效
echo.
echo 配置的镜像源:
echo - 中科大镜像: https://docker.mirrors.ustc.edu.cn
echo - 网易镜像: https://hub-mirror.c.163.com
echo - 百度镜像: https://mirror.baidubce.com
echo.
echo 重启 Docker Desktop 后，请运行以下命令验证:
echo docker info ^| findstr "Registry Mirrors"
echo.
pause