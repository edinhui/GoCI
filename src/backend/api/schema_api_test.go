package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	"goci/backend/storage"
)

// 测试辅助函数：创建临时目录
func createTempDir(t *testing.T) string {
	tempDir, err := os.MkdirTemp("", "schema-api-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	return tempDir
}

// 测试辅助函数：清理临时目录
func cleanupTempDir(t *testing.T, dir string) {
	if err := os.RemoveAll(dir); err != nil {
		t.Fatalf("Failed to cleanup temp dir: %v", err)
	}
}

// 测试辅助函数：设置测试环境
func setupTest(t *testing.T) (*gin.Engine, *storage.SchemaStorage, string) {
	// 创建临时目录
	tempDir := createTempDir(t)

	// 保存当前工作目录
	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current working directory: %v", err)
	}

	// 切换到临时目录
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Failed to change working directory: %v", err)
	}

	// 创建Gin引擎
	gin.SetMode(gin.TestMode)
	r := gin.New()

	// 创建存储实例
	schemaStorage := storage.NewSchemaStorage()

	// 注册API路由
	RegisterRoutes(r, schemaStorage)

	return r, schemaStorage, oldWd
}

// 测试SaveSchema API
func TestSaveSchemaAPI(t *testing.T) {
	// 设置测试环境
	r, _, oldWd := setupTest(t)
	defer os.Chdir(oldWd)

	// 测试数据
	id := "test-schema"
	name := "Test Schema"
	description := "A test schema"
	schemaData := []byte(`{"type": "object", "properties": {"name": {"type": "string"}}}`)

	// 创建请求
	req := httptest.NewRequest(http.MethodPost, "/api/schemas/"+id, bytes.NewBuffer(schemaData))
	req.Header.Set("X-Schema-Name", name)
	req.Header.Set("X-Schema-Description", description)
	req.Header.Set("Content-Type", "application/json")

	// 创建响应记录器
	w := httptest.NewRecorder()

	// 处理请求
	r.ServeHTTP(w, req)

	// 验证响应状态码
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	// 验证响应内容
	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if message, ok := response["message"].(string); !ok || message != "Schema saved successfully" {
		t.Errorf("Expected message 'Schema saved successfully', got %v", response["message"])
	}

	if responseID, ok := response["id"].(string); !ok || responseID != id {
		t.Errorf("Expected id '%s', got %v", id, response["id"])
	}

	// 验证Schema文件是否已创建
	schemaDir := filepath.Join(".", "schemas", id)
	schemaPath := filepath.Join(schemaDir, "schema.json")
	if _, err := os.Stat(schemaPath); os.IsNotExist(err) {
		t.Errorf("Schema file was not created")
	}

	// 验证Schema内容是否正确
	data, err := os.ReadFile(schemaPath)
	if err != nil {
		t.Fatalf("Failed to read schema file: %v", err)
	}

	if string(data) != string(schemaData) {
		t.Errorf("Schema content is incorrect: got %s, want %s", string(data), string(schemaData))
	}
}

// 测试GetSchema API
func TestGetSchemaAPI(t *testing.T) {
	// 设置测试环境
	r, schemaStorage, oldWd := setupTest(t)
	defer os.Chdir(oldWd)

	// 测试数据
	id := "test-schema"
	name := "Test Schema"
	description := "A test schema"
	schemaData := []byte(`{"type": "object", "properties": {"name": {"type": "string"}}}`)

	// 保存Schema
	if err := schemaStorage.SaveSchema(id, name, description, schemaData); err != nil {
		t.Fatalf("Failed to save schema: %v", err)
	}

	// 创建请求
	req := httptest.NewRequest(http.MethodGet, "/api/schemas/"+id, nil)

	// 创建响应记录器
	w := httptest.NewRecorder()

	// 处理请求
	r.ServeHTTP(w, req)

	// 验证响应状态码
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	// 验证响应头
	if w.Header().Get("Content-Type") != "application/json" {
		t.Errorf("Expected Content-Type 'application/json', got %s", w.Header().Get("Content-Type"))
	}

	if w.Header().Get("X-Schema-Name") != name {
		t.Errorf("Expected X-Schema-Name '%s', got %s", name, w.Header().Get("X-Schema-Name"))
	}

	if w.Header().Get("X-Schema-Description") != description {
		t.Errorf("Expected X-Schema-Description '%s', got %s", description, w.Header().Get("X-Schema-Description"))
	}

	// 验证响应内容
	if string(w.Body.Bytes()) != string(schemaData) {
		t.Errorf("Schema content is incorrect: got %s, want %s", string(w.Body.Bytes()), string(schemaData))
	}

	// 测试获取不存在的Schema
	req = httptest.NewRequest(http.MethodGet, "/api/schemas/non-existent", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// 验证响应状态码
	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, w.Code)
	}
}

// 测试ListSchemas API
func TestListSchemasAPI(t *testing.T) {
	// 设置测试环境
	r, schemaStorage, oldWd := setupTest(t)
	defer os.Chdir(oldWd)

	// 测试数据
	schemas := []struct {
		id          string
		name        string
		description string
		data        []byte
	}{
		{
			id:          "schema1",
			name:        "Schema 1",
			description: "First test schema",
			data:        []byte(`{"type": "object", "properties": {"name": {"type": "string"}}}`),
		},
		{
			id:          "schema2",
			name:        "Schema 2",
			description: "Second test schema",
			data:        []byte(`{"type": "object", "properties": {"age": {"type": "integer"}}}`),
		},
	}

	// 保存Schema
	for _, schema := range schemas {
		if err := schemaStorage.SaveSchema(schema.id, schema.name, schema.description, schema.data); err != nil {
			t.Fatalf("Failed to save schema: %v", err)
		}
	}

	// 创建请求
	req := httptest.NewRequest(http.MethodGet, "/api/schemas", nil)

	// 创建响应记录器
	w := httptest.NewRecorder()

	// 处理请求
	r.ServeHTTP(w, req)

	// 验证响应状态码
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	// 验证响应内容
	var response []storage.SchemaMetadata
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	// 验证列表长度是否正确
	if len(response) != len(schemas) {
		t.Errorf("Schema list length is incorrect: got %d, want %d", len(response), len(schemas))
	}

	// 验证列表内容是否正确
	for _, schema := range schemas {
		found := false
		for _, metadata := range response {
			if metadata.ID == schema.id {
				found = true
				if metadata.Name != schema.name {
					t.Errorf("Schema name is incorrect: got %s, want %s", metadata.Name, schema.name)
				}
				if metadata.Description != schema.description {
					t.Errorf("Schema description is incorrect: got %s, want %s", metadata.Description, schema.description)
				}
				break
			}
		}
		if !found {
			t.Errorf("Schema %s not found in list", schema.id)
		}
	}
}

// 测试DeleteSchema API
func TestDeleteSchemaAPI(t *testing.T) {
	// 设置测试环境
	r, schemaStorage, oldWd := setupTest(t)
	defer os.Chdir(oldWd)

	// 测试数据
	id := "test-schema"
	name := "Test Schema"
	description := "A test schema"
	schemaData := []byte(`{"type": "object", "properties": {"name": {"type": "string"}}}`)

	// 保存Schema
	if err := schemaStorage.SaveSchema(id, name, description, schemaData); err != nil {
		t.Fatalf("Failed to save schema: %v", err)
	}

	// 创建请求
	req := httptest.NewRequest(http.MethodDelete, "/api/schemas/"+id, nil)

	// 创建响应记录器
	w := httptest.NewRecorder()

	// 处理请求
	r.ServeHTTP(w, req)

	// 验证响应状态码
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	// 验证响应内容
	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if message, ok := response["message"].(string); !ok || message != "Schema deleted successfully" {
		t.Errorf("Expected message 'Schema deleted successfully', got %v", response["message"])
	}

	if responseID, ok := response["id"].(string); !ok || responseID != id {
		t.Errorf("Expected id '%s', got %v", id, response["id"])
	}

	// 验证Schema目录是否已删除
	schemaDir := filepath.Join(".", "schemas", id)
	if _, err := os.Stat(schemaDir); !os.IsNotExist(err) {
		t.Errorf("Schema directory was not deleted")
	}

	// 测试删除不存在的Schema
	req = httptest.NewRequest(http.MethodDelete, "/api/schemas/non-existent", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// 验证响应状态码
	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, w.Code)
	}
}
