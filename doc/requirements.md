# 配置管理系统需求文档

## 1. 系统概述

本系统旨在提供一个基于约定的配置管理平台，用于管理Golang程序的配置项。系统将提供多种配置访问方式，并支持程序状态监控功能。

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
    components:
      schemas:
        # 基础类型
        String:
          type: string
        Integer:
          type: integer
        Float:
          type: number
        Boolean:
          type: boolean
        DateTime:
          type: string
          format: date-time
        
        # 网络相关
        IPv4Address:
          type: string
          format: ipv4
        IPv6Address:
          type: string
          format: ipv6
        MACAddress:
          type: string
          pattern: "^([0-9A-Fa-f]{2}[:-]){5}([0-9A-Fa-f]{2})$"
        Port:
          type: integer
          minimum: 0
          maximum: 65535
        
        # 文件系统
        FilePath:
          type: string
          pattern: "^[a-zA-Z0-9_./-]+$"
        DirectoryPath:
          type: string
          pattern: "^[a-zA-Z0-9_./-]+$"
        
        # 安全相关
        Password:
          type: string
          minLength: 8
          pattern: "^(?=.*[A-Za-z])(?=.*\d)[A-Za-z\d]{8,}$"
        Email:
          type: string
          format: email
    ```

  - 配置Schema组织方式
    - 单一文件方式（推荐用于简单配置）
      ```yaml
      # app_config.yaml
      openapi: 3.0.0
      info:
        title: Application Configuration
        version: 1.0.0
      components:
        schemas:
          AppConfig:
            type: object
            properties:
              server:
                $ref: './common_types.yaml#/components/schemas/ServerConfig'
              database:
                $ref: './common_types.yaml#/components/schemas/DatabaseConfig'
      ```

    - 多文件方式（推荐用于复杂配置）
      ```yaml
      # config.yaml
      openapi: 3.0.0
      info:
        title: Application Configuration
        version: 1.0.0
      components:
        schemas:
          AppConfig:
            type: object
            properties:
              server:
                $ref: './schemas/server.yaml#/ServerConfig'
              database:
                $ref: './schemas/database.yaml#/DatabaseConfig'
              logging:
                $ref: './schemas/logging.yaml#/LoggingConfig'
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

  - 目录结构建议
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

### 2.2 元数据约定
- 基本信息（必填）
  - 配置项标识（ID）：唯一标识符
  - 配置项名称（Name）：显示名称
  - 配置项描述（Description）：详细说明
  - 配置项分组（Group）：所属分组

- 配置属性（可选）
  - 默认值（Default Value）
  - 必填标记（Required）
  - 只读标记（ReadOnly）
  - 验证规则（Validation Rules）
    - 最小值/最大值
    - 正则表达式
    - 枚举值列表

### 2.3 分组约定
- 分组命名约定
  - 使用小写字母和下划线
  - 反映功能模块
  - 最多三级分组
  - 示例：
    - `system`
    - `network/interfaces`
    - `security/firewall/rules`

- 分组结构约定
  - 支持多级分组（建议不超过3级）
  - 分组名称和描述
  - 分组图标（可选）

### 2.4 前端展示约定
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

### 2.5 权限约定
- 基于分组的权限控制
- 权限级别：
  - `read`: 只读
  - `write`: 可修改
  - `admin`: 完全控制
- 权限继承规则

### 2.7 配置Schema组织
- 预定义通用类型
  - 系统提供的通用类型定义
    ```yaml
    # common_types.yaml
    openapi: 3.0.0
    info:
      title: Common Types Definition
      version: 1.0.0
    components:
      schemas:
        # 基础类型
        String:
          type: string
        Integer:
          type: integer
        Float:
          type: number
        Boolean:
          type: boolean
        DateTime:
          type: string
          format: date-time
        
        # 网络相关
        IPv4Address:
          type: string
          format: ipv4
        IPv6Address:
          type: string
          format: ipv6
        MACAddress:
          type: string
          pattern: "^([0-9A-Fa-f]{2}[:-]){5}([0-9A-Fa-f]{2})$"
        Port:
          type: integer
          minimum: 0
          maximum: 65535
        
        # 文件系统
        FilePath:
          type: string
          pattern: "^[a-zA-Z0-9_./-]+$"
        DirectoryPath:
          type: string
          pattern: "^[a-zA-Z0-9_./-]+$"
        
        # 安全相关
        Password:
          type: string
          minLength: 8
          pattern: "^(?=.*[A-Za-z])(?=.*\d)[A-Za-z\d]{8,}$"
        Email:
          type: string
          format: email
    ```

- 配置Schema组织方式
  - 单一文件方式（推荐用于简单配置）
    ```yaml
    # app_config.yaml
    openapi: 3.0.0
    info:
      title: Application Configuration
      version: 1.0.0
    components:
      schemas:
        AppConfig:
          type: object
          properties:
            server:
              $ref: './common_types.yaml#/components/schemas/ServerConfig'
            database:
              $ref: './common_types.yaml#/components/schemas/DatabaseConfig'
    ```

  - 多文件方式（推荐用于复杂配置）
    ```yaml
    # config.yaml
    openapi: 3.0.0
    info:
      title: Application Configuration
      version: 1.0.0
    components:
      schemas:
        AppConfig:
          type: object
          properties:
            server:
              $ref: './schemas/server.yaml#/ServerConfig'
            database:
              $ref: './schemas/database.yaml#/DatabaseConfig'
            logging:
              $ref: './schemas/logging.yaml#/LoggingConfig'
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

- 目录结构建议
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

## 3. 功能需求

### 3.1 配置管理功能
- 配置项定义
- 配置项验证
- 配置项导入/导出
- 配置项版本控制

### 3.2 访问接口
- 命令行接口（CLI）
- Web界面
- REST API

### 3.3 状态监控
- 程序运行状态显示
- 自定义监控指标
- 实时状态更新
- 状态历史记录

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