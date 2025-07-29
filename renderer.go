package main

import (
	"bytes"
	"fmt"
	"os"
	"sort"
	"strings"
	"text/template"
	"time"
)

// ReportRenderer 报告渲染器
type ReportRenderer struct {
	templateFile string
}

// NewReportRenderer 创建报告渲染器
func NewReportRenderer(templateFile string) *ReportRenderer {
	return &ReportRenderer{
		templateFile: templateFile,
	}
}

// Render 渲染报告
func (rr *ReportRenderer) Render(report *Report) (string, error) {
	var templateContent string
	var err error
	
	if rr.templateFile != "" {
		// 使用自定义模板
		content, err := os.ReadFile(rr.templateFile)
		if err != nil {
			return "", fmt.Errorf("读取模板文件失败: %v", err)
		}
		templateContent = string(content)
	} else {
		// 使用默认模板
		if report.Type == "daily" {
			templateContent = rr.getDefaultDailyTemplate()
		} else {
			templateContent = rr.getDefaultWeeklyTemplate()
		}
	}
	
	// 创建模板函数
	funcMap := template.FuncMap{
		"formatTime": func(t time.Time) string {
			return t.Format("2006-01-02 15:04:05")
		},
		"formatDate": func(t time.Time) string {
			return t.Format("2006年01月02日")
		},
		"formatShortHash": func(hash string) string {
			if len(hash) > 8 {
				return hash[:8]
			}
			return hash
		},
		"join": strings.Join,
		"add": func(a, b int) int {
			return a + b
		},
		"sortedKeys": func(m map[string][]*GitCommit) []string {
			keys := make([]string, 0, len(m))
			for k := range m {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			return keys
		},
		"sortedFileTypes": func(m map[string]int) []string {
			type kv struct {
				key   string
				value int
			}
			var kvs []kv
			for k, v := range m {
				kvs = append(kvs, kv{k, v})
			}
			sort.Slice(kvs, func(i, j int) bool {
				return kvs[i].value > kvs[j].value
			})
			var result []string
			for _, kv := range kvs {
				result = append(result, fmt.Sprintf("%s: %d", kv.key, kv.value))
			}
			return result
		},
		"sub": func(a, b int) int {
			return a - b
		},
	}
	
	// 解析模板
	tmpl, err := template.New("report").Funcs(funcMap).Parse(templateContent)
	if err != nil {
		return "", fmt.Errorf("解析模板失败: %v", err)
	}
	
	// 渲染模板
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, report)
	if err != nil {
		return "", fmt.Errorf("执行模板失败: %v", err)
	}
	
	return buf.String(), nil
}

// getDefaultDailyTemplate 获取默认日报模板
func (rr *ReportRenderer) getDefaultDailyTemplate() string {
	return `{{range $index, $commit := .Commits}}
# {{add $index 1}}. {{$commit.Message}}
{{end}}
`
}

// getDefaultWeeklyTemplate 获取默认周报模板
func (rr *ReportRenderer) getDefaultWeeklyTemplate() string {
	return `{{range $index, $commit := .Commits}}
# {{add $index 1}}. {{$commit.Message}}
{{end}}
`
}