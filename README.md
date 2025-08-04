# Git 提交记录报告生成器

一个基于 Git 提交记录自动生成日报和周报的 Golang 工具，支持命令行和 Web 界面两种使用方式。

## 功能特性

- 📅 **自动生成日报和周报**：基于 Git 提交记录自动生成工作报告
- 🌐 **Web 界面**：现代化的 Next.js 前端界面，支持在线生成和预览报告
- 🖥️ **命令行工具**：传统的 CLI 工具，适合脚本化和自动化场景
- 🔌 **HTTP API**：RESTful API 接口，支持第三方集成
- 👤 **多用户支持**：可指定特定作者或使用当前 Git 用户
- 📊 **详细统计信息**：包含提交次数、代码行数、文件类型分布等
- 🏷️ **智能分类**：自动将提交按功能开发、Bug修复、重构等类别分组
- 🎨 **自定义模板**：支持使用自定义模板定制报告格式
- 💾 **灵活输出**：支持控制台输出、文件保存或在线预览

## 界面预览

### Web 界面

![Git Report Generator Web Interface](https://github.com/phuhao00/auto_okr/blob/main/image.png)

*现代化的 Web 界面，支持在线生成和预览 Git 提交报告*

![Git Report Generator Additional Interface](https://github.com/phuhao00/auto_okr/blob/main/image02.png)

*Git 报告生成器的润色汇报功能界面*

## 安装

### 🐳 Docker 一键部署（推荐）

**前提条件：**
- 安装 [Docker](https://www.docker.com/products/docker-desktop)
- 安装 Docker Compose

**Windows 用户：**
```bash
# 双击运行或在命令行执行
docker-deploy.bat
```

**Linux/macOS 用户：**
```bash
# 给脚本执行权限
chmod +x docker-deploy.sh
# 运行部署脚本
./docker-deploy.sh
```

**手动 Docker 部署：**
```bash
# 构建并启动服务
docker-compose up -d --build

# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs -f

# 停止服务
docker-compose down
```

### 📦 传统安装方式

#### 后端服务

```bash
# 克隆项目
git clone <repository-url>
cd git-report-generator

# 安装 Go 依赖
go mod tidy

# 编译
go build -o git-report.exe
```

#### 前端界面（可选）

```bash
# 进入前端目录
cd frontend

# 安装依赖
npm install

# 启动开发服务器
npm run dev
```

## 使用方法

### Web 界面使用

1. **启动后端服务**
```bash
# 启动 HTTP 服务器（默认端口 8080）
./git-report.exe -server
```

2. **启动前端界面**
```bash
# 在另一个终端中启动前端（默认端口 3000）
cd frontend
npm run dev
```

3. **访问 Web 界面**
   - 打开浏览器访问 `http://localhost:3000`
   - 输入 Git 仓库路径
   - 选择报告类型和日期
   - 在线生成和预览报告
   - 支持下载 Markdown 格式报告

### API 接口

#### 生成报告
```bash
POST /api/generate-report
Content-Type: application/json

{
  "repo_path": "/path/to/repo",
  "report_type": "daily",
  "date": "2024-01-15",
  "author": "张三"
}
```

#### 健康检查
```bash
GET /api/health
```

### 命令行使用

#### 基本用法

```bash
# 生成今日日报
./git-report.exe

# 生成指定日期的日报
./git-report.exe -date 2024-01-15

# 生成本周周报
./git-report.exe -type weekly

# 生成指定作者的报告
./git-report.exe -author "张三"

# 指定Git仓库路径
./git-report.exe -repo /path/to/your/repo

# 保存报告到文件
./git-report.exe -output daily-report.md

# 使用自定义模板
./git-report.exe -template custom-template.tmpl
```

### 命令行参数

| 参数 | 说明 | 默认值 | 示例 |
|------|------|--------|---------|
| `-server` | 启动 HTTP 服务器模式 | false | `-server` |
| `-type` | 报告类型：daily, weekly | daily | `-type weekly` |
| `-date` | 指定日期 (YYYY-MM-DD) | 今天 | `-date 2024-01-15` |
| `-repo` | Git仓库路径 | 当前目录 | `-repo /path/to/repo` |
| `-author` | 指定作者 | 当前Git用户 | `-author "张三"` |
| `-output` | 输出文件路径 | 控制台输出 | `-output report.md` |
| `-template` | 自定义模板文件 | 内置模板 | `-template my-template.tmpl` |

## 报告内容

### 日报包含
- 📊 当日统计信息（提交次数、修改文件数、代码行数）
- 🔥 主要修改文件列表
- 📁 文件类型分布
- 📝 按类别分组的详细提交记录

### 周报包含
- 📊 本周统计信息
- 📅 每日提交分布
- 🔥 主要修改文件列表
- 📁 文件类型分布
- 📝 按类别分组的工作内容
- 💡 本周工作总结

## 提交分类规则

工具会根据提交信息自动将提交分类为：

- **功能开发**：包含 feat, feature, add, 新增, 功能 等关键词
- **Bug修复**：包含 fix, bug, 修复, 修正 等关键词
- **代码重构**：包含 refactor, 重构 等关键词
- **文档更新**：包含 doc, readme, 文档 等关键词
- **测试相关**：包含 test, 测试 等关键词
- **配置修改**：包含 config, 配置 等关键词
- **其他**：不符合以上分类的提交

## 自定义模板

可以创建自定义模板文件来定制报告格式。模板使用 Go 的 `text/template` 语法。

### 可用变量

- `.Type`：报告类型（daily/weekly）
- `.Date`：报告日期
- `.Period`：时间范围描述
- `.Author`：作者名称
- `.RepoInfo`：仓库信息（name, branch, url）
- `.Commits`：提交记录列表
- `.Summary`：统计摘要
- `.Categories`：按类别分组的提交
- `.GeneratedAt`：生成时间

### 可用函数

- `formatTime`：格式化时间
- `formatDate`：格式化日期
- `formatShortHash`：格式化短哈希
- `join`：连接字符串数组
- `sortedKeys`：获取排序后的键
- `sortedFileTypes`：获取排序后的文件类型

### 示例模板

```go
# {{.Author}} 的工作{{if eq .Type "daily"}}日报{{else}}周报{{end}}

**时间：** {{.Period}}
**仓库：** {{.RepoInfo.name}}

## 统计信息
- 提交次数：{{.Summary.TotalCommits}}
- 修改文件：{{.Summary.TotalFiles}}
- 代码变更：+{{.Summary.TotalAdditions}} -{{.Summary.TotalDeletions}}

{{range sortedKeys .Categories}}
### {{.}}
{{range index $.Categories .}}
- {{.Message}} ({{formatShortHash .Hash}})
{{end}}
{{end}}
```

## 项目架构

```
┌─────────────────────────────────────────────────────────────────┐
│                           用户界面层                              │
├─────────────────────────────────────────────────────────────────┤
│  Web 浏览器 (http://localhost:3000)                             │
│  ┌─────────────────────────────────────────────────────────┐    │
│  │                Next.js 前端                              │    │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐      │    │
│  │  │   页面组件   │  │   样式组件   │  │   API 客户端 │      │    │
│  │  │ (page.tsx)  │  │(Tailwind CSS)│  │  (fetch)    │      │    │
│  │  └─────────────┘  └─────────────┘  └─────────────┘      │    │
│  └─────────────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────────────┘
                                │
                                │ HTTP API 调用
                                │ (POST /api/generate-report)
                                ▼
┌─────────────────────────────────────────────────────────────────┐
│                        后端服务层                                 │
├─────────────────────────────────────────────────────────────────┤
│  Go HTTP 服务器 (http://localhost:8080)                         │
│  ┌─────────────────────────────────────────────────────────┐    │
│  │                   HTTP 路由层                            │    │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐      │    │
│  │  │   路由处理   │  │   CORS 中间件│  │   请求验证   │      │    │
│  │  │(Gorilla Mux)│  │   (rs/cors)  │  │             │      │    │
│  │  └─────────────┘  └─────────────┘  └─────────────┘      │    │
│  └─────────────────────────────────────────────────────────┘    │
│  ┌─────────────────────────────────────────────────────────┐    │
│  │                   业务逻辑层                              │    │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐      │    │
│  │  │  报告生成器  │  │  模板渲染器  │  │   Git 操作   │      │    │
│  │  │ (report.go) │  │(renderer.go)│  │  (git.go)   │      │    │
│  │  └─────────────┘  └─────────────┘  └─────────────┘      │    │
│  └─────────────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────────────┘
                                │
                                │ Git 命令调用
                                │ (git log, git show)
                                ▼
┌─────────────────────────────────────────────────────────────────┐
│                        数据源层                                   │
├─────────────────────────────────────────────────────────────────┤
│  ┌─────────────────────────────────────────────────────────┐    │
│  │                   Git 仓库                              │    │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐      │    │
│  │  │  提交历史    │  │   分支信息   │  │   文件变更   │      │    │
│  │  │ (.git/logs) │  │ (.git/refs) │  │             │      │    │
│  │  └─────────────┘  └─────────────┘  └─────────────┘      │    │
│  └─────────────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────────────┘
```

### 架构说明

1. **用户界面层**：基于 Next.js 的现代化 Web 界面，提供直观的用户交互
2. **后端服务层**：Go 语言实现的 HTTP API 服务，处理报告生成逻辑
3. **数据源层**：直接读取 Git 仓库数据，获取提交历史和变更信息

### 数据流向

1. 用户在 Web 界面输入仓库路径和报告参数
2. 前端通过 HTTP API 发送请求到后端服务
3. 后端服务调用 Git 命令获取仓库数据
4. 报告生成器处理数据并使用模板渲染
5. 生成的报告返回给前端进行展示和下载

## 技术栈

### 后端
- **Go 1.21+**：主要编程语言
- **Gorilla Mux**：HTTP 路由
- **CORS**：跨域支持
- **Git**：版本控制系统集成

### 前端
- **Next.js 14**：React 框架
- **TypeScript**：类型安全
- **Tailwind CSS**：样式框架
- **Lucide React**：图标库

### 容器化
- **Docker**：容器化平台
- **Docker Compose**：多容器应用编排
- **Alpine Linux**：轻量级基础镜像
- **Multi-stage Build**：优化镜像大小

## 注意事项

1. 确保在 Git 仓库目录中运行，或使用 `-repo` 参数指定仓库路径
2. 确保已配置 Git 用户信息（`git config user.name`）
3. 工具依赖 `git` 命令，请确保 Git 已正确安装并在 PATH 中
4. 周报默认按周一到周日计算
5. Web 界面需要同时启动后端服务器和前端开发服务器
6. 默认端口：后端 8080，前端 3000
7. Docker 部署时需要将 Git 仓库目录挂载到容器中
8. 使用 Docker 时确保有足够的磁盘空间用于构建镜像

## 示例输出

```markdown
# 张三 的工作日报

**日期：** 2024年01月15日
**仓库：** my-project (main分支)
**生成时间：** 2024-01-15 18:30:00

---

## 📊 今日统计

- **提交次数：** 5 次
- **修改文件：** 8 个
- **新增代码：** 156 行
- **删除代码：** 23 行

### 🔥 主要修改文件
- src/main.go (3次)
- README.md (2次)
- config/app.yaml (1次)

### 📁 文件类型分布
- Go: 5
- Markdown: 2
- 配置文件: 1

---

## 📝 详细提交记录

### 功能开发 (3个提交)

**a1b2c3d4** - 2024-01-15 09:30:00
> feat: 添加用户认证功能
> 📁 修改文件: src/auth.go, src/user.go
> 📈 代码变更: +89 -5

...
```

## 开发

### 项目结构
```
git-report-generator/
├── main.go                    # 主程序入口
├── server.go                  # HTTP服务器
├── git.go                    # Git操作相关
├── report.go                 # 报告生成逻辑
├── renderer.go               # 模板渲染
├── go.mod                    # Go模块依赖
├── Dockerfile                # 后端Docker配置
├── docker-compose.yml        # 容器编排配置
├── docker-deploy.sh          # Linux/macOS部署脚本
├── docker-deploy.bat         # Windows部署脚本
├── .dockerignore             # Docker忽略文件
├── templates/                # 报告模板
│   └── custom-template.tmpl
└── frontend/                 # Next.js前端
    ├── app/
    │   ├── globals.css
    │   ├── layout.tsx
    │   └── page.tsx
    ├── Dockerfile            # 前端Docker配置
    ├── .dockerignore         # 前端Docker忽略文件
    ├── package.json
    ├── next.config.js
    ├── tailwind.config.js
    └── tsconfig.json
```

### 开发环境

**本地开发：**
```bash
# 后端开发
go run . -server

# 前端开发
cd frontend
npm run dev
```

**Docker 开发：**
```bash
# 构建开发镜像
docker-compose -f docker-compose.dev.yml up --build

# 查看日志
docker-compose logs -f
```

### 构建生产版本

**Docker 部署（推荐）：**
```bash
# 一键部署
./docker-deploy.sh  # Linux/macOS
docker-deploy.bat   # Windows

# 或手动部署
docker-compose up -d --build
```

**传统构建：**

```bash
# 构建后端
go build -ldflags "-s -w" -o git-report

# 构建前端
cd frontend
npm run build
```

### Docker 镜像管理

```bash
# 查看镜像
docker images | grep git-report

# 清理未使用的镜像
docker image prune

# 重新构建镜像
docker-compose build --no-cache

# 推送到镜像仓库（可选）
docker tag git-report-backend:latest your-registry/git-report-backend:latest
docker push your-registry/git-report-backend:latest
```

## 故障排除

### 常见问题

#### Docker 相关

**Q: Docker 构建失败，提示网络连接问题**
```bash
# 解决方案：使用国内镜像源
docker-compose -f docker-compose.china.yml up -d --build
```

**Q: 容器启动后无法访问服务**
```bash
# 检查容器状态
docker-compose ps

# 查看容器日志
docker-compose logs backend
docker-compose logs frontend

# 检查端口占用
netstat -an | findstr :8080  # Windows
lsof -i :8080               # Linux/macOS
```

#### Git 相关

**Q: 提示 "not a git repository"**
- 确保在 Git 仓库目录中运行
- 或使用 `-repo` 参数指定正确的仓库路径

**Q: 无法获取 Git 用户信息**
```bash
# 配置 Git 用户信息
git config --global user.name "Your Name"
git config --global user.email "your.email@example.com"
```

**Q: 生成的报告为空**
- 检查指定日期范围内是否有提交记录
- 确认作者名称是否正确
- 使用 `git log --author="作者名" --since="2024-01-01"` 验证

#### 网络相关

**Q: 前端无法连接后端 API**
- 确认后端服务已启动（默认端口 8080）
- 检查防火墙设置
- 确认 CORS 配置正确

### 性能优化

**大型仓库优化：**
```bash
# 限制提交历史深度
git log --since="1 week ago" --oneline

# 使用浅克隆
git clone --depth 100 <repository-url>
```

**内存使用优化：**
- 对于大型仓库，建议使用 Docker 部署以限制内存使用
- 可以通过 `docker-compose.yml` 中的 `mem_limit` 参数控制内存限制

## 贡献指南

我们欢迎社区贡献！请遵循以下步骤：

### 开发流程

1. **Fork 项目**
   ```bash
   git clone https://github.com/your-username/git-report-generator.git
   cd git-report-generator
   ```

2. **创建功能分支**
   ```bash
   git checkout -b feature/your-feature-name
   ```

3. **本地开发**
   ```bash
   # 启动开发环境
   go run . -server
   
   # 在另一个终端启动前端
   cd frontend
   npm run dev
   ```

4. **运行测试**
   ```bash
   # 后端测试
   go test ./...
   
   # 前端测试
   cd frontend
   npm test
   ```

5. **提交更改**
   ```bash
   git add .
   git commit -m "feat: add your feature description"
   git push origin feature/your-feature-name
   ```

6. **创建 Pull Request**
   - 在 GitHub 上创建 PR
   - 详细描述更改内容
   - 确保通过所有检查

### 代码规范

**Go 代码：**
- 使用 `gofmt` 格式化代码
- 遵循 Go 官方编码规范
- 添加必要的注释和文档

**前端代码：**
- 使用 TypeScript 严格模式
- 遵循 ESLint 规则
- 使用 Prettier 格式化代码

**提交信息：**
- 使用约定式提交格式
- 格式：`type(scope): description`
- 类型：feat, fix, docs, style, refactor, test, chore

### 报告问题

在提交 Issue 时，请包含：
- 操作系统和版本
- Go 版本和 Node.js 版本
- 详细的错误信息
- 复现步骤
- 相关的日志输出

## 更新日志

详细的更新历史请查看 [CHANGELOG.md](CHANGELOG.md)

### 最新版本亮点

- ✨ 新增 Web 界面支持
- 🐳 Docker 一键部署
- 🎨 自定义模板功能
- 📊 增强的统计信息
- 🌐 多语言支持
- 🔧 改进的错误处理

## 路线图

### 计划中的功能

- [ ] **AI 增强**：集成 AI 模型自动生成工作总结
- [ ] **多仓库支持**：同时分析多个 Git 仓库
- [ ] **团队报告**：生成团队级别的工作报告
- [ ] **集成支持**：支持 GitLab、GitHub API
- [ ] **数据导出**：支持 PDF、Excel 格式导出
- [ ] **定时任务**：自动定时生成和发送报告
- [ ] **移动端适配**：响应式设计优化
- [ ] **插件系统**：支持第三方插件扩展

### 版本规划

- **v2.0**：AI 增强和多仓库支持
- **v2.1**：团队协作功能
- **v2.2**：企业级集成

## 社区

### 获取帮助

- 📖 [文档](https://github.com/phuhao00/git-report-generator/wiki)
- 💬 [讨论区](https://github.com/phuhao00/git-report-generator/discussions)
- 🐛 [问题反馈](https://github.com/phuhao00/git-report-generator/issues)

### 贡献者

感谢所有为项目做出贡献的开发者！

<!-- 贡献者列表将自动更新 -->

## 许可证

MIT License

Copyright (c) 2024 Git Report Generator Contributors

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

---

<div align="center">

**⭐ 如果这个项目对你有帮助，请给我们一个 Star！**

[🏠 主页](https://github.com/phuhao00/git-report-generator) • 
[📖 文档](https://github.com/phuhao00/git-report-generator/wiki) • 
[🐛 报告问题](https://github.com/phuhao00/git-report-generator/issues) • 
[💡 功能建议](https://github.com/phuhao00/git-report-generator/discussions)

</div>
