package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
)

func doHealthCheck(t *testing.T) map[string]interface{} {
	t.Helper()
	e := echo.New()
	e.GET("/health", HealthCheck)
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
	var body map[string]interface{}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	return body
}

func TestAllHealthy(t *testing.T) {
	body := doHealthCheck(t)
	if body["status"] != "healthy" {
		t.Errorf("expected healthy, got %v", body["status"])
	}
	checks := body["checks"].(map[string]interface{})
	for name, v := range checks {
		m := v.(map[string]interface{})
		if m["status"] != "ok" {
			t.Errorf("check %s: expected ok, got %v", name, m["status"])
		}
	}
}

func TestDatabaseFailure(t *testing.T) {
	orig := CheckDatabase
	CheckDatabase = func() map[string]interface{} {
		return map[string]interface{}{"status": "error"}
	}
	defer func() { CheckDatabase = orig }()

	body := doHealthCheck(t)
	if body["status"] != "unhealthy" {
		t.Errorf("expected unhealthy, got %v", body["status"])
	}
}

func TestRedisFailure(t *testing.T) {
	orig := CheckRedis
	CheckRedis = func() map[string]interface{} {
		return map[string]interface{}{"status": "error"}
	}
	defer func() { CheckRedis = orig }()

	body := doHealthCheck(t)
	if body["status"] != "unhealthy" {
		t.Errorf("expected unhealthy, got %v", body["status"])
	}
}

func TestAllDependenciesDown(t *testing.T) {
	origDB := CheckDatabase
	origRedis := CheckRedis
	origDisk := CheckDisk
	CheckDatabase = func() map[string]interface{} {
		return map[string]interface{}{"status": "error"}
	}
	CheckRedis = func() map[string]interface{} {
		return map[string]interface{}{"status": "error"}
	}
	CheckDisk = func() map[string]interface{} {
		return map[string]interface{}{"status": "error"}
	}
	defer func() {
		CheckDatabase = origDB
		CheckRedis = origRedis
		CheckDisk = origDisk
	}()

	body := doHealthCheck(t)
	if body["status"] != "unhealthy" {
		t.Errorf("expected unhealthy, got %v", body["status"])
	}
}

func TestDiskCheckReturnsAvailable(t *testing.T) {
	body := doHealthCheck(t)
	checks := body["checks"].(map[string]interface{})
	disk := checks["disk"].(map[string]interface{})
	if disk["status"] != "ok" {
		t.Errorf("expected disk ok, got %v", disk["status"])
	}
	if _, exists := disk["available"]; !exists {
		t.Error("disk check missing 'available' field")
	}
}

func TestResponseStructure(t *testing.T) {
	body := doHealthCheck(t)
	for _, key := range []string{"status", "uptime", "checks"} {
		if _, exists := body[key]; !exists {
			t.Errorf("missing top-level key: %s", key)
		}
	}
	checks := body["checks"].(map[string]interface{})
	for _, key := range []string{"database", "redis", "disk"} {
		if _, exists := checks[key]; !exists {
			t.Errorf("missing check: %s", key)
		}
	}
}
