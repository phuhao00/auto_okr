package main

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

// Report 报告结构
type Report struct {
	Type        string            // daily, weekly
	Date        time.Time         // 报告日期
	Period      string            // 时间范围描述
	Author      string            // 作者
	RepoInfo    map[string]string // 仓库信息
	Commits     []*GitCommit      // 提交记录
	Summary     *ReportSummary    // 统计摘要
	Categories  map[string][]*GitCommit // 按类别分组的提交
	GeneratedAt time.Time         // 生成时间
}

// ReportSummary 报告摘要
type ReportSummary struct {
	TotalCommits   int               // 总提交数
	TotalFiles     int               // 总文件数
	TotalAdditions int               // 总新增行数
	TotalDeletions int               // 总删除行数
	FileTypes      map[string]int    // 文件类型统计
	DailyStats     map[string]int    // 每日统计（仅周报）
	TopFiles       []string          // 修改最多的文件
}

// ReportGenerator 报告生成器
type ReportGenerator struct {
	gitParser *GitParser
	author    string
}

// NewReportGenerator 创建报告生成器
func NewReportGenerator(repoPath, author string) *ReportGenerator {
	gitParser := NewGitParser(repoPath)
	
	// 如果没有指定作者，尝试获取当前Git用户
	if author == "" {
		if currentUser, err := gitParser.GetCurrentUser(); err == nil {
			author = currentUser
		}
	}
	
	return &ReportGenerator{
		gitParser: gitParser,
		author:    author,
	}
}

// GenerateDailyReport 生成日报
func (rg *ReportGenerator) GenerateDailyReport(date time.Time) (*Report, error) {
	// 获取当天的提交记录
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24 * time.Hour).Add(-time.Second)
	
	commits, err := rg.gitParser.GetCommits(startOfDay, endOfDay, rg.author)
	if err != nil {
		return nil, err
	}
	
	repoInfo, _ := rg.gitParser.GetRepoInfo()
	
	report := &Report{
		Type:        "daily",
		Date:        date,
		Period:      date.Format("2006年01月02日"),
		Author:      rg.author,
		RepoInfo:    repoInfo,
		Commits:     commits,
		Summary:     rg.generateSummary(commits, false),
		Categories:  rg.categorizeCommits(commits),
		GeneratedAt: time.Now(),
	}
	
	return report, nil
}

// GenerateWeeklyReport 生成周报
func (rg *ReportGenerator) GenerateWeeklyReport(date time.Time) (*Report, error) {
	// 获取本周的开始和结束时间（周一到周日）
	weekday := int(date.Weekday())
	if weekday == 0 { // 周日
		weekday = 7
	}
	startOfWeek := date.AddDate(0, 0, -(weekday-1))
	startOfWeek = time.Date(startOfWeek.Year(), startOfWeek.Month(), startOfWeek.Day(), 0, 0, 0, 0, startOfWeek.Location())
	endOfWeek := startOfWeek.AddDate(0, 0, 7).Add(-time.Second)
	
	commits, err := rg.gitParser.GetCommits(startOfWeek, endOfWeek, rg.author)
	if err != nil {
		return nil, err
	}
	
	repoInfo, _ := rg.gitParser.GetRepoInfo()
	
	report := &Report{
		Type:     "weekly",
		Date:     date,
		Period:   fmt.Sprintf("%s 至 %s", startOfWeek.Format("2006年01月02日"), endOfWeek.Format("2006年01月02日")),
		Author:   rg.author,
		RepoInfo: repoInfo,
		Commits:  commits,
		Summary:  rg.generateSummary(commits, true),
		Categories: rg.categorizeCommits(commits),
		GeneratedAt: time.Now(),
	}
	
	return report, nil
}

// generateSummary 生成统计摘要
func (rg *ReportGenerator) generateSummary(commits []*GitCommit, isWeekly bool) *ReportSummary {
	summary := &ReportSummary{
		FileTypes: make(map[string]int),
		DailyStats: make(map[string]int),
	}
	
	fileCount := make(map[string]int)
	
	for _, commit := range commits {
		summary.TotalCommits++
		summary.TotalAdditions += commit.Additions
		summary.TotalDeletions += commit.Deletions
		
		// 统计每日提交数（仅周报）
		if isWeekly {
			dayKey := commit.Date.Format("01-02")
			summary.DailyStats[dayKey]++
		}
		
		// 统计文件
		for _, file := range commit.Files {
			fileCount[file]++
			
			// 统计文件类型
			ext := rg.getFileExtension(file)
			summary.FileTypes[ext]++
		}
	}
	
	summary.TotalFiles = len(fileCount)
	
	// 获取修改最多的文件（前5个）
	type fileFreq struct {
		file  string
		count int
	}
	
	var fileFreqs []fileFreq
	for file, count := range fileCount {
		fileFreqs = append(fileFreqs, fileFreq{file, count})
	}
	
	sort.Slice(fileFreqs, func(i, j int) bool {
		return fileFreqs[i].count > fileFreqs[j].count
	})
	
	for i, ff := range fileFreqs {
		if i >= 5 {
			break
		}
		summary.TopFiles = append(summary.TopFiles, fmt.Sprintf("%s (%d次)", ff.file, ff.count))
	}
	
	return summary
}

// categorizeCommits 按类别分组提交
func (rg *ReportGenerator) categorizeCommits(commits []*GitCommit) map[string][]*GitCommit {
	categories := make(map[string][]*GitCommit)
	
	for _, commit := range commits {
		category := rg.categorizeCommit(commit)
		categories[category] = append(categories[category], commit)
	}
	
	return categories
}

// categorizeCommit 根据提交信息分类
func (rg *ReportGenerator) categorizeCommit(commit *GitCommit) string {
	message := strings.ToLower(commit.Message)
	
	// 功能开发
	if strings.Contains(message, "feat") || strings.Contains(message, "feature") || 
	   strings.Contains(message, "add") || strings.Contains(message, "新增") ||
	   strings.Contains(message, "功能") {
		return "功能开发"
	}
	
	// Bug修复
	if strings.Contains(message, "fix") || strings.Contains(message, "bug") ||
	   strings.Contains(message, "修复") || strings.Contains(message, "修正") {
		return "Bug修复"
	}
	
	// 重构
	if strings.Contains(message, "refactor") || strings.Contains(message, "重构") {
		return "代码重构"
	}
	
	// 文档
	if strings.Contains(message, "doc") || strings.Contains(message, "readme") ||
	   strings.Contains(message, "文档") {
		return "文档更新"
	}
	
	// 测试
	if strings.Contains(message, "test") || strings.Contains(message, "测试") {
		return "测试相关"
	}
	
	// 配置
	if strings.Contains(message, "config") || strings.Contains(message, "配置") {
		return "配置修改"
	}
	
	return "其他"
}

// getFileExtension 获取文件扩展名
func (rg *ReportGenerator) getFileExtension(filename string) string {
	parts := strings.Split(filename, ".")
	if len(parts) < 2 {
		return "无扩展名"
	}
	ext := strings.ToLower(parts[len(parts)-1])
	
	// 合并一些常见的扩展名
	switch ext {
	case "js", "jsx", "ts", "tsx":
		return "JavaScript/TypeScript"
	case "go":
		return "Go"
	case "py":
		return "Python"
	case "java":
		return "Java"
	case "cpp", "cc", "cxx", "c":
		return "C/C++"
	case "html", "htm":
		return "HTML"
	case "css", "scss", "sass":
		return "CSS"
	case "md", "markdown":
		return "Markdown"
	case "json", "yaml", "yml", "xml":
		return "配置文件"
	default:
		return ext
	}
}