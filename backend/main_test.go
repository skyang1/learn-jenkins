package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestGinHealthCheck(t *testing.T) {
	router := setupRouter()

	req, _ := http.NewRequest(http.MethodGet, "/api/health", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("Expected status code 200, got %d", rr.Code)
	}

	var resp map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to parse response JSON: %v", err)
	}

	if resp["status"] != "UP" {
		t.Errorf("Expected status 'UP', got %v", resp["status"])
	}
}

func TestGinMetrics(t *testing.T) {
	router := setupRouter()

	req, _ := http.NewRequest(http.MethodGet, "/api/metrics", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", rr.Code)
	}

	var resp map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to parse metrics JSON: %v", err)
	}

	if resp["num_goroutine"] == nil || resp["alloc_mb"] == nil {
		t.Errorf("Metrics missing expected fields: %v", resp)
	}
}

func TestGinEcho(t *testing.T) {
	router := setupRouter()

	payload := []byte(`{"test_key":"hello_gin"}`)
	req, _ := http.NewRequest(http.MethodPost, "/api/echo?trace_id=12345", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Custom-Header", "Jenkins-Test")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", rr.Code)
	}

	var resp map[string]any
	_ = json.Unmarshal(rr.Body.Bytes(), &resp)

	if resp["method"] != "POST" {
		t.Errorf("Expected method POST, got %v", resp["method"])
	}
}

func TestGinMockStatus(t *testing.T) {
	router := setupRouter()

	// Test 400
	req400, _ := http.NewRequest(http.MethodGet, "/api/mock/status?code=400", nil)
	rr400 := httptest.NewRecorder()
	router.ServeHTTP(rr400, req400)
	if rr400.Code != http.StatusBadRequest {
		t.Errorf("Expected 400, got %d", rr400.Code)
	}

	// Test 500
	req500, _ := http.NewRequest(http.MethodGet, "/api/mock/status?code=500", nil)
	rr500 := httptest.NewRecorder()
	router.ServeHTTP(rr500, req500)
	if rr500.Code != http.StatusInternalServerError {
		t.Errorf("Expected 500, got %d", rr500.Code)
	}
}

func TestGinMessagesCRUD(t *testing.T) {
	router := setupRouter()

	// 1. Get messages
	getReq, _ := http.NewRequest(http.MethodGet, "/api/messages", nil)
	getRR := httptest.NewRecorder()
	router.ServeHTTP(getRR, getReq)
	if getRR.Code != http.StatusOK {
		t.Fatalf("Expected 200 on GET /api/messages, got %d", getRR.Code)
	}

	// 2. Post new message
	payload := []byte(`{"content":"Gin Framework Test","author":"CI Bot"}`)
	postReq, _ := http.NewRequest(http.MethodPost, "/api/messages", bytes.NewBuffer(payload))
	postReq.Header.Set("Content-Type", "application/json")
	postRR := httptest.NewRecorder()
	router.ServeHTTP(postRR, postReq)
	if postRR.Code != http.StatusCreated {
		t.Fatalf("Expected 201 on POST /api/messages, got %d: %s", postRR.Code, postRR.Body.String())
	}

	var postResp struct {
		Data Message `json:"data"`
	}
	_ = json.Unmarshal(postRR.Body.Bytes(), &postResp)
	createdID := postResp.Data.ID

	// 3. Delete message
	delReq, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/messages/%d", createdID), nil)
	delRR := httptest.NewRecorder()
	router.ServeHTTP(delRR, delReq)
	if delRR.Code != http.StatusOK {
		t.Fatalf("Expected 200 on DELETE /api/messages/%d, got %d", createdID, delRR.Code)
	}
}

func TestGinJobs(t *testing.T) {
	router := setupRouter()

	// 1. Get jobs
	getReq, _ := http.NewRequest(http.MethodGet, "/api/jobs", nil)
	getRR := httptest.NewRecorder()
	router.ServeHTTP(getRR, getReq)
	if getRR.Code != http.StatusOK {
		t.Fatalf("Expected 200 on GET /api/jobs, got %d", getRR.Code)
	}

	// 2. Post new job
	payload := []byte(`{"name":"sonar-code-scan","branch":"main"}`)
	postReq, _ := http.NewRequest(http.MethodPost, "/api/jobs", bytes.NewBuffer(payload))
	postReq.Header.Set("Content-Type", "application/json")
	postRR := httptest.NewRecorder()
	router.ServeHTTP(postRR, postReq)
	if postRR.Code != http.StatusCreated {
		t.Fatalf("Expected 201 on POST /api/jobs, got %d", postRR.Code)
	}
}
