package main

import (
	"net/http"

	"go-redis-mq/internal/dashboard"
	"go-redis-mq/internal/redisclient"

	"github.com/gin-gonic/gin"
)

func main() {
	redisclient.InitRedis()

	r := gin.Default()

	// API 返回 JSON 数据
	r.GET("/api/stats", func(c *gin.Context) {
		data, err := dashboard.GetStats()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, data)
	})

	// 提供静态文件用于展示图表页面
	r.StaticFile("/", "./web/static/index.html")

	r.Run(":8080") // http://localhost:8080
}
