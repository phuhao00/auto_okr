package main

import (
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// GitCommit 表示一个Git提交
type GitCommit struct {
	Hash      string
	Author    string
	Date      time.Time
	Message   string
	Files     []string
	Additions int
	Deletions int
}

// GitParser Git解析器
type GitParser struct {
	repoPath string
}

// NewGitParser 创建新的Git解析器
func NewGitParser(repoPath string) *GitParser {
	return &GitParser{
		repoPath: repoPath,
	}
}

// GetCommits 获取指定时间范围内的提交记录
func (g *GitParser) GetCommits(since, until time.Time, author string) ([]*GitCommit, error) {
	args := []string{
		"log",
		"--pretty=format:%H|%an|%ad|%s",
		"--date=iso",
		"--numstat",
		fmt.Sprintf("--since=%s", since.Format("2006-01-02 00:00:00")),
		fmt.Sprintf("--until=%s", until.Format("2006-01-02 23:59:59")),
	}

	if author != "" {
		args = append(args, fmt.Sprintf("--author=%s", author))
	}

	cmd := exec.Command("git", args...)
	cmd.Dir = g.repoPath

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("执行git命令失败: %v", err)
	}

	return g.parseCommits(string(output))
}

// GetCurrentUser 获取当前Git用户
func (g *GitParser) GetCurrentUser() (string, error) {
	cmd := exec.Command("git", "config", "user.name")
	cmd.Dir = g.repoPath

	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("获取Git用户失败: %v", err)
	}

	return strings.TrimSpace(string(output)), nil
}

// parseCommits 解析Git日志输出
func (g *GitParser) parseCommits(output string) ([]*GitCommit, error) {
	var commits []*GitCommit
	lines := strings.Split(output, "\n")

	var currentCommit *GitCommit
	commitRegex := regexp.MustCompile(`^([a-f0-9]+)\|(.+)\|(.+)\|(.+)$`)
	numstatRegex := regexp.MustCompile(`^(\d+)\s+(\d+)\s+(.+)$`)

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// 检查是否是提交行
		if matches := commitRegex.FindStringSubmatch(line); matches != nil {
			// 保存上一个提交
			if currentCommit != nil {
				commits = append(commits, currentCommit)
			}

			// 解析日期
			date, err := time.Parse("2006-01-02 15:04:05 -0700", matches[3])
			if err != nil {
				return nil, fmt.Errorf("解析日期失败: %v", err)
			}

			currentCommit = &GitCommit{
				Hash:    matches[1],
				Author:  matches[2],
				Date:    date,
				Message: matches[4],
				Files:   []string{},
			}
		} else if currentCommit != nil {
			// 检查是否是numstat行
			if matches := numstatRegex.FindStringSubmatch(line); matches != nil {
				additions, _ := strconv.Atoi(matches[1])
				deletions, _ := strconv.Atoi(matches[2])
				filename := matches[3]

				currentCommit.Additions += additions
				currentCommit.Deletions += deletions
				currentCommit.Files = append(currentCommit.Files, filename)
			}
		}
	}

	// 添加最后一个提交
	if currentCommit != nil {
		commits = append(commits, currentCommit)
	}

	return commits, nil
}

// GetRepoInfo 获取仓库信息
func (g *GitParser) GetRepoInfo() (map[string]string, error) {
	info := make(map[string]string)

	// 获取仓库名称
	cmd := exec.Command("git", "remote", "get-url", "origin")
	cmd.Dir = g.repoPath
	if output, err := cmd.Output(); err == nil {
		remoteURL := strings.TrimSpace(string(output))
		// 从URL中提取仓库名
		parts := strings.Split(remoteURL, "/")
		if len(parts) > 0 {
			repoName := parts[len(parts)-1]
			repoName = strings.TrimSuffix(repoName, ".git")
			info["name"] = repoName
		}
		info["url"] = remoteURL
	}

	// 获取当前分支
	cmd = exec.Command("git", "branch", "--show-current")
	cmd.Dir = g.repoPath
	if output, err := cmd.Output(); err == nil {
		info["branch"] = strings.TrimSpace(string(output))
	}

	return info, nil
}