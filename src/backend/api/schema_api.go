package api

import (
	"io"
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

	// 从请求体获取Schema数据
	schemaData, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
		return
	}

	// 从请求头获取Schema名称和描述
	name := c.GetHeader("X-Schema-Name")
	if name == "" {
		name = id // 如果没有提供名称，使用ID作为名称
	}
	description := c.GetHeader("X-Schema-Description")

	// 保存Schema
	if err := h.storage.SaveSchema(id, name, description, schemaData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Schema saved successfully", "id": id})
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

	// 设置响应头
	c.Header("Content-Type", "application/json")
	c.Header("X-Schema-Name", metadata.Name)
	c.Header("X-Schema-Description", metadata.Description)
	c.Header("X-Schema-Created-At", metadata.CreatedAt)
	c.Header("X-Schema-Updated-At", metadata.UpdatedAt)

	// 返回Schema数据
	c.Data(http.StatusOK, "application/json", schemaData)
}

// ListSchemas 处理列出所有Schema的请求
func (h *SchemaHandler) ListSchemas(c *gin.Context) {
	// 获取所有Schema
	schemas, err := h.storage.ListSchemas()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, schemas)
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
