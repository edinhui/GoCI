<template>
  <div class="schema-editor">
    <div class="editor-header">
      <h2>{{ $t('schemaEditor.title') }}</h2>
      <div class="header-actions">
        <el-button @click="switchLanguage" size="small" type="info">
          {{ $t('language.switch') }}
        </el-button>
      </div>
    </div>
    <el-card class="editor-card">
      <template #header>
        <div class="card-header">
          <h3>{{ $t('schemaEditor.properties') }}</h3>
          <el-button type="primary" @click="openAddPropertyDialog">
            {{ $t('schemaEditor.addProperty') }}
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
          <h3>{{ $t('schemaEditor.preview') }}</h3>
          <el-button type="success" @click="exportSchema">
            {{ $t('schemaEditor.exportSchema') }}
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
      :title="dialogMode === 'add' ? $t('propertyDialog.add') : $t('propertyDialog.edit')"
      width="600px"
    >
      <el-form :model="currentProperty" label-width="120px">
        <el-form-item :label="$t('propertyDialog.name')" required>
          <el-input v-model="currentProperty.name" :placeholder="`${$t('propertyDialog.name')}`" :disabled="currentProperty.isFixed" />
        </el-form-item>
        
        <el-form-item :label="$t('propertyDialog.type')" required>
          <el-select v-model="currentProperty.type" :placeholder="`${$t('propertyDialog.type')}`" style="width: 100%" :disabled="currentProperty.isFixed">
            <el-option :label="$t('propertyTypes.string')" value="string" />
            <el-option :label="$t('propertyTypes.number')" value="number" />
            <el-option :label="$t('propertyTypes.boolean')" value="boolean" />
            <el-option :label="$t('propertyTypes.object')" value="object" />
            <el-option :label="$t('propertyTypes.array')" value="array" />
          </el-select>
        </el-form-item>
        
        <el-form-item :label="$t('propertyDialog.required')">
          <el-switch v-model="currentProperty.required" :disabled="currentProperty.isFixed" />
        </el-form-item>
        
        <el-form-item :label="$t('propertyDialog.description')">
          <el-input 
            v-model="currentProperty.description" 
            type="textarea" 
            :placeholder="`${$t('propertyDialog.description')}`"
            :disabled="currentProperty.isFixed"
          />
        </el-form-item>
        
        <el-form-item v-if="currentProperty.isFixed" :label="$t('propertyDialog.fixedFieldValue')" required>
          <el-input 
            v-model="currentProperty.value" 
            :placeholder="`${$t('propertyDialog.fixedFieldValue')}`"
          />
          <span class="validation-hint">{{ $t('propertyDialog.fixedFieldHint') }}</span>
        </el-form-item>
        
        <!-- Validation rules based on type -->
        <el-divider content-position="left">{{ $t('schemaEditor.validationRules') }}</el-divider>
        
        <!-- String validation -->
        <template v-if="currentProperty.type === 'string'">
          <el-form-item :label="$t('validation.string.minLength')">
            <el-input-number v-model="currentProperty.minLength" :min="0" :step="1" />
          </el-form-item>
          <el-form-item :label="$t('validation.string.maxLength')">
            <el-input-number v-model="currentProperty.maxLength" :min="0" :step="1" />
          </el-form-item>
          <el-form-item :label="$t('validation.string.pattern')">
            <el-input v-model="currentProperty.pattern" :placeholder="$t('validation.string.pattern')" />
            <span class="validation-hint">{{ $t('validation.string.patternHint') }}</span>
          </el-form-item>
          <el-form-item :label="$t('validation.string.format')">
            <el-select v-model="currentProperty.format" :placeholder="$t('validation.string.format')" clearable>
              <el-option :label="$t('validation.formats.email')" value="email" />
              <el-option :label="$t('validation.formats.uri')" value="uri" />
              <el-option :label="$t('validation.formats.date')" value="date" />
              <el-option :label="$t('validation.formats.dateTime')" value="date-time" />
              <el-option :label="$t('validation.formats.hostname')" value="hostname" />
              <el-option :label="$t('validation.formats.ipv4')" value="ipv4" />
              <el-option :label="$t('validation.formats.ipv6')" value="ipv6" />
            </el-select>
          </el-form-item>
        </template>
        
        <!-- Number validation -->
        <template v-if="currentProperty.type === 'number'">
          <el-form-item :label="$t('validation.number.minimum')">
            <el-input-number v-model="currentProperty.minimum" :step="1" />
          </el-form-item>
          <el-form-item :label="$t('validation.number.maximum')">
            <el-input-number v-model="currentProperty.maximum" :step="1" />
          </el-form-item>
          <el-form-item :label="$t('validation.number.multipleOf')">
            <el-input-number v-model="currentProperty.multipleOf" :min="0" :step="1" />
          </el-form-item>
          <el-form-item :label="$t('validation.number.exclusiveMinimum')">
            <el-switch v-model="currentProperty.exclusiveMinimum" />
          </el-form-item>
          <el-form-item :label="$t('validation.number.exclusiveMaximum')">
            <el-switch v-model="currentProperty.exclusiveMaximum" />
          </el-form-item>
        </template>
        
        <!-- Array validation -->
        <template v-if="currentProperty.type === 'array'">
          <el-form-item :label="$t('validation.array.minItems')">
            <el-input-number v-model="currentProperty.minItems" :min="0" :step="1" />
          </el-form-item>
          <el-form-item :label="$t('validation.array.maxItems')">
            <el-input-number v-model="currentProperty.maxItems" :min="0" :step="1" />
          </el-form-item>
          <el-form-item :label="$t('validation.array.uniqueItems')">
            <el-switch v-model="currentProperty.uniqueItems" />
          </el-form-item>
        </template>
        
        <!-- Enum values for all types except object -->
        <template v-if="currentProperty.type !== 'object'">
          <el-form-item :label="$t('validation.enum.title')">
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
                {{ $t('validation.enum.addValue') }}
              </el-button>
            </div>
          </el-form-item>
        </template>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="propertyDialogVisible = false">{{ $t('propertyDialog.cancel') }}</el-button>
          <el-button type="primary" @click="saveProperty">{{ $t('propertyDialog.save') }}</el-button>
        </span>
      </template>
    </el-dialog>
    
    <!-- 固定字段值设置对话框 -->
    <el-dialog
      v-model="fixedFieldDialogVisible"
      title="设置固定字段值"
      width="600px"
    >
      <p>请为以下固定字段设置值，这些值在后续配置中将无法修改。</p>
      
      <el-form label-width="120px">
        <el-form-item 
          v-for="field in currentFixedFields" 
          :key="field.id" 
          :label="field.name" 
          :required="field.required"
        >
          <el-input 
            v-model="field.value" 
            :placeholder="`请输入${field.name}值`" 
          />
          <span class="validation-hint">{{ field.description }}</span>
        </el-form-item>
      </el-form>
      
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="cancelFixedFieldsDialog">{{ $t('propertyDialog.cancel') }}</el-button>
          <el-button type="primary" @click="saveFixedFieldsDialog">{{ $t('propertyDialog.save') }}</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, nextTick, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { useI18n } from 'vue-i18n'
import SchemaPropertyTree from '../components/SchemaPropertyTree.vue'
import fixedFieldsConfig from '../config/fixed-fields.json'

// i18n 实例
const { t, locale } = useI18n()

// 切换语言
const switchLanguage = () => {
  // 切换语言
  locale.value = locale.value === 'zh' ? 'en' : 'zh'
  // 保存语言设置到localStorage
  localStorage.setItem('language', locale.value)
  
  // 提示用户语言已切换
  ElMessage({
    message: locale.value === 'zh' ? '已切换到中文' : 'Switched to English',
    type: 'success'
  })
}

// 固定字段配置 - 从JSON文件加载
const FIXED_FIELDS = fixedFieldsConfig

// Schema properties data structure
const schemaProperties = ref([])

// 初始化Schema，添加根级固定字段
const initSchema = () => {
  // 打开根级固定字段设置对话框
  openFixedFieldsDialog(null, FIXED_FIELDS.root)
}

// 组件挂载时初始化Schema
onMounted(() => {
  initSchema()
})

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
  children: [],
  // 固定字段值
  value: undefined
})
const dialogMode = ref('add')
const parentProperty = ref(null)

// 固定字段值设置对话框状态
const fixedFieldDialogVisible = ref(false)
const currentFixedFields = ref([])
const currentObjectProperty = ref(null)

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
    children: property.children || [],
    // 保留固定字段标记、只读属性和值
    isFixed: property.isFixed || false,
    readOnly: property.readOnly || false,
    value: property.value
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
      addChildProperty(propertyData, parentProperty.value)
    } else {
      // Add as root property
      addProperty(propertyData)
    }
  } else {
    // Edit existing property
    updateProperty(propertyData)
  }
  
  propertyDialogVisible.value = false
  ElMessage.success('Property saved successfully')
}

// Add a property to the schema
const addProperty = (property) => {
  schemaProperties.value.push(property)
  
  // 如果是对象类型，自动添加固定子字段
  if (property.type === 'object' && !property.isFixed) {
    addFixedFieldsToObject(property)
  }
}

// 为对象类型属性添加固定子字段
const addFixedFieldsToObject = (objectProperty) => {
  if (!objectProperty.children) {
    objectProperty.children = []
  }
  
  // 打开对象级固定字段设置对话框
  openFixedFieldsDialog(objectProperty, FIXED_FIELDS.object)
}

// Add a child property to a parent property
const addChildProperty = (childProperty, parentProperty) => {
  if (!parentProperty.children) {
    parentProperty.children = []
  }
  parentProperty.children.push(childProperty)
  
  // 如果添加的是对象类型，自动添加固定子字段
  if (childProperty.type === 'object' && !childProperty.isFixed) {
    addFixedFieldsToObject(childProperty)
  }
  
  updateProperty(parentProperty)
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
      // 不允许删除固定字段
      if (properties[i].isFixed) {
        continue
      }
      
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

// 打开固定字段设置对话框
const openFixedFieldsDialog = (objectProperty, fixedFieldsTemplate) => {
  currentObjectProperty.value = objectProperty
  
  // 创建固定字段列表副本
  currentFixedFields.value = fixedFieldsTemplate.map(field => ({
    ...field,
    id: generateId(),
    value: ''
  }))
  
  fixedFieldDialogVisible.value = true
}

// 保存固定字段设置
const saveFixedFieldsDialog = () => {
  // 验证必填字段
  const missingRequired = currentFixedFields.value.find(
    field => field.required && (!field.value || !field.value.trim())
  )
  
  if (missingRequired) {
    ElMessage.error(`${missingRequired.name} 是必填字段`)
    return
  }
  
  if (currentObjectProperty.value === null) {
    // 根级固定字段
    currentFixedFields.value.forEach(field => {
      schemaProperties.value.push(field)
    })
  } else {
    // 对象级固定字段
    if (!currentObjectProperty.value.children) {
      currentObjectProperty.value.children = []
    }
    
    currentFixedFields.value.forEach(field => {
      currentObjectProperty.value.children.push(field)
    })
    
    // 更新对象属性
    updateProperty(currentObjectProperty.value)
  }
  
  fixedFieldDialogVisible.value = false
  ElMessage.success('固定字段设置成功')
}

// 取消固定字段设置
const cancelFixedFieldsDialog = () => {
  // 如果是初始化时取消，使用默认值
  if (currentObjectProperty.value === null) {
    // 添加根级固定字段，使用默认值
    FIXED_FIELDS.root.forEach(field => {
      const fixedField = { ...field, id: generateId() }
      schemaProperties.value.push(fixedField)
    })
  } else {
    // 添加对象级固定字段，使用默认值
    FIXED_FIELDS.object.forEach(field => {
      const fixedField = { ...field, id: generateId() }
      currentObjectProperty.value.children.push(fixedField)
    })
    
    // 更新对象属性
    updateProperty(currentObjectProperty.value)
  }
  
  fixedFieldDialogVisible.value = false
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
        
        // 如果有值，添加默认值
        // 注意：对于固定字段，在JSON Schema中使用default来设置默认值
        // 这样在配置编辑器中就会默认显示这个值
        // 结合readOnly属性，实现了值不可修改的效果
        if (prop.value !== undefined && prop.value !== null && prop.value !== '') {
          propSchema.default = prop.value
        }
        
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
  padding: 10px;
}

.editor-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  padding: 10px 0;
  border-bottom: 1px solid #ebeef5;
}

.editor-header h2 {
  margin: 0;
  color: #409EFF;
  font-size: 24px;
}

.header-actions {
  display: flex;
  gap: 10px;
  align-items: center;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.editor-card {
  margin-bottom: 20px;
}

.preview-card {
  margin-bottom: 20px;
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
