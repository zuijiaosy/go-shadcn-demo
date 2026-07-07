package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// doRequest 执行一次测试请求并返回响应记录器
func doRequest(t *testing.T, method, path string) *httptest.ResponseRecorder {
	t.Helper()
	server := NewServer()
	req := httptest.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	server.Handler().ServeHTTP(w, req)
	return w
}

// TestHealthCheck 测试健康检查端点
func TestHealthCheck(t *testing.T) {
	w := doRequest(t, http.MethodGet, "/api/health")

	if w.Code != http.StatusOK {
		t.Fatalf("期望状态码 200，实际 %d", w.Code)
	}

	var body map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("响应不是合法 JSON: %v", err)
	}
	if body["status"] != "ok" {
		t.Errorf("期望 status=ok，实际 %q", body["status"])
	}
}

// TestGetVersion 测试版本信息端点
func TestGetVersion(t *testing.T) {
	w := doRequest(t, http.MethodGet, "/api/version")

	if w.Code != http.StatusOK {
		t.Fatalf("期望状态码 200，实际 %d", w.Code)
	}

	var body map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("响应不是合法 JSON: %v", err)
	}
	for _, key := range []string{"version", "git_commit", "build_time", "go_version"} {
		if body[key] == "" {
			t.Errorf("响应缺少字段 %q", key)
		}
	}
}

// TestAPINotFound 测试未知 API 路径返回 404 及错误格式
func TestAPINotFound(t *testing.T) {
	w := doRequest(t, http.MethodGet, "/api/not-exist")

	if w.Code != http.StatusNotFound {
		t.Fatalf("期望状态码 404，实际 %d", w.Code)
	}

	var body map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("响应不是合法 JSON: %v", err)
	}
	if body["error"] == "" {
		t.Error("错误响应应包含 error 字段")
	}
}
