# Configuration Management Component Architecture

## 1. Frontend Technology Selection (Vue 3)

### JSON Schema Creation Page Framework
- **Graphical Editor**:
  - Using a tree-list format similar to **[vue-json-schema-editor-visual](https://github.com/zyqwst/json-schema-editor-vue3)**, providing a graphical editing and generation interface
  - Supporting drag-and-drop operations, node addition/deletion/editing
  - Real-time preview of generated JSON Schema

- **Editor Components**:
  - **Form-based Editors**: Visual form builders that generate JSON Schema without code editing
  - **Drag-and-Drop Interface**: Intuitive UI for constructing schema elements
  - **Property Panels**: Graphical controls for configuring schema properties
  - **Visual Preview**: Real-time rendering of the schema structure

### CSS Style
- **UI Framework**:
  - **Element Plus**: Vue 3 component library with complete design system
    - Rich component library (forms, tables, navigation, etc.)
    - Comprehensive theme customization
    - Internationalization support
    - Active community and thorough documentation

## 2. JSON File Browsing and Editing Framework

- **Form-based Configuration Editor**:
  - Schema-driven form generation
  - Field-specific input controls
  - Validation and error highlighting
  - Configuration version comparison

- **Visualization Components**:
  - Hierarchical navigation tree
  - Configuration diff viewer
  - Relationship diagrams for complex configurations

## 3. Backend Webserver Framework

Recommended frameworks for bundling frontend with Golang:

- **Gin**: High-performance, feature-rich web framework
- **Echo**: Minimalist web framework with elegant APIs
- **Fiber**: Express-inspired framework with excellent performance

**Frontend Bundling Solutions**:
- `go:embed` (Go 1.16+) for embedding compiled frontend resources
- Tools like `statik` or `packr` for bundling static assets

## 4. Directory Structure Design

### Directory Structure Design

```
/schemas/                # JSON Schema files directory
  /schema1/              # Each Schema has its own directory
    schema1.json         # Schema definition
    metadata.json        # Schema metadata
  /schema2/
    schema2.json
    metadata.json
  schema-registry.json   # Schema registry (metadata management)

/configs/                # JSON configuration files directory
  /schema1/              # Each config corresponds to a Schema
    config.json          # Configuration file (3 levels)
                         # Level 1: Navigation bar
                         # Level 2: Sidebar
                         # Level 3: Page
    history/             # Version history
      config_v1.json
      config_v2.json
  /schema2/
    config.json
    history/
      config_v1.json
  config-registry.json   # Configuration registry
```

### Interaction Flow
- JSON Schema defines structure/validation rules
- Frontend generates forms dynamically from Schema
- Users create/modify configs through form interfaces
- Backend handles Schema/config reading, validation, persistence
- Config editors use Schema for real-time validation

## 5. Testability Considerations

### Frontend Testing
- **Unit Tests**: Vitest/Jest for Vue components
- **Component Tests**: Vue Test Utils for interactions
- **E2E Tests**: Cypress/Playwright for UI automation
- **Mocking**: MSW for API responses

### Backend Testing
- **Unit Tests**: Go's testing package
- **HTTP Tests**: httptest for API endpoints
- **Integration Tests**: testify for assertions/mocks
- **Benchmarks**: Go's benchmark testing

### Testing Strategy
- Isolated component/interface testing
- Dependency injection for test simplicity
- Contract testing for API consistency
- CI/CD pipeline integration
