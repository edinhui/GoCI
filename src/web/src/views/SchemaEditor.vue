<template>
  <div class="schema-editor">
    <h2>JSON Schema Editor</h2>
    <el-card class="editor-card">
      <template #header>
        <div class="card-header">
          <h3>Schema Properties</h3>
          <el-button type="primary" @click="openAddPropertyDialog">
            Add Property
          </el-button>
        </div>
      </template>
      <div class="tree-container">
        <schema-property-tree 
          :properties="schemaProperties" 
          @update:property="updateProperty"
          @delete:property="deleteProperty"
          @add:child="openAddChildPropertyDialog"
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
    
    <!-- Property Add/Edit Dialog -->
    <el-dialog
      v-model="propertyDialogVisible"
      :title="dialogMode === 'add' ? 'Add Property' : 'Edit Property'"
      width="500px"
    >
      <el-form :model="currentProperty" label-width="120px">
        <el-form-item label="Name" required>
          <el-input v-model="currentProperty.name" placeholder="Enter property name" />
        </el-form-item>
        <el-form-item label="Type" required>
          <el-select v-model="currentProperty.type" placeholder="Select Type">
            <el-option label="String" value="string" />
            <el-option label="Number" value="number" />
            <el-option label="Boolean" value="boolean" />
            <el-option label="Object" value="object" />
            <el-option label="Array" value="array" />
          </el-select>
        </el-form-item>
        <el-form-item label="Required">
          <el-switch v-model="currentProperty.required" />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="propertyDialogVisible = false">Cancel</el-button>
          <el-button type="primary" @click="saveProperty">Save</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { ElMessage } from 'element-plus'
import SchemaPropertyTree from '../components/SchemaPropertyTree.vue'

// Schema properties data structure
const schemaProperties = ref([])

// Property dialog state
const propertyDialogVisible = ref(false)
const currentProperty = ref({
  name: '',
  type: 'string',
  required: false,
  children: []
})
const dialogMode = ref('add')
const parentProperty = ref(null)

// Generate a unique ID for new properties
const generateId = () => {
  return 'prop_' + Date.now() + '_' + Math.floor(Math.random() * 1000)
}

// Open dialog to add a root property
const openAddPropertyDialog = () => {
  dialogMode.value = 'add'
  currentProperty.value = {
    name: '',
    type: 'string',
    required: false,
    children: []
  }
  parentProperty.value = null
  propertyDialogVisible.value = true
}

// Open dialog to add a child property
const openAddChildPropertyDialog = (parent) => {
  dialogMode.value = 'add'
  currentProperty.value = {
    name: '',
    type: 'string',
    required: false,
    children: []
  }
  parentProperty.value = parent
  propertyDialogVisible.value = true
}

// Save property from dialog
const saveProperty = () => {
  if (!currentProperty.value.name.trim()) {
    ElMessage.error('Property name is required')
    return
  }
  
  if (dialogMode.value === 'add') {
    const newProperty = {
      id: generateId(),
      name: currentProperty.value.name,
      type: currentProperty.value.type,
      required: currentProperty.value.required,
      children: []
    }
    
    if (parentProperty.value) {
      // Add as child property
      if (parentProperty.value.type !== 'object') {
        parentProperty.value.type = 'object'
      }
      if (!parentProperty.value.children) {
        parentProperty.value.children = []
      }
      parentProperty.value.children.push(newProperty)
      updateProperty(parentProperty.value)
    } else {
      // Add as root property
      schemaProperties.value.push(newProperty)
    }
  } else {
    // Edit existing property
    updateProperty({
      id: currentProperty.value.id,
      name: currentProperty.value.name,
      type: currentProperty.value.type,
      required: currentProperty.value.required
    })
  }
  
  propertyDialogVisible.value = false
  ElMessage.success('Property saved successfully')
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
