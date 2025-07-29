package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type GenerateReportRequest struct {
	RepoPath string `json:"repoPath"`
	Type     string `json:"type"`
	Date     string `json:"date"`
	Author   string `json:"author,omitempty"`
}

type GenerateReportResponse struct {
	Content string `json:"content"`
	Type    string `json:"type"`
	Date    string `json:"date"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type OptimizeReportRequest struct {
	Content string `json:"content"`
}

type OptimizeReportResponse struct {
	OptimizedContent string `json:"optimizedContent"`
}

func enableCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
}

func generateReportHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req GenerateReportRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid JSON format"})
		return
	}

	// 验证输入
	if req.RepoPath == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Repository path is required"})
		return
	}

	// 检查仓库路径是否存在
	if _, err := os.Stat(req.RepoPath); os.IsNotExist(err) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Repository path does not exist"})
		return
	}

	// 检查是否为Git仓库
	gitDir := filepath.Join(req.RepoPath, ".git")
	if _, err := os.Stat(gitDir); os.IsNotExist(err) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Path is not a Git repository"})
		return
	}

	// 解析日期
	targetDate, err := parseDate(req.Date)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: fmt.Sprintf("Invalid date format: %v", err)})
		return
	}

	// 创建报告生成器
	generator := NewReportGenerator(req.RepoPath, req.Author)

	// 生成报告
	var report *Report
	switch req.Type {
	case "daily":
		report, err = generator.GenerateDailyReport(targetDate)
	case "weekly":
		report, err = generator.GenerateWeeklyReport(targetDate)
	default:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid report type. Use 'daily' or 'weekly'"})
		return
	}

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: fmt.Sprintf("Failed to generate report: %v", err)})
		return
	}

	// 渲染报告
	renderer := NewReportRenderer("")
	content, err := renderer.Render(report)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: fmt.Sprintf("Failed to render report: %v", err)})
		return
	}

	// 返回成功响应
	response := GenerateReportResponse{
		Content: content,
		Type:    req.Type,
		Date:    req.Date,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func optimizeReportHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req OptimizeReportRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid JSON format"})
		return
	}

	if req.Content == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Content is required"})
		return
	}

	// 调用免费大模型API进行优化
	optimizedContent, err := optimizeWithAI(req.Content)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: fmt.Sprintf("Failed to optimize report: %v", err)})
		return
	}

	// 返回优化后的内容
	response := OptimizeReportResponse{
		OptimizedContent: optimizedContent,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func optimizeWithAI(content string) (string, error) {
	// 使用智谱AI的Chat Completions API
	apiURL := getAIAPIURL()
	apiKey := getAIAPIKey()
	
	if apiURL == "" || apiKey == "" {
		return "", fmt.Errorf("AI API configuration not found. Please set AI_API_URL and AI_API_KEY environment variables")
	}

	// 构建请求体 - 智谱AI格式
	requestBody := map[string]interface{}{
		"model": "glm-4-flash", // 智谱AI的免费模型
		"messages": []map[string]string{
			{
				"role": "system",
				"content": "你是一个专业的技术文档优化助手。请优化以下Git提交报告，使其更加清晰、专业和易读。保持原有的结构和信息完整性，但改进语言表达、格式和可读性。请用中文回复。",
			},
			{
				"role": "user",
				"content": content,
			},
		},
		"max_tokens": 2000,
		"temperature": 0.7,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %v", err)
	}

	// 发送HTTP请求
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// 解析响应
	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("failed to parse response: %v", err)
	}

	// 提取优化后的内容
	choices, ok := response["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		return "", fmt.Errorf("invalid response format: no choices found")
	}

	firstChoice, ok := choices[0].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("invalid response format: invalid choice format")
	}

	message, ok := firstChoice["message"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("invalid response format: no message found")
	}

	optimizedContent, ok := message["content"].(string)
	if !ok {
		return "", fmt.Errorf("invalid response format: no content found")
	}

	return strings.TrimSpace(optimizedContent), nil
}

func getAIAPIURL() string {
	// 优先使用环境变量，如果没有则使用默认的免费API
	if url := os.Getenv("AI_API_URL"); url != "" {
		return url
	}
	
	// 默认使用DeepSeek API (提供免费额度)
	return "https://api.deepseek.com/v1/chat/completions"
}

func getAIAPIKey() string {
	// 优先使用环境变量
	if key := os.Getenv("AI_API_KEY"); key != "" {
		return key
	}
	
	// 如果没有设置环境变量，返回空字符串
	// 用户需要自己设置API密钥
	return ""
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func startServer() {
	r := mux.NewRouter()

	// API routes
	r.HandleFunc("/api/generate-report", generateReportHandler).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/optimize-report", optimizeReportHandler).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/health", healthHandler).Methods("GET", "OPTIONS")

	// Setup CORS
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"*"},
	})

	handler := c.Handler(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server starting on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}