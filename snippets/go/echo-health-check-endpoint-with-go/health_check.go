package main

import (
	"fmt"
	"net/http"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
)

var startTime = time.Now()

// CheckDatabase pings the database. Replace with real logic.
var CheckDatabase = func() map[string]interface{} {
	// e.g. err := db.Ping(); if err != nil { return error status }
	return map[string]interface{}{"status": "ok"}
}

// CheckRedis pings Redis. Replace with real logic.
var CheckRedis = func() map[string]interface{} {
	// e.g. _, err := rdb.Ping(ctx).Result()
	return map[string]interface{}{"status": "ok"}
}

// CheckDisk reports available disk space for the root partition.
var CheckDisk = func() map[string]interface{} {
	var stat syscall.Statfs_t
	if err := syscall.Statfs("/", &stat); err != nil {
		return map[string]interface{}{"status": "error", "detail": err.Error()}
	}
	avail := stat.Bavail * uint64(stat.Bsize)
	return map[string]interface{}{
		"status":    "ok",
		"available": fmt.Sprintf("%d MB", avail/1024/1024),
	}
}

// HealthCheck handles GET /health.
func HealthCheck(c echo.Context) error {
	checks := map[string]interface{}{
		"database": CheckDatabase(),
		"redis":    CheckRedis(),
		"disk":     CheckDisk(),
	}
	healthy := true
	for _, v := range checks {
		if m, ok := v.(map[string]interface{}); ok {
			if s, _ := m["status"].(string); s != "ok" {
				healthy = false
			}
		}
	}
	status := "healthy"
	if !healthy {
		status = "unhealthy"
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": status,
		"uptime": time.Since(startTime).String(),
		"checks": checks,
	})
}

func main() {
	e := echo.New()
	e.GET("/health", HealthCheck)
	e.Logger.Fatal(e.Start(":8080"))
}
