@echo off
echo ====================================
echo Git 提交记录报告生成器 演示
echo ====================================
echo.

echo 1. 显示帮助信息:
echo.
git-report.exe -h
echo.

echo 2. 生成今日日报 (如果当前目录是Git仓库):
echo.
echo git-report.exe
echo.

echo 3. 生成本周周报:
echo.
echo git-report.exe -type weekly
echo.

echo 4. 生成指定日期的日报:
echo.
echo git-report.exe -date 2024-01-15
echo.

echo 5. 指定作者生成报告:
echo.
echo git-report.exe -author "张三"
echo.

echo 6. 保存报告到文件:
echo.
echo git-report.exe -output daily-report.md
echo.

echo 7. 使用自定义模板:
echo.
echo git-report.exe -template templates/custom-template.tmpl
echo.

echo 8. 指定Git仓库路径:
echo.
echo git-report.exe -repo /path/to/your/repo
echo.

echo ====================================
echo 注意事项:
echo - 确保在Git仓库目录中运行
echo - 确保已配置Git用户信息
echo - 确保Git命令可用
echo ====================================
pause