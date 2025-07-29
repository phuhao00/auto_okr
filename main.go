package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	var (
		reportType = flag.String("type", "daily", "报告类型: daily, weekly")
		date = flag.String("date", "", "指定日期 (YYYY-MM-DD), 默认为今天")
		repoPath = flag.String("repo", ".", "Git仓库路径")
		author = flag.String("author", "", "指定作者，默认为当前Git用户")
		output = flag.String("output", "", "输出文件路径，默认输出到控制台")
		template = flag.String("template", "", "自定义模板文件路径")
		server = flag.Bool("server", false, "启动HTTP服务器模式")
	)
	flag.Parse()

	// 如果是服务器模式，启动HTTP服务器
	if *server {
		startServer()
		return
	}

	// 解析日期
	targetDate, err := parseDate(*date)
	if err != nil {
		log.Fatalf("日期解析错误: %v", err)
	}

	// 创建报告生成器
	generator := NewReportGenerator(*repoPath, *author)

	// 生成报告
	var report *Report
	switch *reportType {
	case "daily":
		report, err = generator.GenerateDailyReport(targetDate)
	case "weekly":
		report, err = generator.GenerateWeeklyReport(targetDate)
	default:
		log.Fatalf("不支持的报告类型: %s", *reportType)
	}

	if err != nil {
		log.Fatalf("生成报告失败: %v", err)
	}

	// 渲染报告
	renderer := NewReportRenderer(*template)
	content, err := renderer.Render(report)
	if err != nil {
		log.Fatalf("渲染报告失败: %v", err)
	}

	// 输出报告
	if *output != "" {
		err = os.WriteFile(*output, []byte(content), 0644)
		if err != nil {
			log.Fatalf("写入文件失败: %v", err)
		}
		fmt.Printf("报告已保存到: %s\n", *output)
	} else {
		fmt.Print(content)
	}
}

func parseDate(dateStr string) (time.Time, error) {
	if dateStr == "" {
		return time.Now(), nil
	}
	return time.Parse("2006-01-02", dateStr)
}