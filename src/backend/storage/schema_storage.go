package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// SchemaStorage 处理JSON Schema的存储和检索
type SchemaStorage struct {
	mutex sync.RWMutex
	// 配置目录路径
	schemasDir string
	// Schema注册表文件路径
	registryPath string
	// Schema注册表（内存中的缓存）
	registry map[string]SchemaMetadata
}

// SchemaMetadata 表示Schema的元数据
type SchemaMetadata struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

// NewSchemaStorage 创建一个新的SchemaStorage实例
func NewSchemaStorage() *SchemaStorage {
	// 创建存储目录
	schemasDir := filepath.Join(".", "schemas")
	os.MkdirAll(schemasDir, os.ModePerm)

	// 注册表文件路径
	registryPath := filepath.Join(schemasDir, "schema-registry.json")

	// 初始化存储
	storage := &SchemaStorage{
		schemasDir:   schemasDir,
		registryPath: registryPath,
		registry:     make(map[string]SchemaMetadata),
	}

	// 加载注册表（无锁版本，避免初始化时的死锁）
	storage.loadRegistryNoLock()

	return storage
}

// loadRegistryNoLock 从文件加载Schema注册表（无锁版本，仅在初始化时使用）
func (s *SchemaStorage) loadRegistryNoLock() {
	// 检查注册表文件是否存在
	if _, err := os.Stat(s.registryPath); os.IsNotExist(err) {
		// 如果不存在，创建一个空的注册表
		s.registry = make(map[string]SchemaMetadata)
		s.saveRegistryNoLock() // 保存空注册表
		return
	}

	// 读取注册表文件
	data, err := os.ReadFile(s.registryPath)
	if err != nil {
		fmt.Printf("Error reading registry file: %v\n", err)
		return
	}

	// 解析JSON
	if err := json.Unmarshal(data, &s.registry); err != nil {
		fmt.Printf("Error parsing registry file: %v\n", err)
		return
	}
}

// loadRegistry 从文件加载Schema注册表
func (s *SchemaStorage) loadRegistry() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// 检查注册表文件是否存在
	if _, err := os.Stat(s.registryPath); os.IsNotExist(err) {        
		// 如果不存在，创建一个空的注册表
		s.registry = make(map[string]SchemaMetadata)
		// 使用无锁版本，避免死锁
		s.saveRegistryNoLock()
		return
	}

	// 读取注册表文件
	data, err := os.ReadFile(s.registryPath)
	if err != nil {
		fmt.Printf("Error reading registry file: %v\n", err)      
		return
	}

	// 解析JSON
	if err := json.Unmarshal(data, &s.registry); err != nil {
		fmt.Printf("Error parsing registry file: %v\n", err)      
		return
	}
}

// saveRegistryNoLock 将Schema注册表保存到文件（无锁版本）
func (s *SchemaStorage) saveRegistryNoLock() error {
	// 将注册表转换为JSON
	data, err := json.MarshalIndent(s.registry, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling registry: %w", err)
	}

	// 写入文件
	if err := os.WriteFile(s.registryPath, data, 0644); err != nil {
		return fmt.Errorf("error writing registry file: %w", err)
	}

	return nil
}

// saveRegistry 将Schema注册表保存到文件
func (s *SchemaStorage) saveRegistry() error {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	// 将注册表转换为JSON
	data, err := json.MarshalIndent(s.registry, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling registry: %w", err)
	}

	// 写入文件
	if err := os.WriteFile(s.registryPath, data, 0644); err != nil {
		return fmt.Errorf("error writing registry file: %w", err)
	}

	return nil
}

// SaveSchema 保存JSON Schema
func (s *SchemaStorage) SaveSchema(id string, name string, description string, schemaData []byte) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// 创建Schema目录
	schemaDir := filepath.Join(s.schemasDir, id)
	if err := os.MkdirAll(schemaDir, os.ModePerm); err != nil {
		return fmt.Errorf("error creating schema directory: %w", err)
	}

	// 保存Schema文件
	schemaPath := filepath.Join(schemaDir, "schema.json")
	if err := os.WriteFile(schemaPath, schemaData, 0644); err != nil {
		return fmt.Errorf("error writing schema file: %w", err)
	}

	// 更新元数据
	now := time.Now().Format(time.RFC3339)
	metadata := SchemaMetadata{
		ID:          id,
		Name:        name,
		Description: description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// 如果已存在，保留创建时间
	if existing, exists := s.registry[id]; exists {
		metadata.CreatedAt = existing.CreatedAt
	}

	// 更新注册表
	s.registry[id] = metadata

	// 保存注册表（使用无锁版本，避免死锁）
	if err := s.saveRegistryNoLock(); err != nil {
		return fmt.Errorf("error saving registry: %w", err)
	}

	return nil
}

// GetSchema 获取指定ID的Schema
func (s *SchemaStorage) GetSchema(id string) ([]byte, SchemaMetadata, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	// 检查Schema是否存在
	metadata, exists := s.registry[id]
	if !exists {
		return nil, SchemaMetadata{}, fmt.Errorf("schema not found: %s", id)
	}

	// 读取Schema文件
	schemaPath := filepath.Join(s.schemasDir, id, "schema.json")
	data, err := os.ReadFile(schemaPath)
	if err != nil {
		return nil, SchemaMetadata{}, fmt.Errorf("error reading schema file: %w", err)
	}

	return data, metadata, nil
}

// ListSchemas 列出所有可用的Schema
func (s *SchemaStorage) ListSchemas() ([]SchemaMetadata, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	// 将注册表中的所有元数据转换为切片
	schemas := make([]SchemaMetadata, 0, len(s.registry))
	for _, metadata := range s.registry {
		schemas = append(schemas, metadata)
	}

	return schemas, nil
}

// DeleteSchema 删除指定ID的Schema
func (s *SchemaStorage) DeleteSchema(id string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// 检查Schema是否存在
	if _, exists := s.registry[id]; !exists {
		return fmt.Errorf("schema not found: %s", id)
	}

	// 删除Schema目录
	schemaDir := filepath.Join(s.schemasDir, id)
	if err := os.RemoveAll(schemaDir); err != nil {
		return fmt.Errorf("error deleting schema directory: %w", err)
	}

	// 从注册表中删除
	delete(s.registry, id)

	// 保存注册表（使用无锁版本，避免死锁）
	if err := s.saveRegistryNoLock(); err != nil {
		return fmt.Errorf("error saving registry: %w", err)
	}

	return nil
}
