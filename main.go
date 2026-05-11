package main

import (
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
)

var startTime = time.Now()

type HealthResponse struct {
	Status    string  `json:"status"`
	Timestamp string  `json:"timestamp"`
	Uptime    float64 `json:"uptime"`
}

type DetailedHealthResponse struct {
	Status      string      `json:"status"`
	Timestamp   string      `json:"timestamp"`
	Uptime      float64     `json:"uptime"`
	Environment string      `json:"environment"`
	GoVersion   string      `json:"goVersion"`
	Memory      MemoryStats `json:"memory"`
}

type MemoryStats struct {
	Alloc      uint64 `json:"alloc"`
	TotalAlloc uint64 `json:"totalAlloc"`
	Sys        uint64 `json:"sys"`
	NumGC      uint32 `json:"numGc"`
}

type RootResponse struct {
	Message   string            `json:"message"`
	Endpoints map[string]string `json:"endpoints"`
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	// Root endpoint
	router.GET("/", func(c *gin.Context) {
		response := RootResponse{
			Message: "Demo API is running",
			Endpoints: map[string]string{
				"health":         "/health",
				"detailedHealth": "/health/detailed",
			},
		}
		c.JSON(http.StatusOK, response)
	})

	// Basic health check endpoint
	router.GET("/health", func(c *gin.Context) {
		response := HealthResponse{
			Status:    "UP",
			Timestamp: time.Now().UTC().Format(time.RFC3339),
			Uptime:    time.Since(startTime).Seconds(),
		}
		c.JSON(http.StatusOK, response)
	})

	// Detailed health check endpoint
	router.GET("/health/detailed", func(c *gin.Context) {
		var m runtime.MemStats
		runtime.ReadMemStats(&m)

		env := os.Getenv("ENV")
		if env == "" {
			env = "development"
		}

		response := DetailedHealthResponse{
			Status:      "UP",
			Timestamp:   time.Now().UTC().Format(time.RFC3339),
			Uptime:      time.Since(startTime).Seconds(),
			Environment: env,
			GoVersion:   runtime.Version(),
			Memory: MemoryStats{
				Alloc:      m.Alloc,
				TotalAlloc: m.TotalAlloc,
				Sys:        m.Sys,
				NumGC:      m.NumGC,
			},
		}
		c.JSON(http.StatusOK, response)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	router.Run(":" + port)
}
