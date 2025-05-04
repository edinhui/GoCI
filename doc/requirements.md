# 配置管理系统需求文档

## 1. 系统概述

### 1.1 系统目标

GoCI (Golang Configuration Infrastructure) 是一个专为 Golang 应用程序设计的配置管理平台，旨在解决现代应用程序开发和运维过程中的配置管理挑战。本系统提供基于约定的配置定义、存储、访问和监控能力，帮助开发者和运维人员简化配置管理流程，提高系统可靠性和开发效率。

### 1.2 核心功能

- **配置声明与验证**：基于 OpenAPI 3.0 规范的配置项定义和自动验证
- **多环境配置管理**：支持开发、测试、生产等不同环境的配置隔离与继承
- **配置版本控制**：记录配置变更历史，支持版本对比和回滚
- **动态配置更新**：支持运行时配置热更新，无需重启服务
- **统一访问接口**：提供 CLI、Web 界面和 REST API 多种配置访问方式
- **状态监控集成**：支持应用程序运行状态的监控和可视化展示

### 1.3 技术特点

- **库集成模式**：作为库集成到 Go 应用中，无需额外部署服务
- **单一二进制**：所有功能集成在单一二进制中，简化部署和维护
- **跨平台支持**：同时支持 Linux 和 Windows 环境
- **无外部依赖**：核心功能不依赖外部服务，可在隔离环境运行
- **模块化设计**：功能模块化，支持按需集成
- **安全机制**：内置权限控制和敏感信息保护
- **可扩展架构**：预留扩展点，支持自定义存储和访问接口

### 1.4 适用场景

- **微服务应用**：管理分布式系统中的服务配置
- **多环境部署**：处理从开发到生产的配置差异
- **DevOps 集成**：支持配置自动化和 CI/CD 流程
- **遗留系统现代化**：为现有系统添加现代配置管理能力
- **边缘计算**：支持在资源受限环境中的配置管理需求

### 1.5 设计理念

GoCI 遵循"约定优于配置"的设计理念，通过预定义的配置模式和标准化的接口，降低配置管理的复杂度。系统强调可维护性和自文档化，使配置定义和使用对开发者友好且易于理解。同时，通过将配置管理功能直接集成到应用程序中，避免了额外的基础设施依赖，简化了部署流程，提高了系统的可靠性。

## 2. 系统约定

### 2.1 配置约定
- 配置项命名约定
  - 使用小写字母和下划线
  - 名称应具有描述性
  - 避免使用特殊字符
  - 示例：`network_interface`, `system_timezone`

- 配置Schema定义（基于OpenAPI 3.0规范）
  - 预定义通用类型
    ```yaml
    # common_types.yaml
    openapi: 3.0.0
    info:
      title: Common Types Definition
      version: 1.0.0
      description: 系统预定义的通用类型
    components:
      schemas:
        # 基础类型
        String:
          type: string
          description: 字符串类型
        Integer:
          type: integer
          description: 整数类型
        Float:
          type: number
          description: 浮点数类型
        Boolean:
          type: boolean
          description: 布尔类型
        DateTime:
          type: string
          format: date-time
          description: 日期时间类型
        
        # 网络相关
        IPv4Address:
          type: string
          format: ipv4
          description: IPv4地址
        IPv6Address:
          type: string
          format: ipv6
          description: IPv6地址
        MACAddress:
          type: string
          pattern: "^([0-9A-Fa-f]{2}[:-]){5}([0-9A-Fa-f]{2})$"
          description: MAC地址
        Port:
          type: integer
          minimum: 0
          maximum: 65535
          description: 端口号
        
        # 文件系统
        FilePath:
          type: string
          pattern: "^[a-zA-Z0-9_./-]+$"
          description: 文件路径
        DirectoryPath:
          type: string
          pattern: "^[a-zA-Z0-9_./-]+$"
          description: 目录路径
        
        # 安全相关
        Password:
          type: string
          minLength: 8
          pattern: "^(?=.*[A-Za-z])(?=.*\d)[A-Za-z\d]{8,}$"
          description: 密码（至少8位，包含字母和数字）
        Email:
          type: string
          format: email
          description: 电子邮件地址
    ```

  - 配置Schema组织方式
    - 单一文件方式（推荐用于简单配置）
      ```yaml
      # app_config.yaml
      openapi: 3.0.0
      info:
        title: Application Configuration
        version: 1.0.0
        description: 应用程序配置定义
      tags:
        - name: system
          description: 系统配置
        - name: network
          description: 网络配置
      components:
        schemas:
          AppConfig:
            type: object
            description: 应用程序配置
            properties:
              server:
                $ref: './common_types.yaml#/components/schemas/ServerConfig'
                description: 服务器配置
              database:
                $ref: './common_types.yaml#/components/schemas/DatabaseConfig'
                description: 数据库配置
      ```

    - 多文件方式（推荐用于复杂配置）
      ```yaml
      # config.yaml
      openapi: 3.0.0
      info:
        title: Application Configuration
        version: 1.0.0
        description: 应用程序配置定义
      tags:
        - name: system
          description: 系统配置
        - name: network
          description: 网络配置
        - name: security
          description: 安全配置
      components:
        schemas:
          AppConfig:
            type: object
            description: 应用程序配置
            properties:
              server:
                $ref: './schemas/server.yaml#/ServerConfig'
                description: 服务器配置
              database:
                $ref: './schemas/database.yaml#/DatabaseConfig'
                description: 数据库配置
              logging:
                $ref: './schemas/logging.yaml#/LoggingConfig'
                description: 日志配置
      ```

  - 配置Schema编写建议
    - 简单配置（<10个配置项）
      - 使用单一文件
      - 直接引用通用类型
      - 保持结构扁平

    - 中等配置（10-50个配置项）
      - 考虑按功能模块拆分
      - 使用通用类型组合
      - 保持合理的嵌套层级

    - 复杂配置（>50个配置项）
      - 必须按功能模块拆分
      - 建立清晰的目录结构
      - 使用多级引用
      - 添加详细的注释说明

  - 目录结构约定
    ```
    config/
    ├── common_types.yaml      # 系统预定义通用类型
    ├── app_config.yaml        # 主配置文件（简单配置）
    └── schemas/               # 复杂配置的Schema文件
        ├── server.yaml
        ├── database.yaml
        ├── logging.yaml
        └── security.yaml
    ```

### 2.2 前端展示约定
> 注：本节所述为系统推荐的前端展示约定，适用于默认前端实现。系统支持用户根据实际需求开发自定义前端界面，采用不同的布局、分组和交互方式。只需遵循后端接口规范，即可实现多样化的前端展示。

- 配置项编辑器选择
  - 字符串：文本输入框
  - 数字：数字输入框
  - 布尔值：开关/复选框
  - 枚举值：下拉选择框
  - 数组：可扩展列表
  - 对象：折叠面板

- 分组展示约定
  - 一级分组：顶部导航
  - 二级分组：侧边栏
  - 三级分组：标签页

- 表单布局约定
  - 标签左对齐
  - 输入框右对齐
  - 必填项标记
  - 验证错误提示

- 响应式设计约定
  - 移动端：单列布局
  - 桌面端：多列布局
  - 自适应间距

### 2.3 权限约定
- 基于分组的权限控制
- 权限级别：
  - `read`: 只读
  - `write`: 可修改
  - `admin`: 完全控制
- 权限继承规则

## 3. 功能需求

### 3.1 配置管理功能

#### 3.1.1 核心配置管理
- 配置项定义
  - 支持通过Schema文件定义配置项的名称、类型、默认值、描述、校验规则等
  - 支持分组、嵌套、枚举、数组等复杂结构
  - 支持通用类型引用
  - 支持配置项启用/禁用状态管理
  - 日志配置Schema定义示例
    ```yaml
    LoggingConfig:
      type: object
      description: 日志配置项
      properties:
        default_level:
          type: string
          enum: ["debug", "info", "warn", "error", "fatal"]
          default: "info"
          description: 默认日志级别
        modules:
          type: object
          additionalProperties:
            type: object
            properties:
              level:
                type: string
                enum: ["debug", "info", "warn", "error", "fatal"]
                description: 模块日志级别
              output:
                oneOf:
                  - type: string
                    enum: ["console", "file", "syslog"]
                  - type: array
                    items:
                      type: string
                      enum: ["console", "file", "syslog"]
                default: "console"
                description: 日志输出目标
              file_path:
                type: string
                description: 日志文件路径，仅在output包含file时有效
            required: ["level"]
          description: 按模块配置日志
        rotation:
          type: object
          properties:
            max_size:
              type: integer
              minimum: 1
              default: 100
              description: 单个日志文件最大大小(MB)
            max_age:
              type: integer
              minimum: 1
              default: 7
              description: 日志文件保留天数
            max_backups:
              type: integer
              minimum: 0
              default: 10
              description: 保留的旧日志文件数量
            compress:
              type: boolean
              default: true
              description: 是否压缩旧日志文件
          description: 日志轮转设置
        format:
          type: string
          enum: ["text", "json"]
          default: "text"
          description: 日志输出格式
        timestamp_format:
          type: string
          enum: ["RFC3339", "RFC3339Nano", "Unix", "UnixMs"]
          default: "RFC3339"
          description: 日志时间戳格式
      required: ["default_level"]
    ```

- 配置项操作
  - 支持配置项的增删改查（CRUD）
  - 支持批量操作
  - 支持配置项导入/导出（YAML/JSON格式）
  - 支持导入校验与冲突处理

- 配置项校验
  - 自动校验类型、范围、必填性
  - 支持自定义校验规则
  - 校验失败详细提示

#### 3.1.2 配置版本与变更管理
- 版本控制
  - 支持历史版本记录、对比、回滚
  - 支持变更说明与标签
  - 支持一键回滚到历史版本

- 变更审计
  - 记录变更操作、操作人、时间、内容
  - 支持历史查询与导出
  - 支持变更影响分析

#### 3.1.3 配置可视化工具
- 可视化配置项定义生成器
  - 提供Web页面可视化编辑器，支持图形化创建和编辑配置项定义
  - 支持字段类型、描述、默认值、校验规则等属性的可视化配置
  - 支持分组、嵌套、枚举、数组等复杂结构的可视化编辑
  - 支持通用类型引用与模板复用
  - 实时预览和校验Schema结构

- Excel格式导入导出
  - 支持将配置项（包括对象、数组等复杂结构）导出为Excel文件
  - 支持多Sheet导出，按配置分组或类型分Sheet
  - 支持导出配置项元信息和当前配置值
  - 支持从Excel文件导入配置项定义和配置值
  - 提供Excel模板功能

#### 3.1.4 配置依赖与权限管理
- 配置项引用与依赖
  - 支持$ref引用
  - 支持依赖关系校验与可视化
  - 支持依赖变更影响分析

- 配置项权限管理
  - 支持细粒度权限设置与继承
  - 支持基于分组的权限控制
  - 支持权限变更审计

#### 3.1.5 配置分发与生效
- 配置动态加载与推送
  - 支持变更后自动推送
  - 支持服务端/客户端动态加载
  - 支持配置变更通知
  - 日志配置动态更新示例
    ```yaml
    logging:
      default_level: "info"      # 默认日志级别
      modules:                   # 模块特定配置
        config_manager:
          level: "info"
          output: "file"
          file_path: "logs/config.log"
        web_server:
          level: "warn"
          output: "console"
        api:
          level: "debug"
          output: ["file", "console"]
          file_path: "logs/api.log"
      rotation:                  # 日志轮转配置
        max_size: 100            # 单位：MB
        max_age: 7               # 单位：天
        max_backups: 10          # 保留文件数
        compress: true           # 是否压缩
      format: "json"             # 日志格式: text/json
      timestamp_format: "RFC3339" # 时间戳格式
    ```

- 配置生效控制
  - 支持立即或延迟生效
  - 支持配置变更预览
  - 支持配置回滚

### 3.2 访问接口

#### 3.2.1 命令行接口（CLI）
- 配置查询
  - 支持查询单个配置项：`goci config get <path>`
  - 支持查询配置组：`goci config list <group>`
  - 支持模糊搜索：`goci config search <keyword>`
  - 支持查看配置项描述：`goci config describe <path>`
  - 支持查看配置历史：`goci config history <path>`

- 配置修改
  - 支持设置配置值：`goci config set <path> <value>`
  - 支持批量设置：`goci config set-batch <file>`
  - 支持删除配置：`goci config delete <path>`
  - 支持配置回滚：`goci config rollback <version>`

- 配置导入导出
  - 支持导出配置：`goci config export <format> <path>`
  - 支持导入配置：`goci config import <format> <file>`
  - 支持导出Schema：`goci schema export <path>`
  - 支持导入Schema：`goci schema import <file>`

- 系统管理
  - 支持查看系统状态：`goci status`
  - 支持查看配置变更：`goci changes`
  - 支持查看依赖关系：`goci dependencies <path>`
  - 支持查看权限信息：`goci permissions <path>`

- 帮助与文档
  - 支持查看命令帮助：`goci help`
  - 支持查看配置说明：`goci docs <path>`
  - 支持查看示例：`goci examples`

#### 3.2.2 Web界面
- 提供图形化配置管理界面
- 支持配置项的可视化编辑
- 支持配置变更的实时预览
- 支持配置历史查看与回滚
- 支持权限管理与审计日志查看

#### 3.2.3 REST API
- 提供标准的RESTful API接口
- 支持配置的CRUD操作
- 支持批量操作
- 支持版本控制
- 支持权限验证

### 3.3 状态监控

#### 3.3.1 自动状态展示
- 基于OpenAPI规范自动生成状态展示页面
  - 支持从OpenAPI文档自动识别状态监控接口
  - 自动生成状态指标展示页面
  - 支持自定义展示模板
  - 支持实时数据刷新

- 状态指标定义
  - 支持在OpenAPI中定义状态指标
    ```yaml
    paths:
      /api/v1/status:
        get:
          summary: 获取系统状态
          tags: [status]
          responses:
            200:
              description: 系统状态信息
              content:
                application/json:
                  schema:
                    type: object
                    properties:
                      cpu_usage:
                        type: number
                        format: float
                        description: CPU使用率
                        x-status-display:
                          type: gauge
                          unit: "%"
                          warning: 80
                          critical: 90
                      memory_usage:
                        type: number
                        format: float
                        description: 内存使用率
                        x-status-display:
                          type: gauge
                          unit: "%"
                          warning: 85
                          critical: 95
                      active_connections:
                        type: integer
                        description: 活动连接数
                        x-status-display:
                          type: counter
                          unit: "个"
    ```

#### 3.3.2 状态展示组件
- 支持多种展示组件
  - 仪表盘（Gauge）
  - 计数器（Counter）
  - 趋势图（Trend）
  - 状态灯（Status Light）
  - 表格（Table）
  - 列表（List）

- 展示组件配置
  - 支持设置阈值和告警级别
  - 支持自定义刷新间隔
  - 支持自定义样式
  - 支持组件布局调整

#### 3.3.3 状态数据管理
- 数据采集
  - 支持定时采集
  - 支持事件触发采集
  - 支持自定义采集间隔

- 数据存储
  - 支持实时数据缓存
  - 支持历史数据存储
  - 支持数据聚合

- 数据展示
  - 支持实时数据展示
  - 支持历史数据查询
  - 支持数据对比分析

#### 3.3.4 告警与通知
- 告警规则
  - 支持阈值告警
  - 支持变化率告警
  - 支持组合条件告警

- 告警通知
  - 支持邮件通知
  - 支持Webhook通知
  - 支持自定义通知方式

#### 3.3.5 集成与扩展
- 支持自定义状态指标
- 支持自定义展示组件
- 支持自定义数据源
- 支持第三方系统集成

### 3.4 业务程序集成

#### 3.4.1 库集成模式（推荐方式）
- 提供配置管理库（goci）
  - 支持通过内嵌方式集成到业务程序
  - 使业务程序仅需提供单一二进制文件
  - 无需额外部署配置管理服务
  - 配置管理功能与业务程序一体化
  
- 核心功能
  - 支持本地配置文件解析与校验
  - 支持配置变更的实时监听
  - 支持配置的热加载
  - 支持配置历史版本记录与回滚
  - 提供内置状态页面，无需额外服务

- 实例标识与管理
  - 进程隔离与识别机制（跨平台实现）
    ```go
    // 实例管理器，用于识别和隔离不同实例
    type InstanceManager struct {
        ID        string   // 实例唯一标识
        PidFile   string   // PID文件路径
        LockFile  *os.File // 文件锁
        ConfigDir string   // 实例配置目录
    }
    
    // 获取系统配置目录（跨平台）
    func getSystemConfigDir() string {
        if dirPath := os.Getenv("GOCI_CONFIG_DIR"); dirPath != "" {
            return dirPath
        }
        
        // 使用用户配置目录作为基础路径（跨平台）
        configDir, err := os.UserConfigDir()
        if err != nil {
            // 回退方案：使用临时目录
            return filepath.Join(os.TempDir(), "goci")
        }
        return filepath.Join(configDir, "goci")
    }
    
    // 获取PID文件目录（跨平台）
    func getSystemPidDir() string {
        if dirPath := os.Getenv("GOCI_PID_DIR"); dirPath != "" {
            return dirPath
        }
        
        // 平台特定处理
        switch runtime.GOOS {
        case "windows":
            // Windows使用ProgramData目录或临时目录
            if dataDir := os.Getenv("PROGRAMDATA"); dataDir != "" {
                return filepath.Join(dataDir, "goci", "pids")
            }
            return filepath.Join(os.TempDir(), "goci", "pids")
        default:
            // Linux/MacOS使用/var/run或/tmp
            if _, err := os.Stat("/var/run"); err == nil {
                return "/var/run"
            }
            return filepath.Join(os.TempDir(), "goci", "pids")
        }
    }
    
    // 创建新实例管理器（跨平台）
    func NewInstanceManager(id string) (*InstanceManager, error) {
        if id == "" {
            // 默认使用主机名作为实例ID
            hostname, err := os.Hostname()
            if err != nil {
                return nil, err
            }
            id = hostname
        }
        
        // 确保PID目录存在
        pidDir := getSystemPidDir()
        if err := os.MkdirAll(pidDir, 0755); err != nil {
            return nil, fmt.Errorf("无法创建PID目录: %v", err)
        }
        
        pidFile := filepath.Join(pidDir, fmt.Sprintf("goci_%s.pid", id))
        
        // 检查PID文件是否存在
        if _, err := os.Stat(pidFile); err == nil {
            // 检查进程是否存在（跨平台）
            if isProcessRunning(pidFile) {
                return nil, fmt.Errorf("实例 %s 已在运行中", id)
            }
            // PID文件存在但进程不存在，可以覆盖
            os.Remove(pidFile)
        }
        
        // 创建并锁定PID文件（跨平台）
        lockFile, err := createLockedPidFile(pidFile)
        if err != nil {
            return nil, fmt.Errorf("无法创建或锁定PID文件: %v", err)
        }
        
        // 写入当前PID
        pid := os.Getpid()
        if _, err := lockFile.WriteString(fmt.Sprintf("%d\n", pid)); err != nil {
            releaseLock(lockFile)
            os.Remove(pidFile)
            return nil, err
        }
        
        // 创建配置目录
        configDir := filepath.Join(getSystemConfigDir(), id)
        if err := os.MkdirAll(configDir, 0755); err != nil {
            releaseLock(lockFile)
            os.Remove(pidFile)
            return nil, fmt.Errorf("无法创建配置目录: %v", err)
        }
        
        return &InstanceManager{
            ID:        id,
            PidFile:   pidFile,
            LockFile:  lockFile,
            ConfigDir: configDir,
        }, nil
    }
    
    // 释放实例锁并清理资源（跨平台）
    func (im *InstanceManager) Release() {
        if im.LockFile != nil {
            releaseLock(im.LockFile)
            os.Remove(im.PidFile)
        }
    }
    
    //
    // 平台特定实现（使用构建标签分离）
    //
    
    // +build windows
    
    package goci
    
    import (
        "os"
        "strconv"
        "strings"
        "syscall"
        "unsafe"
    )
    
    // Windows版本的文件锁实现
    func createLockedPidFile(path string) (*os.File, error) {
        // Windows使用创建时独占打开文件的方式实现文件锁
        f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_EXCL, 0644)
        if err != nil {
            // 如果文件已存在但无法独占打开，尝试强制解锁
            if os.IsExist(err) {
                // 尝试删除可能的僵尸文件
                if err := os.Remove(path); err != nil {
                    return nil, err
                }
                // 重新尝试创建
                return os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_EXCL, 0644)
            }
            return nil, err
        }
        return f, nil
    }
    
    // Windows版本的文件锁释放
    func releaseLock(file *os.File) {
        if file != nil {
            file.Close()
        }
    }
    
    // Windows版本的进程检查
    func isProcessRunning(pidFile string) bool {
        data, err := os.ReadFile(pidFile)
        if err != nil {
            return false
        }
        
        pidStr := strings.TrimSpace(string(data))
        pid, err := strconv.Atoi(pidStr)
        if err != nil {
            return false
        }
        
        // Windows需要使用OpenProcess API
        kernel32 := syscall.NewLazyDLL("kernel32.dll")
        OpenProcess := kernel32.NewProc("OpenProcess")
        CloseHandle := kernel32.NewProc("CloseHandle")
        
        const PROCESS_QUERY_INFORMATION = 0x0400
        
        handle, _, _ := OpenProcess.Call(
            uintptr(PROCESS_QUERY_INFORMATION),
            uintptr(0),
            uintptr(pid),
        )
        
        if handle == 0 {
            return false
        }
        
        CloseHandle.Call(handle)
        return true
    }
    
    // +build !windows
    
    package goci
    
    import (
        "os"
        "strconv"
        "strings"
        "syscall"
    )
    
    // Unix/Linux版本的文件锁实现
    func createLockedPidFile(path string) (*os.File, error) {
        lockFile, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0644)
        if err != nil {
            return nil, err
        }
        
        // 尝试获取文件锁
        if err := syscall.Flock(int(lockFile.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); err != nil {
            lockFile.Close()
            return nil, err
        }
        
        return lockFile, nil
    }
    
    // Unix/Linux版本的文件锁释放
    func releaseLock(file *os.File) {
        if file != nil {
            syscall.Flock(int(file.Fd()), syscall.LOCK_UN)
            file.Close()
        }
    }
    
    // Unix/Linux版本的进程检查
    func isProcessRunning(pidFile string) bool {
        data, err := os.ReadFile(pidFile)
        if err != nil {
            return false
        }
        
        pidStr := strings.TrimSpace(string(data))
        pid, err := strconv.Atoi(pidStr)
        if err != nil {
            return false
        }
        
        // 使用Signal(0)检查进程是否存在
        process, err := os.FindProcess(pid)
        if err != nil {
            return false
        }
        
        return process.Signal(syscall.Signal(0)) == nil
    }
    ```
  
  - 实例发现与通信（跨平台）
    ```go
    // 列出所有运行中的实例（跨平台）
    func listInstances() ([]InstanceInfo, error) {
        // 获取PID目录（跨平台）
        pidDir := getSystemPidDir()
        
        // 确保目录存在
        if _, err := os.Stat(pidDir); os.IsNotExist(err) {
            return []InstanceInfo{}, nil // 目录不存在时返回空列表
        }
        
        // 构建模式匹配（跨平台）
        pattern := filepath.Join(pidDir, "goci_*.pid")
        
        // 获取匹配文件
        files, err := filepath.Glob(pattern)
        if err != nil {
            return nil, err
        }
        
        instances := make([]InstanceInfo, 0, len(files))
        for _, file := range files {
            // 从文件名提取实例ID
            filename := filepath.Base(file)
            id := strings.TrimPrefix(filename, "goci_")
            id = strings.TrimSuffix(id, ".pid")
            
            // 读取PID
            data, err := os.ReadFile(file)
            if err != nil {
                continue
            }
            
            pidStr := strings.TrimSpace(string(data))
            pid, err := strconv.Atoi(pidStr)
            if err != nil {
                continue
            }
            
            // 验证进程是否存在（跨平台）
            if !isProcessRunning(file) {
                // 进程不存在，清理僵尸PID文件
                os.Remove(file)
                continue
            }
            
            // 获取配置目录（跨平台）
            configDir := filepath.Join(getSystemConfigDir(), id)
            
            instances = append(instances, InstanceInfo{
                ID:  id,
                PID: pid,
                ConfigDir: configDir,
            })
        }
        
        return instances, nil
    }
    
    // 实例信息结构
    type InstanceInfo struct {
        ID        string // 实例ID
        PID       int    // 进程ID
        ConfigDir string // 配置目录
    }
    ```
  
  - 多实例配置隔离（跨平台）
    ```go
    // 获取实例特定的配置路径（跨平台）
    func (im *InstanceManager) GetConfigPath(fileName string) string {
        return filepath.Join(im.ConfigDir, fileName)
    }
    
    // 加载实例特定的配置（跨平台）
    func (im *InstanceManager) LoadConfig() (*Config, error) {
        configPath := im.GetConfigPath("config.yaml")
        
        // 检查实例特定配置是否存在，不存在则复制全局配置
        if _, err := os.Stat(configPath); os.IsNotExist(err) {
            // 尝试多个位置查找全局配置
            globalConfigPaths := []string{
                filepath.Join(getSystemConfigDir(), "config", "config.yaml"),
                "/etc/goci/config/config.yaml",
                filepath.Join(os.Getenv("PROGRAMDATA"), "goci", "config", "config.yaml"), // Windows
            }
            
            for _, path := range globalConfigPaths {
                if _, err := os.Stat(path); err == nil {
                    data, err := os.ReadFile(path)
                    if err == nil {
                        // 创建目录（如果需要）
                        os.MkdirAll(filepath.Dir(configPath), 0755)
                        
                        // 复制配置到实例特定路径
                        if err := os.WriteFile(configPath, data, 0644); err == nil {
                            break
                        }
                    }
                }
            }
        }
        
        // 加载配置（如果存在）
        if _, err := os.Stat(configPath); os.IsNotExist(err) {
            // 创建默认配置
            defaultConfig := &Config{
                // 添加默认值
            }
            
            // 序列化为YAML
            data, err := yaml.Marshal(defaultConfig)
            if err != nil {
                return nil, err
            }
            
            // 创建目录（如果需要）
            os.MkdirAll(filepath.Dir(configPath), 0755)
            
            // 写入默认配置
            if err := os.WriteFile(configPath, data, 0644); err != nil {
                return nil, err
            }
            
            return defaultConfig, nil
        }
        
        // 读取配置文件
        data, err := os.ReadFile(configPath)
        if err != nil {
            return nil, err
        }
        
        // 解析YAML
        cfg := &Config{}
        if err := yaml.Unmarshal(data, cfg); err != nil {
            return nil, err
        }
        
        return cfg, nil
    }
    
    // 保存实例配置（原子写入，跨平台）
    func (im *InstanceManager) SaveConfig(config *Config) error {
        configPath := im.GetConfigPath("config.yaml")
        
        // 序列化为YAML
        data, err := yaml.Marshal(config)
        if err != nil {
            return err
        }
        
        // 创建临时文件
        tempFile := configPath + ".tmp"
        if err := os.WriteFile(tempFile, data, 0644); err != nil {
            return err
        }
        
        // 原子重命名（跨平台安全写入）
        if err := os.Rename(tempFile, configPath); err != nil {
            os.Remove(tempFile) // 清理临时文件
            return err
        }
        
        return nil
    }
    ```

- 集成示例
  ```go
  // 1. 添加依赖
  // go get github.com/your-org/goci
  
  // 2. 在main函数中初始化配置
  func main() {
      // 创建实例管理器，指定实例ID或使用默认主机名
      instanceID := os.Getenv("GOCI_INSTANCE")
      instance, err := goci.NewInstanceManager(instanceID)
      if err != nil {
          log.Fatalf("实例初始化失败: %v", err)
      }
      defer instance.Release() // 确保退出时释放锁并清理资源
      
      // 初始化配置管理器
      config := goci.NewManager(goci.Options{
          SchemaPath: instance.GetConfigPath("app_config.yaml"),  // 实例特定Schema定义
          ConfigPath: instance.GetConfigPath("config.yaml"),      // 实例特定配置文件路径
          EnableWeb: true,                                        // 启用Web管理界面
          WebPort: 8080,                                          // Web管理端口
          Instance: instance,                                     // 传递实例管理器
      })
      
      // 加载配置
      if err := config.Load(); err != nil {
          log.Fatalf("加载配置失败: %v", err)
      }
      
      // 初始化日志系统
      if err := initLogger(config); err != nil {
          log.Printf("警告: 初始化日志系统失败: %v，将使用默认配置", err)
      }
      
      // 获取logger实例并使用
      logger := goci.GetLogger("main")
      logger.Info("应用启动成功", map[string]interface{}{
          "instance_id": instanceID,
          "web_port": config.GetInt("web_port"),
          "version": "1.0.0",
      })
      
      // 注册配置变更回调
      config.OnChange(func(newConfig *goci.Config) {
          // 处理配置变更
          logger.Info("配置已更新")
          
          // 如果日志配置变更，重新加载日志系统
          if newConfig.HasChanged("logging") {
              logger.Debug("检测到日志配置变更，正在重新加载")
              if err := reloadLogger(newConfig); err != nil {
                  logger.Error("重新加载日志配置失败", map[string]interface{}{
                      "error": err.Error(),
                  })
              } else {
                  logger.Info("日志配置已重新加载")
              }
          }
      })
      
      // 启动业务服务
      startApplication(config)
  }
  
  // 初始化日志系统
  func initLogger(config *goci.Config) error {
      // 获取日志配置
      logConfig, err := config.GetObject("logging")
      if err != nil {
          return err
      }
      
      // 创建日志管理器
      logManager := goci.NewLogManager(logConfig)
      
      // 注册到全局
      goci.SetLogManager(logManager)
      
      return nil
  }
  
  // 重新加载日志配置
  func reloadLogger(config *goci.Config) error {
      logConfig, err := config.GetObject("logging")
      if err != nil {
          return err
      }
      
      // 获取日志管理器并重新加载配置
      logManager := goci.GetLogManager()
      return logManager.Reload(logConfig)
  }
  ```

#### 3.4.2 单一二进制部署优势
- 业务程序与配置管理集成为单一二进制
  - 简化部署流程，无需额外服务
  - 降低运维复杂度
  - 减少依赖和潜在故障点
  
- 内置Web管理界面
  - 业务二进制内置配置管理Web界面
  - 通过配置开关和端口控制Web界面
  - 提供简化版配置管理功能，满足基本需求

### 3.5 数据存储与备份

#### 3.5.1 文件存储
- 存储格式与组织
  - 使用YAML/JSON格式存储配置数据
  - 按照功能模块划分配置文件
  - 支持包含和引用机制
  - 存储路径约定（默认位于应用同目录的`config/`文件夹）

- 文件存储优势
  - 简单直观，易于手动编辑和检查
  - 无需额外数据库依赖
  - 方便与版本控制系统（如Git）集成
  - 配置文件可纳入代码审核流程
  - 与单一二进制部署模式天然契合

- 文件操作安全性
  - 支持文件操作原子性（通过临时文件和重命名实现）
  - 支持文件锁定，防止并发写入冲突
  - 文件修改前自动备份
  - 定期检查文件完整性

- 版本控制实现
  - 每次修改生成带时间戳的备份文件（例如：`config.yaml.20230624120000`）
  - 保留可配置数量的历史版本
  - 支持差异对比和还原指定版本

#### 3.5.2 备份与恢复
- 自动备份策略
  - 每次修改前自动创建备份
  - 定时创建备份（可配置频率）
  - 可配置备份保留数量或时间

- 备份存储位置
  - 本地备份：默认位于`config/backups/`目录
  - 可选配置远程备份（如S3、SFTP等）

- 恢复机制
  - 支持一键回滚到任意历史版本
  - 支持指定时间点恢复
  - 恢复前自动备份当前配置（安全保障）

### 3.6 多环境支持

#### 3.6.1 环境隔离机制
- 配置隔离
  - 支持不同环境（开发、测试、生产等）的配置隔离
  - 支持配置继承关系（如基础配置+环境特定配置）
  - 提供环境配置模板，方便快速创建新环境
  - 支持跨平台（Windows/Linux）环境一致性

- 多实例支持（同机部署）
  - 支持在同一台机器上运行多个配置管理支持的业务程序
  - 每个实例拥有独立的配置空间和存储路径
  - 支持实例间的配置共享与隔离
  - 同时支持Windows和Linux平台的多实例管理

#### 3.6.2 CLI交互增强
- 实例识别与选择
  - 支持通过实例ID/名称指定目标实例：`goci --instance <name> <command>`
  - 支持列出当前机器上所有实例：`goci instance list`
  - 提供实例环境变量：`GOCI_INSTANCE=<name> goci <command>`

- 实例识别机制
  - 基于PID文件实现实例识别和互斥运行
    - 每个实例在启动时创建唯一PID文件
      - Linux：`/var/run/goci_<instance_name>.pid` 或 `/tmp/goci/pids/goci_<instance_name>.pid` 
      - Windows：`%PROGRAMDATA%\goci\pids\goci_<instance_name>.pid` 或 `%TEMP%\goci\pids\goci_<instance_name>.pid`
    - PID文件包含进程ID和实例标识信息
    - 自动检测"僵尸"PID文件并清理
  - 实例名称作为唯一标识符
    - 默认使用主机名作为实例名称
    - 支持通过配置文件或环境变量自定义实例名称：`GOCI_INSTANCE=<name>`
    - 支持通过CLI参数临时指定实例名称：`--instance-name <name>`
  - 跨平台环境变量支持
    - 支持自定义PID文件目录：`GOCI_PID_DIR=<directory>`
    - 支持自定义配置目录：`GOCI_CONFIG_DIR=<directory>`

- 实例管理指令
  - 支持查看实例详情：`goci instance info <name>`
  - 支持检查实例状态：`goci instance status <name>`
  - 支持实例锁定/解锁：`goci instance lock/unlock <name>`

- 智能CLI提示
  - 当存在多个实例时，命令失败会提示可用实例列表
  - 在实例上下文中自动完成配置路径
  - 支持记住常用实例（通过配置文件）

#### 3.6.3 Web端口管理
- 动态端口分配
  - 支持自动检测可用端口并分配
  - 支持优先使用指定端口，若被占用则自动增量查找
  - 启动时清晰显示已分配端口：`Config UI available at http://localhost:8080`

- 端口配置方式
  - 命令行参数指定：`--web-port <port>`
  - 配置文件指定：`web_port: <port>`
  - 环境变量指定：`GOCI_WEB_PORT=<port>`
  - 支持禁用Web界面：`--web-disable`

- 端口冲突处理
  - 端口被占用时提供明确错误信息
  - 提供自动端口搜索选项：`--web-port-auto`

- 实现示例
  ```go
  // 端口管理器
  type PortManager struct {
      BasePort     int    // 基础端口号
      MaxPortTries int    // 最大尝试次数
      InstanceID   string // 实例ID
  }
  
  // 创建新的端口管理器
  func NewPortManager(basePort int, instanceID string) *PortManager {
      return &PortManager{
          BasePort:     basePort,
          MaxPortTries: 100, // 尝试100个端口
          InstanceID:   instanceID,
      }
  }
  
  // 寻找可用端口
  func (pm *PortManager) FindAvailablePort() (int, error) {
      // 首先尝试使用基础端口
      port := pm.BasePort
      
      // 最多尝试MaxPortTries次
      for i := 0; i < pm.MaxPortTries; i++ {
          // 尝试绑定端口
          listener, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
          if err == nil {
              // 端口可用，关闭监听并返回端口号
              listener.Close()
              return port, nil
          }
          
          // 端口被占用，尝试下一个端口
          port++
      }
      
      return 0, fmt.Errorf("无法找到可用端口，已尝试从 %d 到 %d", pm.BasePort, pm.BasePort+pm.MaxPortTries-1)
  }
  
  // 保存端口信息到实例元数据
  func (pm *PortManager) SavePortInfo(port int) error {
      // 端口信息文件路径（跨平台）
      var portFile string
      if runtime.GOOS == "windows" {
          // Windows平台
          pidDir := getSystemPidDir() // 复用之前定义的获取PID目录函数
          portFile = filepath.Join(pidDir, fmt.Sprintf("goci_%s.port", pm.InstanceID))
      } else {
          // Linux/Unix平台
          if _, err := os.Stat("/var/run"); err == nil {
              portFile = fmt.Sprintf("/var/run/goci_%s.port", pm.InstanceID)
          } else {
              portFile = filepath.Join(os.TempDir(), "goci", "pids", fmt.Sprintf("goci_%s.port", pm.InstanceID))
          }
      }
      
      // 确保目录存在
      if err := os.MkdirAll(filepath.Dir(portFile), 0755); err != nil {
          return err
      }
      
      // 写入端口信息
      return os.WriteFile(portFile, []byte(fmt.Sprintf("%d", port)), 0644)
  }
  
  // 获取实例端口信息
  func (pm *PortManager) GetPortInfo() (int, error) {
      // 端口信息文件路径（跨平台）
      var portFile string
      if runtime.GOOS == "windows" {
          // Windows平台
          pidDir := getSystemPidDir() // 复用之前定义的获取PID目录函数
          portFile = filepath.Join(pidDir, fmt.Sprintf("goci_%s.port", pm.InstanceID))
      } else {
          // Linux/Unix平台
          if _, err := os.Stat("/var/run"); err == nil {
              portFile = fmt.Sprintf("/var/run/goci_%s.port", pm.InstanceID)
          } else {
              portFile = filepath.Join(os.TempDir(), "goci", "pids", fmt.Sprintf("goci_%s.port", pm.InstanceID))
          }
      }
      
      // 读取端口信息
      data, err := os.ReadFile(portFile)
      if err != nil {
          return 0, err
      }
      
      // 解析端口
      port, err := strconv.Atoi(strings.TrimSpace(string(data)))
      if err != nil {
          return 0, err
      }
      
      return port, nil
  }
  
  // 配置Web服务器
  func ConfigureWebServer(opts *Options, instanceManager *InstanceManager) error {
      if opts.DisableWeb {
          return nil // Web界面已禁用
      }
      
      // 创建端口管理器
      portManager := NewPortManager(opts.WebPort, instanceManager.ID)
      
      // 查找可用端口
      port, err := portManager.FindAvailablePort()
      if err != nil {
          return err
      }
      
      // 更新端口配置
      opts.WebPort = port
      
      // 保存端口信息
      if err := portManager.SavePortInfo(port); err != nil {
          return err
      }
      
      // 打印Web界面地址
      log.Printf("配置UI可访问地址: http://localhost:%d", port)
      
      return nil
  }
  ```

- 使用示例
  ```go
  func main() {
      // 初始化实例管理器
      instance, err := goci.NewInstanceManager(os.Getenv("GOCI_INSTANCE"))
      if err != nil {
          log.Fatalf("实例初始化失败: %v", err)
      }
      defer instance.Release()
      
      // 解析配置选项
      options := goci.NewOptions()
      
      // 设置Web端口
      webPort := 8080 // 默认端口
      if portStr := os.Getenv("GOCI_WEB_PORT"); portStr != "" {
          if port, err := strconv.Atoi(portStr); err == nil {
              webPort = port
          }
      }
      options.WebPort = webPort
      
      // 配置Web服务器，自动处理端口冲突
      if err := goci.ConfigureWebServer(options, instance); err != nil {
          log.Printf("Web界面配置失败: %v", err)
          options.DisableWeb = true
      }
      
      // 初始化配置管理器
      config := goci.NewManager(options)
      
      // ... 其他初始化代码 ...
  }
  ```

#### 3.6.4 实例CLI交互示例

以下是CLI命令实现多实例交互的具体示例，包括命令解析、实例识别和交互处理的代码：

```go
package cli

import (
    "fmt"
    "os"
    "strings"
    "github.com/spf13/cobra"
)

// CLI配置
type CliConfig struct {
    InstanceName  string // 当前实例名称
    GlobalFlags   bool   // 是否使用全局标志
    Interactive   bool   // 是否交互模式
    AvailableInstances []string // 可用实例列表
}

// 创建主命令
func NewRootCommand() *cobra.Command {
    config := &CliConfig{}
    
    // 获取环境变量中的实例名称
    if instanceName := os.Getenv("GOCI_INSTANCE"); instanceName != "" {
        config.InstanceName = instanceName
    }
    
    // 根命令
    rootCmd := &cobra.Command{
        Use:   "goci",
        Short: "配置管理工具",
        Long:  "Golang配置管理工具，支持多实例并发运行",
        PersistentPreRun: func(cmd *cobra.Command, args []string) {
            // 在命令执行前识别和验证实例
            if err := validateInstance(config); err != nil {
                fmt.Fprintf(os.Stderr, "错误: %v\n", err)
                // 如果未找到实例，提示可用实例
                listAvailableInstances(config)
                os.Exit(1)
            }
        },
    }
    
    // 全局实例标志
    rootCmd.PersistentFlags().StringVar(&config.InstanceName, "instance", "", 
        "指定目标实例名称")
    
    // 添加实例命令组
    rootCmd.AddCommand(newInstanceCommand(config))
    
    // 添加配置命令组
    rootCmd.AddCommand(newConfigCommand(config))
    
    return rootCmd
}

// 验证指定的实例
func validateInstance(config *CliConfig) error {
    // 如果未指定实例名称，尝试使用默认实例
    if config.InstanceName == "" {
        hostname, err := os.Hostname()
        if err != nil {
            return fmt.Errorf("无法获取主机名作为默认实例: %v", err)
        }
        config.InstanceName = hostname
    }
    
    // 检查实例是否存在并运行中
    instances, err := listInstances()
    if err != nil {
        return err
    }
    
    // 保存可用实例列表，用于错误提示
    config.AvailableInstances = make([]string, 0, len(instances))
    for _, inst := range instances {
        config.AvailableInstances = append(config.AvailableInstances, inst.ID)
        if inst.ID == config.InstanceName {
            return nil // 实例存在且运行中
        }
    }
    
    return fmt.Errorf("实例 '%s' 未找到或未运行", config.InstanceName)
}

// 列出可用实例
func listAvailableInstances(config *CliConfig) {
    if len(config.AvailableInstances) == 0 {
        fmt.Println("当前没有运行中的实例")
        return
    }
    
    fmt.Println("可用实例:")
    for _, name := range config.AvailableInstances {
        fmt.Printf("  - %s\n", name)
    }
    fmt.Println("\n使用 --instance <名称> 指定目标实例")
}

// 实例管理命令组
func newInstanceCommand(config *CliConfig) *cobra.Command {
    instCmd := &cobra.Command{
        Use:   "instance",
        Short: "实例管理命令",
        Long:  "用于管理多个配置实例的命令",
        Aliases: []string{"inst"},
    }
    
    // 列出实例的命令
    listCmd := &cobra.Command{
        Use:   "list",
        Short: "列出所有运行中的实例",
        Run: func(cmd *cobra.Command, args []string) {
            instances, err := listInstances()
            if err != nil {
                fmt.Fprintf(os.Stderr, "列出实例时出错: %v\n", err)
                os.Exit(1)
            }
            
            if len(instances) == 0 {
                fmt.Println("当前没有运行中的实例")
                return
            }
            
            fmt.Println("运行中的实例:")
            for _, inst := range instances {
                fmt.Printf("  - %s (PID: %d)\n", inst.ID, inst.PID)
            }
        },
    }
    
    // 查看实例详情的命令
    infoCmd := &cobra.Command{
        Use:   "info [实例名称]",
        Short: "查看实例详细信息",
        Args:  cobra.MaximumNArgs(1),
        Run: func(cmd *cobra.Command, args []string) {
            var instanceName string
            if len(args) > 0 {
                instanceName = args[0]
            } else {
                instanceName = config.InstanceName
            }
            
            // 获取实例信息
            inst, err := getInstanceInfo(instanceName)
            if err != nil {
                fmt.Fprintf(os.Stderr, "获取实例信息时出错: %v\n", err)
                os.Exit(1)
            }
            
            // 显示实例信息
            fmt.Printf("实例信息: %s\n", inst.ID)
            fmt.Printf("  - PID: %d\n", inst.PID)
            
            // 获取端口信息
            portManager := NewPortManager(0, inst.ID)
            if port, err := portManager.GetPortInfo(); err == nil {
                fmt.Printf("  - Web端口: %d\n", port)
                fmt.Printf("  - Web地址: http://localhost:%d\n", port)
            }
            
            // 显示实例配置目录（跨平台）
            configDir := getInstanceConfigDir(inst.ID)
            fmt.Printf("  - 配置目录: %s\n", configDir)
        },
    }
    
    // 检查实例状态的命令
    statusCmd := &cobra.Command{
        Use:   "status [实例名称]",
        Short: "检查实例状态",
        Args:  cobra.MaximumNArgs(1),
        Run: func(cmd *cobra.Command, args []string) {
            var instanceName string
            if len(args) > 0 {
                instanceName = args[0]
            } else {
                instanceName = config.InstanceName
            }
            
            // 获取实例信息
            inst, err := getInstanceInfo(instanceName)
            if err != nil {
                fmt.Fprintf(os.Stderr, "获取实例信息时出错: %v\n", err)
                os.Exit(1)
            }
            
            // 检查进程状态
            process, _ := os.FindProcess(inst.PID)
            if err := process.Signal(syscall.Signal(0)); err != nil {
                fmt.Printf("实例 %s 已停止运行\n", inst.ID)
                return
            }
            
            fmt.Printf("实例 %s 正在运行 (PID: %d)\n", inst.ID, inst.PID)
        },
    }
    
    instCmd.AddCommand(listCmd, infoCmd, statusCmd)
    return instCmd
}

// 配置管理命令组
func newConfigCommand(config *CliConfig) *cobra.Command {
    configCmd := &cobra.Command{
        Use:   "config",
        Short: "配置管理命令",
        Long:  "用于管理配置项的命令",
    }
    
    // 获取配置项的命令
    getCmd := &cobra.Command{
        Use:   "get <配置路径>",
        Short: "获取配置项值",
        Args:  cobra.ExactArgs(1),
        Run: func(cmd *cobra.Command, args []string) {
            // 连接到指定实例
            client, err := connectToInstance(config.InstanceName)
            if err != nil {
                fmt.Fprintf(os.Stderr, "连接到实例 %s 失败: %v\n", config.InstanceName, err)
                os.Exit(1)
            }
            defer client.Close()
            
            // 获取配置值
            value, err := client.GetConfig(args[0])
            if err != nil {
                fmt.Fprintf(os.Stderr, "获取配置项失败: %v\n", err)
                os.Exit(1)
            }
            
            fmt.Println(value)
        },
    }
    
    // 设置配置项的命令
    setCmd := &cobra.Command{
        Use:   "set <配置路径> <配置值>",
        Short: "设置配置项值",
        Args:  cobra.ExactArgs(2),
        Run: func(cmd *cobra.Command, args []string) {
            // 连接到指定实例
            client, err := connectToInstance(config.InstanceName)
            if err != nil {
                fmt.Fprintf(os.Stderr, "连接到实例 %s 失败: %v\n", config.InstanceName, err)
                os.Exit(1)
            }
            defer client.Close()
            
            // 设置配置值
            err = client.SetConfig(args[0], args[1])
            if err != nil {
                fmt.Fprintf(os.Stderr, "设置配置项失败: %v\n", err)
                os.Exit(1)
            }
            
            fmt.Printf("配置项 %s 已设置为: %s\n", args[0], args[1])
        },
    }
    
    configCmd.AddCommand(getCmd, setCmd)
    return configCmd
}

// 连接到指定实例
func connectToInstance(instanceName string) (*Client, error) {
    // 获取实例信息
    instance, err := getInstanceInfo(instanceName)
    if err != nil {
        return nil, err
    }
    
    // 创建客户端连接
    client, err := NewClient(instance.ID)
    if err != nil {
        return nil, err
    }
    
    return client, nil
}

// 实例信息结构
type InstanceInfo struct {
    ID  string
    PID int
}

// 获取特定实例信息
func getInstanceInfo(instanceName string) (*InstanceInfo, error) {
    instances, err := listInstances()
    if err != nil {
        return nil, err
    }
    
    for _, inst := range instances {
        if inst.ID == instanceName {
            return &inst, nil
        }
    }
    
    return nil, fmt.Errorf("实例 '%s' 未找到或未运行", instanceName)
}

// 列出所有运行中的实例（跨平台）
func listInstances() ([]InstanceInfo, error) {
    // 获取PID目录（跨平台）
    pidDir := getSystemPidDir()
    
    // 确保目录存在
    if _, err := os.Stat(pidDir); os.IsNotExist(err) {
        return []InstanceInfo{}, nil // 目录不存在时返回空列表
    }
    
    // 构建模式匹配（跨平台）
    pattern := filepath.Join(pidDir, "goci_*.pid")
    
    // 获取匹配文件
    files, err := filepath.Glob(pattern)
    if err != nil {
        return nil, err
    }
    
    instances := make([]InstanceInfo, 0, len(files))
    for _, file := range files {
        // 从文件名提取实例ID
        filename := filepath.Base(file)
        id := strings.TrimPrefix(filename, "goci_")
        id = strings.TrimSuffix(id, ".pid")
        
        // 读取PID
        data, err := os.ReadFile(file)
        if err != nil {
            continue
        }
        
        pidStr := strings.TrimSpace(string(data))
        pid, err := strconv.Atoi(pidStr)
        if err != nil {
            continue
        }
        
        // 验证进程是否存在（跨平台）
        if !isProcessRunning(file) {
            // 进程不存在，清理僵尸PID文件
            os.Remove(file)
            continue
        }
        
        // 获取配置目录（跨平台）
        configDir := filepath.Join(getSystemConfigDir(), id)
        
        instances = append(instances, InstanceInfo{
            ID:  id,
            PID: pid,
            ConfigDir: configDir,
        })
    }
    
    return instances, nil
}
```

**使用示例：**

1. 列出当前运行中的所有实例：

```bash
$ goci instance list
运行中的实例:
  - app1 (PID: 12345)
  - app2 (PID: 23456)
```

2. 查看特定实例详情：

```bash
$ goci --instance app1 instance info
实例信息: app1
  - PID: 12345
  - Web端口: 8080
  - Web地址: http://localhost:8080
  - 配置目录: /etc/goci/app1
```

3. 对特定实例操作配置：

```bash
$ goci --instance app1 config get database.host
localhost

$ goci --instance app1 config set database.port 5432
配置项 database.port 已设置为: 5432
```

4. 使用环境变量指定目标实例：

```bash
$ export GOCI_INSTANCE=app2
$ goci config get server.port
8000
```

5. 当指定的实例不存在时提供友好提示：

```bash
$ goci --instance app3 config get server.port
错误: 实例 'app3' 未找到或未运行
可用实例:
  - app1
  - app2

使用 --instance <名称> 指定目标实例
```

#### 3.6.5 跨平台支持
- Windows和Linux平台支持
  - 系统核心功能在Windows和Linux平台下保持一致
  - 配置文件路径处理兼容Windows和Linux路径格式
  - 文件锁、进程检测等机制提供平台特定实现
  - 环境变量命名与使用遵循跨平台最佳实践

- 路径与目录管理
  - Windows标准目录
    - 配置目录：`%APPDATA%\goci` 或 `%PROGRAMDATA%\goci`
    - PID文件目录：`%PROGRAMDATA%\goci\pids` 或 `%TEMP%\goci\pids`
    - 备份目录：`%APPDATA%\goci\backups`
  - Linux标准目录
    - 配置目录：`/etc/goci` 或 `~/.config/goci`
    - PID文件目录：`/var/run` 或 `/tmp/goci/pids`
    - 备份目录：`/var/lib/goci/backups` 或 `~/.local/share/goci/backups`

- 平台特定实现
  - 文件锁实现
    - Windows：使用独占文件模式 (`O_EXCL`) 模拟文件锁
    - Linux：使用系统文件锁 (`flock`)
  - 进程检测
    - Windows：使用 `OpenProcess` API 检测进程存在性
    - Linux：使用信号 (`Signal(0)`) 检测进程存在性
  - 目录权限
    - Windows：遵循 Windows 权限模型
    - Linux：遵循 POSIX 权限模型，支持用户组权限

- 跨平台编译与部署
  - 使用条件编译标签 (`// +build`) 分离平台特定代码
  - 提供统一的API接口，屏蔽平台实现差异
  - 支持单一源代码库构建不同平台的二进制文件
  - 自动检测目标平台并调整行为（如路径分隔符使用）

- 测试与兼容性保证
  - 提供Windows和Linux平台的自动化测试
  - 确保核心功能在两个平台上表现一致
  - 定期进行跨平台兼容性验证
  - 在CI/CD管道中包含跨平台测试步骤

## 4. 非功能需求

### 4.1 性能需求
- 配置读取延迟 < 100ms
- 配置写入延迟 < 200ms
- 支持并发配置访问
- 支持大规模配置项管理

### 4.2 安全需求
- 配置访问权限控制
- 配置修改审计日志
- 敏感配置加密存储
- 配置备份恢复

### 4.3 可用性需求
- 7x24小时可用
- 配置变更不影响服务
- 支持配置回滚
- 友好的错误提示

### 4.4 可维护性需求
- 模块化设计
- 清晰的代码结构
- 完善的文档
- 易于扩展
- 日志管理
  - 分模块日志
    - 支持按功能模块分离日志（配置管理、Web服务、CLI等）
    - 支持为每个模块单独配置日志级别
    - 支持模块日志独立输出到不同目标
  - 分级日志
    - 支持多级日志（DEBUG、INFO、WARN、ERROR、FATAL）
    - 支持按环境自动调整默认日志级别（开发环境详细，生产环境简洁）
    - 支持运行时动态调整日志级别，无需重启服务
  - 日志格式
    - 支持多种日志格式（文本、JSON、结构化）
    - 支持自定义日志字段和输出格式
    - 支持上下文信息（如请求ID、会话ID）贯穿整个日志链路
  - 日志配置
    - 支持通过配置文件控制所有日志行为
    - 支持通过环境变量覆盖日志配置
    - 支持热加载日志配置，无需重启服务
    - 支持按时间、大小自动轮转日志文件
  - 日志接口
    - 提供统一的日志接口，便于代码中使用
    - 支持与主流日志框架集成（如zap、logrus）
    - 支持将日志发送到第三方服务（如ELK、云日志服务）

## 5. 开发计划

### 5.1 第一阶段
- 实现核心约定
- 实现配置解析器
- 实现文件存储

### 5.2 第二阶段
- 实现REST API
- 实现CLI接口
- 实现基础Web界面

### 5.3 第三阶段
- 实现状态监控
- 完善Web界面
- 实现高级功能

## 6. 待确认事项
- 配置项的具体类型需求
- 监控指标的具体定义
- 权限控制的具体要求
- 性能指标的具体要求
- 部署环境的具体要求