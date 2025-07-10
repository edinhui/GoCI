<template>
  <div class="schema-editor">
    <h2>JSON Schema Editor</h2>
    <el-card class="editor-card">
      <template #header>
        <div class="card-header">
          <h3>Schema Properties</h3>
          <el-button type="primary" @click="addRootProperty">
            Add Property
          </el-button>
        </div>
      </template>
      <div class="tree-container">
        <schema-property-tree 
          :properties="schemaProperties" 
          @update:property="updateProperty"
          @delete:property="deleteProperty"
        />
      </div>
    </el-card>

    <el-card class="preview-card">
      <template #header>
        <div class="card-header">
          <h3>Schema Preview</h3>
          <el-button type="success" @click="exportSchema">
            Export Schema
          </el-button>
        </div>
      </template>
      <div class="preview-container">
        <pre>{{ schemaPreview }}</pre>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { ElMessage } from 'element-plus'
import SchemaPropertyTree from '../components/SchemaPropertyTree.vue'

// Schema properties data structure
const schemaProperties = ref([])

// Generate a unique ID for new properties
const generateId = () => {
  return 'prop_' + Date.now() + '_' + Math.floor(Math.random() * 1000)
}

// Add a new root property
const addRootProperty = () => {
  schemaProperties.value.push({
    id: generateId(),
    name: 'newProperty',
    type: 'string',
    required: false,
    children: []
  })
}

// Update a property in the schema
const updateProperty = (property) => {
  const findAndUpdate = (properties, id, updatedProperty) => {
    for (let i = 0; i < properties.length; i++) {
      if (properties[i].id === id) {
        properties[i] = { ...properties[i], ...updatedProperty }
        return true
      }
      if (properties[i].children && properties[i].children.length > 0) {
        if (findAndUpdate(properties[i].children, id, updatedProperty)) {
          return true
        }
      }
    }
    return false
  }

  findAndUpdate(schemaProperties.value, property.id, property)
}

// Delete a property from the schema
const deleteProperty = (id) => {
  const findAndDelete = (properties, id) => {
    for (let i = 0; i < properties.length; i++) {
      if (properties[i].id === id) {
        properties.splice(i, 1)
        return true
      }
      if (properties[i].children && properties[i].children.length > 0) {
        if (findAndDelete(properties[i].children, id)) {
          return true
        }
      }
    }
    return false
  }

  findAndDelete(schemaProperties.value, id)
}

// Generate a preview of the JSON Schema
const schemaPreview = computed(() => {
  const generateSchema = (properties) => {
    const schema = {
      $schema: "http://json-schema.org/draft-07/schema#",
      type: "object",
      properties: {},
      required: []
    }

    const processProperties = (properties, schemaObj) => {
      for (const prop of properties) {
        if (prop.type === 'object' && prop.children && prop.children.length > 0) {
          schemaObj.properties[prop.name] = {
            type: 'object',
            properties: {}
          }
          if (prop.required) {
            schemaObj.required.push(prop.name)
          }
          processProperties(prop.children, schemaObj.properties[prop.name])
        } else {
          schemaObj.properties[prop.name] = { type: prop.type }
          if (prop.required) {
            schemaObj.required.push(prop.name)
          }
        }
      }
    }

    processProperties(properties, schema)
    return schema
  }

  try {
    const schema = generateSchema(schemaProperties.value)
    return JSON.stringify(schema, null, 2)
  } catch (error) {
    console.error('Error generating schema preview:', error)
    return '{ "error": "Could not generate schema preview" }'
  }
})

// Export the schema as a JSON file
const exportSchema = () => {
  try {
    const schema = JSON.parse(schemaPreview.value)
    const dataStr = "data:text/json;charset=utf-8," + encodeURIComponent(JSON.stringify(schema, null, 2))
    const downloadAnchorNode = document.createElement('a')
    downloadAnchorNode.setAttribute("href", dataStr)
    downloadAnchorNode.setAttribute("download", "schema.json")
    document.body.appendChild(downloadAnchorNode)
    downloadAnchorNode.click()
    downloadAnchorNode.remove()
    ElMessage({
      message: 'Schema exported successfully',
      type: 'success'
    })
  } catch (error) {
    console.error('Error exporting schema:', error)
    ElMessage({
      message: 'Error exporting schema',
      type: 'error'
    })
  }
}
</script>

<style scoped>
.schema-editor {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.editor-card, .preview-card {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.tree-container {
  min-height: 300px;
}

.preview-container {
  background-color: #f5f7fa;
  border-radius: 4px;
  padding: 10px;
  max-height: 400px;
  overflow: auto;
}

pre {
  margin: 0;
  white-space: pre-wrap;
  word-wrap: break-word;
}
</style>
