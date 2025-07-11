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
          @edit:property="openEditPropertyDialog"
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
      width="600px"
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
        <el-form-item label="Description">
          <el-input v-model="currentProperty.description" type="textarea" :rows="2" placeholder="Property description" />
        </el-form-item>
        
        <!-- Validation rules based on type -->
        <el-divider content-position="left">Validation Rules</el-divider>
        
        <!-- String validation -->
        <template v-if="currentProperty.type === 'string'">
          <el-form-item label="Min Length">
            <el-input-number v-model="currentProperty.minLength" :min="0" :step="1" />
          </el-form-item>
          <el-form-item label="Max Length">
            <el-input-number v-model="currentProperty.maxLength" :min="0" :step="1" />
          </el-form-item>
          <el-form-item label="Pattern">
            <el-input v-model="currentProperty.pattern" placeholder="Regular expression pattern" />
            <span class="validation-hint">Example: ^[a-zA-Z0-9]+$</span>
          </el-form-item>
          <el-form-item label="Format">
            <el-select v-model="currentProperty.format" placeholder="Select Format" clearable>
              <el-option label="Email" value="email" />
              <el-option label="URI" value="uri" />
              <el-option label="Date" value="date" />
              <el-option label="Date-Time" value="date-time" />
              <el-option label="Hostname" value="hostname" />
              <el-option label="IPv4" value="ipv4" />
              <el-option label="IPv6" value="ipv6" />
            </el-select>
          </el-form-item>
        </template>
        
        <!-- Number validation -->
        <template v-if="currentProperty.type === 'number'">
          <el-form-item label="Minimum">
            <el-input-number v-model="currentProperty.minimum" :step="1" />
          </el-form-item>
          <el-form-item label="Maximum">
            <el-input-number v-model="currentProperty.maximum" :step="1" />
          </el-form-item>
          <el-form-item label="Multiple Of">
            <el-input-number v-model="currentProperty.multipleOf" :min="0" :step="1" />
          </el-form-item>
          <el-form-item label="Exclusive Min">
            <el-switch v-model="currentProperty.exclusiveMinimum" />
          </el-form-item>
          <el-form-item label="Exclusive Max">
            <el-switch v-model="currentProperty.exclusiveMaximum" />
          </el-form-item>
        </template>
        
        <!-- Array validation -->
        <template v-if="currentProperty.type === 'array'">
          <el-form-item label="Min Items">
            <el-input-number v-model="currentProperty.minItems" :min="0" :step="1" />
          </el-form-item>
          <el-form-item label="Max Items">
            <el-input-number v-model="currentProperty.maxItems" :min="0" :step="1" />
          </el-form-item>
          <el-form-item label="Unique Items">
            <el-switch v-model="currentProperty.uniqueItems" />
          </el-form-item>
        </template>
        
        <!-- Enum values for all types except object -->
        <template v-if="currentProperty.type !== 'object'">
          <el-form-item label="Enum Values">
            <div class="enum-input-container">
              <el-tag
                v-for="(enumValue, index) in currentProperty.enum || []"
                :key="index"
                closable
                @close="removeEnum(index)"
                class="enum-tag"
              >
                {{ enumValue }}
              </el-tag>
              <el-input
                v-if="enumInputVisible"
                ref="enumInputRef"
                v-model="enumInputValue"
                class="enum-input"
                size="small"
                @keyup.enter="addEnum"
                @blur="addEnum"
              />
              <el-button v-else class="button-new-enum" size="small" @click="showEnumInput">
                + Add Value
              </el-button>
            </div>
          </el-form-item>
        </template>
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
import { ref, computed, nextTick } from 'vue'
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
  description: '',
  // String validations
  minLength: null,
  maxLength: null,
  pattern: '',
  format: '',
  // Number validations
  minimum: null,
  maximum: null,
  multipleOf: null,
  exclusiveMinimum: false,
  exclusiveMaximum: false,
  // Array validations
  minItems: null,
  maxItems: null,
  uniqueItems: false,
  // Enum values
  enum: [],
  children: []
})
const dialogMode = ref('add')
const parentProperty = ref(null)

// Enum input handling
const enumInputVisible = ref(false)
const enumInputValue = ref('')
const enumInputRef = ref(null)

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
    description: '',
    // String validations
    minLength: null,
    maxLength: null,
    pattern: '',
    format: '',
    // Number validations
    minimum: null,
    maximum: null,
    multipleOf: null,
    exclusiveMinimum: false,
    exclusiveMaximum: false,
    // Array validations
    minItems: null,
    maxItems: null,
    uniqueItems: false,
    // Enum values
    enum: [],
    children: []
  }
  parentProperty.value = parent
  propertyDialogVisible.value = true
}

// Helper function to ensure a property has all the necessary validation fields
const ensureValidationFields = (property) => {
  const result = { ...property }
  
  // Common fields
  if (!('description' in result)) result.description = ''
  if (!('enum' in result)) result.enum = []
  if (!('children' in result)) result.children = []
  
  // Type-specific validation fields
  if (property.type === 'string') {
    if (!('minLength' in result)) result.minLength = null
    if (!('maxLength' in result)) result.maxLength = null
    if (!('pattern' in result)) result.pattern = ''
    if (!('format' in result)) result.format = ''
  } else if (property.type === 'number') {
    if (!('minimum' in result)) result.minimum = null
    if (!('maximum' in result)) result.maximum = null
    if (!('multipleOf' in result)) result.multipleOf = null
    if (!('exclusiveMinimum' in result)) result.exclusiveMinimum = false
    if (!('exclusiveMaximum' in result)) result.exclusiveMaximum = false
  } else if (property.type === 'array') {
    if (!('minItems' in result)) result.minItems = null
    if (!('maxItems' in result)) result.maxItems = null
    if (!('uniqueItems' in result)) result.uniqueItems = false
  }
  
  return result
}

// Open dialog to edit a property
const openEditPropertyDialog = (property) => {
  dialogMode.value = 'edit'
  parentProperty.value = null
  
  // Create a deep copy of the property and ensure all validation fields are initialized
  currentProperty.value = ensureValidationFields(JSON.parse(JSON.stringify(property)))
  
  propertyDialogVisible.value = true
}

// Show enum input field
const showEnumInput = () => {
  enumInputVisible.value = true
  nextTick(() => {
    enumInputRef.value?.focus()
  })
}

// Add enum value
const addEnum = () => {
  const value = enumInputValue.value.trim()
  if (value) {
    if (!currentProperty.value.enum) {
      currentProperty.value.enum = []
    }
    currentProperty.value.enum.push(value)
  }
  enumInputVisible.value = false
  enumInputValue.value = ''
}

// Remove enum value
const removeEnum = (index) => {
  currentProperty.value.enum.splice(index, 1)
}

// Helper function to extract validation fields from a property based on its type
const extractValidationFields = (property) => {
  const result = {
    id: property.id || generateId(),
    name: property.name,
    type: property.type,
    required: property.required,
    children: property.children || []
  }
  
  // Add description if available
  if (property.description) result.description = property.description
  
  // Add validation rules based on type
  if (property.type === 'string') {
    if (property.minLength !== null && property.minLength !== undefined) result.minLength = property.minLength
    if (property.maxLength !== null && property.maxLength !== undefined) result.maxLength = property.maxLength
    if (property.pattern) result.pattern = property.pattern
    if (property.format) result.format = property.format
  }
  
  if (property.type === 'number') {
    if (property.minimum !== null && property.minimum !== undefined) result.minimum = property.minimum
    if (property.maximum !== null && property.maximum !== undefined) result.maximum = property.maximum
    if (property.multipleOf !== null && property.multipleOf !== undefined) result.multipleOf = property.multipleOf
    if (property.exclusiveMinimum) result.exclusiveMinimum = true
    if (property.exclusiveMaximum) result.exclusiveMaximum = true
  }
  
  if (property.type === 'array') {
    if (property.minItems !== null && property.minItems !== undefined) result.minItems = property.minItems
    if (property.maxItems !== null && property.maxItems !== undefined) result.maxItems = property.maxItems
    if (property.uniqueItems) result.uniqueItems = true
  }
  
  // Add enum values if present
  if (property.enum && property.enum.length > 0) {
    result.enum = [...property.enum]
  }
  
  return result
}

// Save property from dialog
const saveProperty = () => {
  if (!currentProperty.value.name.trim()) {
    ElMessage.error('Property name is required')
    return
  }
  
  // Create property object with validation fields
  const propertyData = extractValidationFields(currentProperty.value)
  
  if (dialogMode.value === 'add') {
    if (parentProperty.value) {
      // Add as child property
      if (parentProperty.value.type !== 'object') {
        parentProperty.value.type = 'object'
      }
      if (!parentProperty.value.children) {
        parentProperty.value.children = []
      }
      parentProperty.value.children.push(propertyData)
      updateProperty(parentProperty.value)
    } else {
      // Add as root property
      schemaProperties.value.push(propertyData)
    }
  } else {
    // Edit existing property
    updateProperty(propertyData)
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
        // Create property schema with type and validation rules
        const propSchema = { type: prop.type }
        
        // Apply validation rules using the same helper function
        const validationFields = extractValidationFields(prop)
        
        // Copy all validation fields except id, name, type, required, and children
        // which are handled separately in the schema generation
        Object.keys(validationFields).forEach(key => {
          if (!['id', 'name', 'type', 'required', 'children'].includes(key)) {
            propSchema[key] = validationFields[key]
          }
        })
        
        // Handle object type with children
        if (prop.type === 'object' && prop.children && prop.children.length > 0) {
          propSchema.properties = {}
          propSchema.required = []
          processProperties(prop.children, propSchema)
          
          // If no required properties were found, remove the required array
          if (propSchema.required.length === 0) {
            delete propSchema.required
          }
          
          schemaObj.properties[prop.name] = propSchema
        } else {
          schemaObj.properties[prop.name] = propSchema
        }
        
        // Add to parent's required array if this property is required
        if (prop.required) {
          schemaObj.required.push(prop.name)
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

.validation-hint {
  color: #909399;
  font-size: 12px;
  line-height: 1.4;
  display: block;
  margin-top: 4px;
}

.enum-input-container {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  align-items: center;
}

.enum-tag {
  margin-right: 4px;
}

.enum-input {
  width: 120px;
  margin-right: 8px;
  vertical-align: bottom;
}

.button-new-enum {
  height: 32px;
}
</style>
