# Git 提交记录报告生成器 Makefile

# 变量定义
APP_NAME = git-report
VERSION = 1.0.0
BUILD_TIME = $(shell date +"%Y-%m-%d %H:%M:%S")
GIT_COMMIT = $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# Go 相关变量
GOCMD = go
GOBUILD = $(GOCMD) build
GOCLEAN = $(GOCMD) clean
GOTEST = $(GOCMD) test
GOGET = $(GOCMD) get
GOMOD = $(GOCMD) mod

# 构建标志
LDFLAGS = -ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -X main.GitCommit=$(GIT_COMMIT)"

# 默认目标
.PHONY: all
all: clean build

# 构建
.PHONY: build
build:
	@echo "构建 $(APP_NAME)..."
	$(GOBUILD) $(LDFLAGS) -o $(APP_NAME).exe .
	@echo "构建完成: $(APP_NAME).exe"

# 构建所有平台
.PHONY: build-all
build-all: clean
	@echo "构建所有平台版本..."
	@mkdir -p dist
	# Windows
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o dist/$(APP_NAME)-windows-amd64.exe .
	GOOS=windows GOARCH=386 $(GOBUILD) $(LDFLAGS) -o dist/$(APP_NAME)-windows-386.exe .
	# Linux
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o dist/$(APP_NAME)-linux-amd64 .
	GOOS=linux GOARCH=386 $(GOBUILD) $(LDFLAGS) -o dist/$(APP_NAME)-linux-386 .
	# macOS
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o dist/$(APP_NAME)-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o dist/$(APP_NAME)-darwin-arm64 .
	@echo "所有平台构建完成，文件位于 dist/ 目录"

# 清理
.PHONY: clean
clean:
	@echo "清理构建文件..."
	$(GOCLEAN)
	@if exist $(APP_NAME).exe del $(APP_NAME).exe
	@if exist dist rmdir /s /q dist
	@echo "清理完成"

# 测试
.PHONY: test
test:
	@echo "运行测试..."
	$(GOTEST) -v ./...

# 格式化代码
.PHONY: fmt
fmt:
	@echo "格式化代码..."
	$(GOCMD) fmt ./...

# 代码检查
.PHONY: vet
vet:
	@echo "代码检查..."
	$(GOCMD) vet ./...

# 更新依赖
.PHONY: deps
deps:
	@echo "更新依赖..."
	$(GOMOD) tidy
	$(GOMOD) download

# 安装到系统
.PHONY: install
install: build
	@echo "安装到系统..."
	@if not exist "%USERPROFILE%\bin" mkdir "%USERPROFILE%\bin"
	copy $(APP_NAME).exe "%USERPROFILE%\bin\"
	@echo "已安装到 %USERPROFILE%\bin\$(APP_NAME).exe"
	@echo "请确保 %USERPROFILE%\bin 在您的 PATH 环境变量中"

# 卸载
.PHONY: uninstall
uninstall:
	@echo "从系统卸载..."
	@if exist "%USERPROFILE%\bin\$(APP_NAME).exe" del "%USERPROFILE%\bin\$(APP_NAME).exe"
	@echo "卸载完成"

# 运行示例
.PHONY: demo
demo: build
	@echo "运行演示..."
	@echo "生成今日日报:"
	.\$(APP_NAME).exe
	@echo ""
	@echo "生成本周周报:"
	.\$(APP_NAME).exe -type weekly

# 创建发布包
.PHONY: release
release: build-all
	@echo "创建发布包..."
	@mkdir -p release
	copy README.md release\
	copy config.example.json release\
	xcopy templates release\templates\ /E /I
	copy dist\* release\
	@echo "发布包已创建在 release/ 目录"

# 显示帮助
.PHONY: help
help:
	@echo "可用的 make 命令:"
	@echo "  build      - 构建应用程序"
	@echo "  build-all  - 构建所有平台版本"
	@echo "  clean      - 清理构建文件"
	@echo "  test       - 运行测试"
	@echo "  fmt        - 格式化代码"
	@echo "  vet        - 代码检查"
	@echo "  deps       - 更新依赖"
	@echo "  install    - 安装到系统"
	@echo "  uninstall  - 从系统卸载"
	@echo "  demo       - 运行演示"
	@echo "  release    - 创建发布包"
	@echo "  help       - 显示此帮助信息"