package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// Build variables injected during CI/CD with -ldflags
var (
	Version   = "dev"
	GitCommit = "unknown"
	BuildTime = "unknown"
	startTime = time.Now()
)

type Message struct {
	ID        int       `json:"id"`
	Content   string    `json:"content" binding:"required"`
	Author    string    `json:"author"`
	CreatedAt time.Time `json:"created_at"`
}

type CIJob struct {
	ID          int       `json:"id"`
	Name        string    `json:"name" binding:"required"`
	Branch      string    `json:"branch"`
	Status      string    `json:"status"` // pending, running, success, failed
	DurationSec int       `json:"duration_sec"`
	CreatedAt   time.Time `json:"created_at"`
}

type AppStore struct {
	mu       sync.RWMutex
	messages []Message
	nextMsgID int
	jobs     []CIJob
	nextJobID int
}

var store = &AppStore{
	messages: []Message{
		{ID: 1, Content: "Hello from Gin Engine & Jenkins CI/CD!", Author: "DevOps Bot", CreatedAt: time.Now().Add(-10 * time.Minute)},
		{ID: 2, Content: "Vue 3 + Gin backend online.", Author: "Antigravity", CreatedAt: time.Now().Add(-5 * time.Minute)},
	},
	nextMsgID: 3,
	jobs: []CIJob{
		{ID: 101, Name: "build-backend-image", Branch: "main", Status: "success", DurationSec: 42, CreatedAt: time.Now().Add(-15 * time.Minute)},
		{ID: 102, Name: "vue-unit-tests", Branch: "main", Status: "success", DurationSec: 18, CreatedAt: time.Now().Add(-10 * time.Minute)},
		{ID: 103, Name: "deploy-k8s-staging", Branch: "feature/ci", Status: "running", DurationSec: 65, CreatedAt: time.Now().Add(-2 * time.Minute)},
	},
	nextJobID: 104,
}

// corsMiddleware configures CORS headers for local development
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization, X-Requested-With, X-Custom-Header")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

func setupRouter() *gin.Engine {
	// Set release mode if APP_ENV is production
	if os.Getenv("APP_ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.Logger())
	r.Use(corsMiddleware())

	api := r.Group("/api")
	{
		// 1. Health & Meta
		api.GET("/health", handleHealth)
		api.GET("/info", handleInfo)
		api.GET("/metrics", handleMetrics)

		// 2. Message CRUD
		api.GET("/messages", handleGetMessages)
		api.POST("/messages", handleCreateMessage)
		api.DELETE("/messages/:id", handleDeleteMessage)

		// 3. Request Echo & Inspector (Header/Body testing)
		api.POST("/echo", handleEcho)

		// 4. HTTP Status Code Simulator (Error & Circuit Break testing)
		api.GET("/mock/status", handleMockStatus)

		// 5. CI/CD Job Task Simulator
		api.GET("/jobs", handleGetJobs)
		api.POST("/jobs", handleCreateJob)
	}

	return r
}

func handleHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "UP",
		"timestamp": time.Now().Format(time.RFC3339),
		"uptime":    time.Since(startTime).String(),
	})
}

func handleInfo(c *gin.Context) {
	hostname, _ := os.Hostname()
	c.JSON(http.StatusOK, gin.H{
		"app":         "Vue + Go (Gin) Demo Service",
		"framework":   "Gin " + gin.Version,
		"version":     Version,
		"git_commit":  GitCommit,
		"build_time":  BuildTime,
		"go_version":  runtime.Version(),
		"os_arch":     fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
		"hostname":    hostname,
		"environment": getEnv("APP_ENV", "production"),
	})
}

func handleMetrics(c *gin.Context) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	c.JSON(http.StatusOK, gin.H{
		"num_cpu":        runtime.NumCPU(),
		"num_goroutine":  runtime.NumGoroutine(),
		"alloc_mb":       fmt.Sprintf("%.2f MB", float64(m.Alloc)/(1024*1024)),
		"total_alloc_mb": fmt.Sprintf("%.2f MB", float64(m.TotalAlloc)/(1024*1024)),
		"sys_mb":         fmt.Sprintf("%.2f MB", float64(m.Sys)/(1024*1024)),
		"num_gc":         m.NumGC,
	})
}

func handleGetMessages(c *gin.Context) {
	store.mu.RLock()
	msgs := make([]Message, len(store.messages))
	copy(msgs, store.messages)
	store.mu.RUnlock()

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": msgs,
	})
}

func handleCreateMessage(c *gin.Context) {
	var req struct {
		Content string `json:"content" binding:"required"`
		Author  string `json:"author"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数校验失败: content 字段必填"})
		return
	}
	if req.Author == "" {
		req.Author = "Anonymous"
	}

	store.mu.Lock()
	newMsg := Message{
		ID:        store.nextMsgID,
		Content:   req.Content,
		Author:    req.Author,
		CreatedAt: time.Now(),
	}
	store.nextMsgID++
	store.messages = append(store.messages, newMsg)
	store.mu.Unlock()

	c.JSON(http.StatusCreated, gin.H{
		"code": 0,
		"msg":  "created",
		"data": newMsg,
	})
}

func handleDeleteMessage(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "非法 message ID"})
		return
	}

	store.mu.Lock()
	found := false
	newMsgs := make([]Message, 0, len(store.messages))
	for _, m := range store.messages {
		if m.ID == id {
			found = true
		} else {
			newMsgs = append(newMsgs, m)
		}
	}
	if found {
		store.messages = newMsgs
	}
	store.mu.Unlock()

	if !found {
		c.JSON(http.StatusNotFound, gin.H{"error": "未找到对应留言"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "deleted"})
}

func handleEcho(c *gin.Context) {
	var body any
	_ = c.ShouldBindJSON(&body)

	headers := make(map[string]string)
	for k, v := range c.Request.Header {
		if len(v) > 0 {
			headers[k] = v[0]
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"method":     c.Request.Method,
		"path":       c.Request.URL.Path,
		"query":      c.Request.URL.Query(),
		"client_ip":  c.ClientIP(),
		"headers":    headers,
		"body":       body,
		"received_at": time.Now().Format(time.RFC3339),
	})
}

func handleMockStatus(c *gin.Context) {
	codeStr := c.DefaultQuery("code", "200")
	statusCode, err := strconv.Atoi(codeStr)
	if err != nil || statusCode < 100 || statusCode > 599 {
		statusCode = http.StatusBadRequest
	}

	message := fmt.Sprintf("模拟响应: HTTP %d (%s)", statusCode, http.StatusText(statusCode))
	if statusCode >= 400 {
		c.JSON(statusCode, gin.H{
			"code":    statusCode,
			"error":   message,
			"details": "该接口由前端主动请求模拟，用于测试前端错误处理及 Jenkins 探测熔断",
		})
		return
	}

	c.JSON(statusCode, gin.H{
		"code":    0,
		"message": message,
	})
}

func handleGetJobs(c *gin.Context) {
	store.mu.RLock()
	jobs := make([]CIJob, len(store.jobs))
	copy(jobs, store.jobs)
	store.mu.RUnlock()

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": jobs,
	})
}

func handleCreateJob(c *gin.Context) {
	var req struct {
		Name   string `json:"name" binding:"required"`
		Branch string `json:"branch"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数校验失败: name 字段必填"})
		return
	}
	if req.Branch == "" {
		req.Branch = "main"
	}

	store.mu.Lock()
	newJob := CIJob{
		ID:          store.nextJobID,
		Name:        req.Name,
		Branch:      req.Branch,
		Status:      "running",
		DurationSec: 0,
		CreatedAt:   time.Now(),
	}
	store.nextJobID++
	store.jobs = append([]CIJob{newJob}, store.jobs...) // prepend
	store.mu.Unlock()

	c.JSON(http.StatusCreated, gin.H{
		"code": 0,
		"msg":  "CI Job 模拟任务已创建",
		"data": newJob,
	})
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func main() {
	port := getEnv("PORT", "8080")
	r := setupRouter()

	log.Printf("Starting Gin Web Server on port %s...", port)
	log.Printf("Version: %s, GitCommit: %s, BuildTime: %s", Version, GitCommit, BuildTime)

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
