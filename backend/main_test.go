package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthCheck(t *testing.T) {
	handler := setupRoutes()

	req, err := http.NewRequest(http.MethodGet, "/api/health", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, status)
	}

	var resp map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to parse response JSON: %v", err)
	}

	if resp["status"] != "UP" {
		t.Errorf("Expected status 'UP', got %v", resp["status"])
	}
}

func TestInfoEndpoint(t *testing.T) {
	handler := setupRoutes()

	req, err := http.NewRequest(http.MethodGet, "/api/info", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, status)
	}

	var resp map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to parse response JSON: %v", err)
	}

	if resp["app"] == nil {
		t.Errorf("Expected 'app' field in info response")
	}
}

func TestMessagesCRUD(t *testing.T) {
	handler := setupRoutes()

	// 1. Get messages
	getReq, _ := http.NewRequest(http.MethodGet, "/api/messages", nil)
	getRR := httptest.NewRecorder()
	handler.ServeHTTP(getRR, getReq)
	if getRR.Code != http.StatusOK {
		t.Fatalf("Expected 200 on GET /api/messages, got %d", getRR.Code)
	}

	// 2. Post new message
	payload := []byte(`{"content":"CI/CD Test Message","author":"Jenkins Test"}`)
	postReq, _ := http.NewRequest(http.MethodPost, "/api/messages", bytes.NewBuffer(payload))
	postReq.Header.Set("Content-Type", "application/json")
	postRR := httptest.NewRecorder()
	handler.ServeHTTP(postRR, postReq)
	if postRR.Code != http.StatusCreated {
		t.Fatalf("Expected 201 on POST /api/messages, got %d", postRR.Code)
	}

	var postResp struct {
		Data Message `json:"data"`
	}
	if err := json.Unmarshal(postRR.Body.Bytes(), &postResp); err != nil {
		t.Fatalf("Failed to unmarshal post response: %v", err)
	}
	createdID := postResp.Data.ID

	// 3. Delete created message
	delReq, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/messages?id=%d", createdID), nil)
	delRR := httptest.NewRecorder()
	handler.ServeHTTP(delRR, delReq)
	if delRR.Code != http.StatusOK {
		t.Fatalf("Expected 200 on DELETE /api/messages, got %d", delRR.Code)
	}
}
