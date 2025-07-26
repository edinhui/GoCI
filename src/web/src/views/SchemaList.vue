<template>
  <div class="schema-list">
    <div class="page-header">
      <h2>{{ $t('schemaList.title') }}</h2>
      <div class="header-actions">
        <el-button @click="switchLanguage" size="small" type="info">
          {{ $t('language.switch') }}
        </el-button>
      </div>
    </div>

    <el-card class="list-card">
      <template #header>
        <div class="card-header">
          <div class="header-left">
            <h3>{{ $t('schemaList.availableSchemas') }}</h3>
            <el-input
              v-model="searchQuery"
              :placeholder="$t('schemaList.search')"
              prefix-icon="el-icon-search"
              clearable
              style="width: 250px"
              class="search-input"
            />
          </div>
          <div class="header-right">
            <el-button type="primary" @click="createNewSchema">
              <el-icon><Plus /></el-icon> {{ $t('schemaList.createNew') }}
            </el-button>
          </div>
        </div>
      </template>

      <el-table
        v-loading="loading"
        :data="filteredSchemas"
        style="width: 100%"
        :empty-text="$t('schemaList.noSchemas')"
      >
        <el-table-column prop="id" :label="$t('schemaList.id')" width="180" />
        <el-table-column prop="name" :label="$t('schemaList.name')" width="180" />
        <el-table-column prop="description" :label="$t('schemaList.description')" />
        <el-table-column prop="updatedAt" :label="$t('schemaList.updatedAt')" width="180" />
        <el-table-column :label="$t('schemaList.actions')" width="250" fixed="right">
          <template #default="scope">
            <el-button size="small" @click="viewSchema(scope.row)">
              {{ $t('schemaList.view') }}
            </el-button>
            <el-button size="small" type="primary" @click="editSchema(scope.row)">
              {{ $t('schemaList.edit') }}
            </el-button>
            <el-button size="small" type="danger" @click="confirmDelete(scope.row)">
              {{ $t('schemaList.delete') }}
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import { useI18n } from 'vue-i18n'
import { schemaService } from '../services/api'

const { t, locale } = useI18n()
const router = useRouter()

// 数据
const schemas = ref([])
const loading = ref(false)
const searchQuery = ref('')

// 计算属性：过滤后的Schema列表
const filteredSchemas = computed(() => {
  if (!searchQuery.value) {
    return schemas.value
  }
  const query = searchQuery.value.toLowerCase()
  return schemas.value.filter(
    schema => 
      (schema.id && schema.id.toLowerCase().includes(query)) ||
      (schema.name && schema.name.toLowerCase().includes(query)) ||
      (schema.description && schema.description.toLowerCase().includes(query))
  )
})

// 切换语言
const switchLanguage = () => {
  locale.value = locale.value === 'zh-CN' ? 'en-US' : 'zh-CN'
}

// 加载Schema列表
const loadSchemas = async () => {
  loading.value = true
  try {
    const response = await schemaService.listSchemas()
    // 后端API现在返回统一格式，包含schemas字段
    schemas.value = response.data.schemas || []
    console.log('Loaded schemas:', schemas.value)
  } catch (error) {
    console.error('Error loading schemas:', error)
    ElMessage.error(t('schemaList.loadError'))
  } finally {
    loading.value = false
  }
}

// 创建新Schema
const createNewSchema = () => {
  router.push({ name: 'SchemaEditor' })
}

// 查看Schema
const viewSchema = (schema) => {
  router.push({ 
    name: 'SchemaViewer', 
    params: { id: schema.id }
  })
}

// 编辑Schema
const editSchema = (schema) => {
  router.push({ 
    name: 'SchemaEditorEdit', 
    params: { id: schema.id }
  })
}

// 确认删除
const confirmDelete = (schema) => {
  ElMessageBox.confirm(
    t('schemaList.deleteConfirm', { name: schema.name }),
    t('schemaList.deleteTitle'),
    {
      confirmButtonText: t('schemaList.confirmDelete'),
      cancelButtonText: t('schemaList.cancel'),
      type: 'warning'
    }
  )
    .then(() => {
      deleteSchema(schema)
    })
    .catch(() => {
      // 用户取消删除
    })
}

// 删除Schema
const deleteSchema = async (schema) => {
  try {
    await schemaService.deleteSchema(schema.id)
    ElMessage.success(t('schemaList.deleteSuccess'))
    loadSchemas() // 重新加载列表
  } catch (error) {
    console.error('Error deleting schema:', error)
    ElMessage.error(t('schemaList.deleteError'))
  }
}

// 组件挂载时加载数据
onMounted(() => {
  loadSchemas()
})
</script>

<style scoped>
.schema-list {
  width: 100%;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.page-header h2 {
  margin: 0;
  font-size: 1.8rem;
}

.list-card {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-left {
  display: flex;
  align-items: center;
}

.header-left h3 {
  margin: 0;
  margin-right: 15px;
}

.search-input {
  margin-left: 15px;
}

.header-right {
  display: flex;
  gap: 10px;
}
</style>
