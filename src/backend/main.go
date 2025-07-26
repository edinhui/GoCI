package main

import (
	"log"
	"net/http"

	"goci/backend/api"
	"goci/backend/storage"

	"github.com/gin-gonic/gin"
)

// corsMiddleware 创建一个CORS中间件，允许前端与后端API进行通信
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Schema-Name, X-Schema-Description")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func main() {
	// 创建Gin引擎
	r := gin.Default()

	// 添加CORS中间件
	r.Use(corsMiddleware())

	// 创建存储服务
	schemaStorage := storage.NewSchemaStorage()

	// 注册API路由
	api.RegisterRoutes(r, schemaStorage)

	// 启动服务器
	log.Println("Starting server on :8080...")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
