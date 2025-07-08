# Configuration Management Component Architecture

## 1. 前端技术选型 (Vue 3)

### JSON Schema 创建页面相关框架
- **图形化编辑器**：
  - 采用类似 **[vue-json-schema-editor-visual](https://github.com/zyqwst/json-schema-editor-vue3)** 的树状列表形式，提供图形化的编辑和生成界面
  - 支持拖拽操作、节点添加/删除/编辑等功能
  - 实时预览生成的 JSON Schema

- **编辑器组件**：
  - **Monaco Editor**：用于高级 JSON 编辑，与 VS Code 使用相同的编辑器引擎，支持语法高亮和自动完成
  - **JSON Editor**：功能强大的 JSON 编辑器，支持树视图、代码视图和表单视图
  - **vue-json-editor**：轻量级 JSON 编辑器组件，支持基本的编辑功能

### CSS 风格
- **UI 框架**：
  - **Element Plus**：专为 Vue 3 设计的组件库，提供完整的设计系统
    - 丰富的组件库，包括表单、表格、导航等
    - 全面的主题定制能力
    - 国际化支持
    - 活跃的社区和完善的文档

## 2. JSON 文件浏览和修改框架

- **JSON 编辑器**：
  - **vue-json-editor**：轻量级 JSON 编辑器组件
  - **jsoneditor**：功能强大的 JSON 编辑器，支持树视图、代码视图和表单视图
  - **vue-json-pretty**：美观的 JSON 查看器和编辑器

- **树形结构展示**：
  - **vue3-tree-vue**：用于展示层级结构
  - **vue3-json-viewer**：专门用于 JSON 数据的可折叠视图

## 3. 后端 Webserver 框架

考虑到需要将前端页面打包进 Golang 程序，推荐以下框架：

- **Gin**：高性能 Web 框架，易于使用且功能丰富，支持静态文件服务
- **Echo**：简洁高效的 Web 框架，API 设计优雅
- **Fiber**：Express 风格的 Web 框架，性能优异

**前端打包集成方案**：
- 使用 `go:embed` (Go 1.16+) 将编译后的前端资源嵌入到二进制文件中
- 使用 `statik` 或 `packr` 等工具将静态资源打包到 Go 二进制文件中

## 4. 目录结构设计

### JSON Schema 文件存储结构
```
/schemas/
  /navigation/           # 导航栏级别的 schema
    schema1.json
    schema2.json
  /sidebar/              # 侧边栏级别的 schema
    schema1.json
    schema2.json
  /page/                 # 页面级别的 schema
    schema1.json
    schema2.json
  schema-registry.json   # schema 注册表，管理所有 schema 的元数据
```

### JSON 配置文件存储结构
```
/configs/
  /navigation/           # 对应导航栏级别的配置
    config1.json
    config2.json
  /sidebar/              # 对应侧边栏级别的配置
    config1.json
    config2.json
  /page/                 # 对应页面级别的配置
    config1.json
    config2.json
  config-registry.json   # 配置注册表，管理所有配置的元数据和关系
```

### 交互关系
- JSON Schema 定义了配置文件的结构和验证规则
- 前端页面根据 Schema 动态生成表单界面
- 用户通过表单界面创建或修改配置文件
- 后端 API 负责读取 Schema 和配置文件，并提供验证和持久化服务
- 配置修改页面通过 Schema 获取字段定义，实时验证用户输入

## 5. 可测试性考虑

### 前端组件测试
- **单元测试**：使用 **Vitest** 或 **Jest** 进行 Vue 组件的单元测试
- **组件测试**：使用 **Vue Test Utils** 测试组件交互
- **端到端测试**：使用 **Cypress** 或 **Playwright** 进行 UI 自动化测试
- **模拟数据**：使用 **MSW (Mock Service Worker)** 模拟 API 响应

### 后端接口测试
- **单元测试**：使用 Go 标准库 `testing` 包
- **HTTP 测试**：使用 `httptest` 包测试 API 端点
- **集成测试**：使用 `testify` 库进行断言和模拟
- **基准测试**：使用 Go 的基准测试功能评估性能

### 测试策略
- 组件和接口分离测试
- 使用依赖注入简化测试
- 实现契约测试确保前后端接口一致性
- 自动化测试集成到 CI/CD 流程
