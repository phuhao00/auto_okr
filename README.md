# Git 提交记录报告生成器

一个基于 Git 提交记录自动生成日报和周报的 Golang 工具。

## 功能特性

- 📅 **自动生成日报和周报**：基于 Git 提交记录自动生成工作报告
- 👤 **多用户支持**：可指定特定作者或使用当前 Git 用户
- 📊 **详细统计信息**：包含提交次数、代码行数、文件类型分布等
- 🏷️ **智能分类**：自动将提交按功能开发、Bug修复、重构等类别分组
- 🎨 **自定义模板**：支持使用自定义模板定制报告格式
- 💾 **灵活输出**：支持控制台输出或保存到文件

## 安装

```bash
# 克隆项目
git clone <repository-url>
cd git-report-generator

# 编译
go build -o git-report.exe
```

## 使用方法

### 基本用法

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

## 注意事项

1. 确保在 Git 仓库目录中运行，或使用 `-repo` 参数指定仓库路径
2. 确保已配置 Git 用户信息（`git config user.name`）
3. 工具依赖 `git` 命令，请确保 Git 已正确安装并在 PATH 中
4. 周报默认按周一到周日计算

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

## 许可证

MIT License