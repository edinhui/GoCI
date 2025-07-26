package storage

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

// 测试辅助函数：创建临时目录
func createTempDir(t *testing.T) string {
	// 使用系统临时目录
	tempDir, err := os.MkdirTemp(os.TempDir(), "schema-storage-test")
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

// 测试SchemaStorage的创建
func TestNewSchemaStorage(t *testing.T) {
	// 创建临时目录
	tempDir := createTempDir(t)
	defer cleanupTempDir(t, tempDir)

	// 保存当前工作目录
	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current working directory: %v", err)
	}

	// 切换到临时目录
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Failed to change working directory: %v", err)
	}
	defer os.Chdir(oldWd)

	// 创建存储实例
	storage := NewSchemaStorage()
	
	// 确保存储实例创建成功
	if storage == nil {
		t.Fatalf("Failed to create storage instance")
	}

	// 验证存储目录是否已创建
	schemasDir := filepath.Join(".", "schemas")
	if _, err := os.Stat(schemasDir); os.IsNotExist(err) {
		t.Errorf("Schemas directory was not created")
	}

	// 验证注册表文件是否已创建
	registryPath := filepath.Join(schemasDir, "schema-registry.json")
	if _, err := os.Stat(registryPath); os.IsNotExist(err) {
		t.Errorf("Schema registry file was not created")
	}

	// 验证注册表是否为空
	data, err := os.ReadFile(registryPath)
	if err != nil {
		t.Fatalf("Failed to read registry file: %v", err)
	}

	var registry map[string]SchemaMetadata
	if err := json.Unmarshal(data, &registry); err != nil {
		t.Fatalf("Failed to parse registry file: %v", err)
	}

	if len(registry) != 0 {
		t.Errorf("Registry is not empty: %v", registry)
	}
}

// 测试保存Schema
func TestSaveSchema(t *testing.T) {
	// 创建临时目录
	tempDir := createTempDir(t)
	defer cleanupTempDir(t, tempDir)

	// 保存当前工作目录
	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current working directory: %v", err)
	}

	// 切换到临时目录
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Failed to change working directory: %v", err)
	}
	defer os.Chdir(oldWd)

	// 创建存储实例
	storage := NewSchemaStorage()

	// 测试数据
	id := "test-schema"
	name := "Test Schema"
	description := "A test schema"
	schemaData := []byte(`{"type": "object", "properties": {"name": {"type": "string"}}}`)

	// 保存Schema
	if err := storage.SaveSchema(id, name, description, schemaData); err != nil {
		t.Fatalf("Failed to save schema: %v", err)
	}

	// 验证Schema目录是否已创建
	schemaDir := filepath.Join(".", "schemas", id)
	if _, err := os.Stat(schemaDir); os.IsNotExist(err) {
		t.Errorf("Schema directory was not created")
	}

	// 验证Schema文件是否已创建
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

	// 验证注册表是否已更新
	registryPath := filepath.Join(".", "schemas", "schema-registry.json")
	data, err = os.ReadFile(registryPath)
	if err != nil {
		t.Fatalf("Failed to read registry file: %v", err)
	}

	var registry map[string]SchemaMetadata
	if err := json.Unmarshal(data, &registry); err != nil {
		t.Fatalf("Failed to parse registry file: %v", err)
	}

	metadata, exists := registry[id]
	if !exists {
		t.Errorf("Schema metadata not found in registry")
	}

	if metadata.ID != id {
		t.Errorf("Schema ID is incorrect: got %s, want %s", metadata.ID, id)
	}

	if metadata.Name != name {
		t.Errorf("Schema name is incorrect: got %s, want %s", metadata.Name, name)
	}

	if metadata.Description != description {
		t.Errorf("Schema description is incorrect: got %s, want %s", metadata.Description, description)
	}

	// 验证创建时间和更新时间是否已设置
	if metadata.CreatedAt == "" {
		t.Errorf("Schema created time is not set")
	}

	if metadata.UpdatedAt == "" {
		t.Errorf("Schema updated time is not set")
	}
}

// 测试获取Schema
func TestGetSchema(t *testing.T) {
	// 创建临时目录
	tempDir := createTempDir(t)
	defer cleanupTempDir(t, tempDir)

	// 保存当前工作目录
	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current working directory: %v", err)
	}

	// 切换到临时目录
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Failed to change working directory: %v", err)
	}
	defer os.Chdir(oldWd)

	// 创建存储实例
	storage := NewSchemaStorage()

	// 测试数据
	id := "test-schema"
	name := "Test Schema"
	description := "A test schema"
	schemaData := []byte(`{"type": "object", "properties": {"name": {"type": "string"}}}`)

	// 保存Schema
	if err := storage.SaveSchema(id, name, description, schemaData); err != nil {
		t.Fatalf("Failed to save schema: %v", err)
	}

	// 获取Schema
	data, metadata, err := storage.GetSchema(id)
	if err != nil {
		t.Fatalf("Failed to get schema: %v", err)
	}

	// 验证Schema内容是否正确
	if string(data) != string(schemaData) {
		t.Errorf("Schema content is incorrect: got %s, want %s", string(data), string(schemaData))
	}

	// 验证元数据是否正确
	if metadata.ID != id {
		t.Errorf("Schema ID is incorrect: got %s, want %s", metadata.ID, id)
	}

	if metadata.Name != name {
		t.Errorf("Schema name is incorrect: got %s, want %s", metadata.Name, name)
	}

	if metadata.Description != description {
		t.Errorf("Schema description is incorrect: got %s, want %s", metadata.Description, description)
	}

	// 验证不存在的Schema
	_, _, err = storage.GetSchema("non-existent")
	if err == nil {
		t.Errorf("Expected error when getting non-existent schema")
	}
}

// 测试列出所有Schema
func TestListSchemas(t *testing.T) {
	// 创建临时目录
	tempDir := createTempDir(t)
	defer cleanupTempDir(t, tempDir)

	// 保存当前工作目录
	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current working directory: %v", err)
	}

	// 切换到临时目录
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Failed to change working directory: %v", err)
	}
	defer os.Chdir(oldWd)

	// 创建存储实例
	storage := NewSchemaStorage()

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
		if err := storage.SaveSchema(schema.id, schema.name, schema.description, schema.data); err != nil {
			t.Fatalf("Failed to save schema: %v", err)
		}
	}

	// 列出所有Schema
	list, err := storage.ListSchemas()
	if err != nil {
		t.Fatalf("Failed to list schemas: %v", err)
	}

	// 验证列表长度是否正确
	if len(list) != len(schemas) {
		t.Errorf("Schema list length is incorrect: got %d, want %d", len(list), len(schemas))
	}

	// 验证列表内容是否正确
	for _, schema := range schemas {
		found := false
		for _, metadata := range list {
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

// 测试删除Schema
func TestDeleteSchema(t *testing.T) {
	// 创建临时目录
	tempDir := createTempDir(t)
	defer cleanupTempDir(t, tempDir)

	// 保存当前工作目录
	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current working directory: %v", err)
	}

	// 切换到临时目录
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Failed to change working directory: %v", err)
	}
	defer os.Chdir(oldWd)

	// 创建存储实例
	storage := NewSchemaStorage()

	// 测试数据
	id := "test-schema"
	name := "Test Schema"
	description := "A test schema"
	schemaData := []byte(`{"type": "object", "properties": {"name": {"type": "string"}}}`)

	// 保存Schema
	if err := storage.SaveSchema(id, name, description, schemaData); err != nil {
		t.Fatalf("Failed to save schema: %v", err)
	}

	// 删除Schema
	if err := storage.DeleteSchema(id); err != nil {
		t.Fatalf("Failed to delete schema: %v", err)
	}

	// 验证Schema目录是否已删除
	schemaDir := filepath.Join(".", "schemas", id)
	if _, err := os.Stat(schemaDir); !os.IsNotExist(err) {
		t.Errorf("Schema directory was not deleted")
	}

	// 验证注册表是否已更新
	registryPath := filepath.Join(".", "schemas", "schema-registry.json")
	data, err := os.ReadFile(registryPath)
	if err != nil {
		t.Fatalf("Failed to read registry file: %v", err)
	}

	var registry map[string]SchemaMetadata
	if err := json.Unmarshal(data, &registry); err != nil {
		t.Fatalf("Failed to parse registry file: %v", err)
	}

	if _, exists := registry[id]; exists {
		t.Errorf("Schema metadata still exists in registry")
	}

	// 验证删除不存在的Schema
	err = storage.DeleteSchema("non-existent")
	if err == nil {
		t.Errorf("Expected error when deleting non-existent schema")
	}
}
