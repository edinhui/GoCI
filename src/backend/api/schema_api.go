package api

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"goci/backend/storage"
)

// SchemaHandler 处理Schema相关的API请求
type SchemaHandler struct {
	storage *storage.SchemaStorage
}

// NewSchemaHandler 创建一个新的SchemaHandler实例
func NewSchemaHandler(storage *storage.SchemaStorage) *SchemaHandler {
	return &SchemaHandler{
		storage: storage,
	}
}

// SaveSchema 处理保存Schema的请求
func (h *SchemaHandler) SaveSchema(c *gin.Context) {
	// 从URL参数获取Schema ID
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Schema ID is required"})
		return
	}

	// 解析请求体
	var requestBody struct {
		Metadata struct {
			Name        string `json:"name"`
			Description string `json:"description"`
		} `json:"metadata"`
		Schema json.RawMessage `json:"schema"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse request body: " + err.Error()})
		return
	}

	// 获取元数据
	name := requestBody.Metadata.Name
	if name == "" {
		name = id // 如果没有提供名称，使用ID作为名称
	}
	description := requestBody.Metadata.Description

	// 将schema转换为字节数组
	schemaData, err := json.Marshal(requestBody.Schema)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to process schema data: " + err.Error()})
		return
	}

	// 保存Schema
	if err := h.storage.SaveSchema(id, name, description, schemaData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 构建统一格式的响应
	c.JSON(http.StatusOK, gin.H{
		"metadata": gin.H{
			"id":          id,
			"name":        name,
			"description": description,
		},
		"message": "Schema saved successfully",
	})
}

// GetSchema 处理获取Schema的请求
func (h *SchemaHandler) GetSchema(c *gin.Context) {
	// 从URL参数获取Schema ID
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Schema ID is required"})
		return
	}

	// 获取Schema
	schemaData, metadata, err := h.storage.GetSchema(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// 解析schema数据为JSON
	var schemaJSON map[string]interface{}
	if err := json.Unmarshal(schemaData, &schemaJSON); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse schema data"})
		return
	}

	// 构建响应对象
	response := gin.H{
		"metadata": gin.H{
			"id":          id,
			"name":        metadata.Name,
			"description": metadata.Description,
			"createdAt":   metadata.CreatedAt,
			"updatedAt":   metadata.UpdatedAt,
		},
		"schema": schemaJSON,
	}

	// 返回统一格式的JSON响应
	c.JSON(http.StatusOK, response)
}

// ListSchemas 处理列出所有Schema的请求
func (h *SchemaHandler) ListSchemas(c *gin.Context) {
	// 获取所有Schema
	schemas, err := h.storage.ListSchemas()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 构建统一格式的响应
	c.JSON(http.StatusOK, gin.H{"schemas": schemas})
}

// DeleteSchema 处理删除Schema的请求
func (h *SchemaHandler) DeleteSchema(c *gin.Context) {
	// 从URL参数获取Schema ID
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Schema ID is required"})
		return
	}

	// 删除Schema
	if err := h.storage.DeleteSchema(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Schema deleted successfully", "id": id})
}

// RegisterRoutes 注册API路由
func RegisterRoutes(r *gin.Engine, storage *storage.SchemaStorage) {
	// 创建处理器
	handler := NewSchemaHandler(storage)

	// 创建API组
	api := r.Group("/api")
	{
		// Schema API
		schemas := api.Group("/schemas")
		{
			// 保存Schema
			schemas.POST("/:id", handler.SaveSchema)
			// 获取Schema
			schemas.GET("/:id", handler.GetSchema)
			// 列出所有Schema
			schemas.GET("", handler.ListSchemas)
			// 删除Schema
			schemas.DELETE("/:id", handler.DeleteSchema)
		}
	}
}
