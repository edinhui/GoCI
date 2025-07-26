<template>
  <div class="schema-viewer">
    <div class="viewer-header">
      <h2>{{ $t('schemaViewer.title') }}</h2>
      <div class="header-actions">
        <el-button @click="backToList" size="small">
          {{ $t('schemaViewer.backToList') }}
        </el-button>
        <el-button @click="switchLanguage" size="small" type="info">
          {{ $t('language.switch') }}
        </el-button>
      </div>
    </div>

    <el-card v-loading="loading" class="viewer-card">
      <template #header>
        <div class="card-header">
          <h3>{{ schema.name || $t('schemaViewer.schemaDetails') }}</h3>
          <el-button type="primary" @click="editSchema">
            {{ $t('schemaViewer.edit') }}
          </el-button>
        </div>
      </template>

      <div v-if="schema.id" class="schema-details">
        <div class="detail-item">
          <span class="label">{{ $t('schemaViewer.id') }}:</span>
          <span>{{ schema.id }}</span>
        </div>
        <div class="detail-item">
          <span class="label">{{ $t('schemaViewer.name') }}:</span>
          <span>{{ schema.name }}</span>
        </div>
        <div class="detail-item">
          <span class="label">{{ $t('schemaViewer.description') }}:</span>
          <span>{{ schema.description }}</span>
        </div>
        <div class="detail-item">
          <span class="label">{{ $t('schemaViewer.createdAt') }}:</span>
          <span>{{ schema.createdAt }}</span>
        </div>
        <div class="detail-item">
          <span class="label">{{ $t('schemaViewer.updatedAt') }}:</span>
          <span>{{ schema.updatedAt }}</span>
        </div>
      </div>

      <el-tabs type="border-card">
        <el-tab-pane :label="$t('schemaViewer.schemaContent')">
          <div class="schema-content">
            <pre>{{ schemaContent }}</pre>
          </div>
        </el-tab-pane>
        <el-tab-pane :label="$t('schemaViewer.configExample')">
          <div class="config-example">
            <pre>{{ configExample }}</pre>
          </div>
        </el-tab-pane>
      </el-tabs>
    </el-card>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { useI18n } from 'vue-i18n'
import { schemaService } from '../services/api'

const { t, locale } = useI18n()
const route = useRoute()
const router = useRouter()

// 数据
const schema = ref({})
const schemaData = ref({})
const loading = ref(false)

// 计算属性：格式化的Schema内容
const schemaContent = computed(() => {
  if (!schemaData.value) return ''
  try {
    return JSON.stringify(schemaData.value, null, 2)
  } catch (error) {
    return '{}'
  }
})

// 计算属性：根据Schema生成的配置示例
const configExample = computed(() => {
  if (!schemaData.value) return ''
  try {
    // 简单生成一个基于Schema的配置示例
    const example = generateConfigExample(schemaData.value)
    return JSON.stringify(example, null, 2)
  } catch (error) {
    return '{}'
  }
})

// 切换语言
const switchLanguage = () => {
  locale.value = locale.value === 'zh-CN' ? 'en-US' : 'zh-CN'
}

// 返回列表页
const backToList = () => {
  router.push({ name: 'SchemaList' })
}

// 编辑Schema
const editSchema = () => {
  router.push({ 
    name: 'SchemaEditor', 
    params: { id: schema.value.id }
  })
}

// 加载Schema详情
const loadSchema = async (id) => {
  loading.value = true
  try {
    const response = await schemaService.getSchema(id)
    
        // 响应体现在是包含 metadata 和 schema 的统一格式
    const responseData = response.data
    
    // 获取元数据
    schema.value = responseData.metadata || {
      id: id,
      name: '',
      description: '',
      createdAt: '',
      updatedAt: ''
    }
    
    // 获取Schema数据
    schemaData.value = responseData.schema || {}
    
    console.log('Schema loaded:', {
      metadata: schema.value,
      schema: schemaData.value
    })
  } catch (error) {
    console.error('Error loading schema:', error)
    ElMessage.error(t('schemaViewer.loadError'))
  } finally {
    loading.value = false
  }
}

// 根据Schema生成配置示例
const generateConfigExample = (schema) => {
  // 简单实现，根据Schema类型生成示例值
  const generateValue = (propSchema) => {
    if (!propSchema) return null
    
    if (propSchema.const) return propSchema.const
    if (propSchema.default) return propSchema.default
    
    switch (propSchema.type) {
      case 'string':
        return propSchema.title || 'Example string'
      case 'number':
      case 'integer':
        return 42
      case 'boolean':
        return true
      case 'array':
        if (propSchema.items) {
          return [generateValue(propSchema.items)]
        }
        return []
      case 'object':
        const obj = {}
        if (propSchema.properties) {
          Object.keys(propSchema.properties).forEach(key => {
            obj[key] = generateValue(propSchema.properties[key])
          })
        }
        return obj
      default:
        return null
    }
  }
  
  return generateValue(schema)
}

// 组件挂载时加载数据
onMounted(() => {
  const id = route.params.id
  if (id) {
    loadSchema(id)
  }
})
</script>

<style scoped>
.schema-viewer {
  padding: 20px;
}

.viewer-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.header-actions {
  display: flex;
  gap: 10px;
}

.viewer-card {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.schema-details {
  margin-bottom: 20px;
  padding: 15px;
  background-color: #f8f9fa;
  border-radius: 4px;
}

.detail-item {
  margin-bottom: 10px;
}

.label {
  font-weight: bold;
  margin-right: 10px;
}

.schema-content, .config-example {
  padding: 15px;
  background-color: #f8f9fa;
  border-radius: 4px;
  overflow: auto;
  max-height: 500px;
}

pre {
  margin: 0;
  white-space: pre-wrap;
}
</style>
