package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

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