<template>
  <div class="schema-property-tree">
    <el-tree
      :data="formattedProperties"
      node-key="id"
      default-expand-all
    >
      <template #default="{ node, data }">
        <div class="property-node">
          <div class="property-info">
            <span class="property-name">{{ data.name }}</span>
            <el-tag size="small" :type="getTypeColor(data.type)">
              {{ data.type }}
            </el-tag>
            <el-tag v-if="data.required" size="small" type="danger">Required</el-tag>
          </div>
          <div class="property-actions">
            <el-button-group>
              <el-button size="small" @click.stop="editProperty(data)">
                <el-icon><Edit /></el-icon>
              </el-button>
              <el-button 
                size="small" 
                v-if="data.type === 'object'" 
                @click.stop="emitAddChild(data)"
              >
                <el-icon><Plus /></el-icon>
              </el-button>
              <el-button size="small" type="danger" @click.stop="deleteProperty(data)">
                <el-icon><Delete /></el-icon>
              </el-button>
            </el-button-group>
          </div>
        </div>
      </template>
    </el-tree>

    <!-- Property Edit Dialog -->
    <el-dialog
      v-model="editDialogVisible"
      title="Edit Property"
      width="500px"
    >
      <el-form :model="currentProperty" label-width="120px">
        <el-form-item label="Name">
          <el-input v-model="currentProperty.name" />
        </el-form-item>
        <el-form-item label="Type">
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
          <el-button @click="editDialogVisible = false">Cancel</el-button>
          <el-button type="primary" @click="saveProperty">Save</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { Edit, Plus, Delete } from '@element-plus/icons-vue'

const props = defineProps({
  properties: {
    type: Array,
    required: true
  }
})

const emit = defineEmits(['update:property', 'delete:property', 'add:child'])

// Format properties for the tree component
const formattedProperties = computed(() => {
  return props.properties.map(prop => ({
    ...prop,
    label: prop.name
  }))
})

// Property editing
const editDialogVisible = ref(false)
const currentProperty = ref({})
const currentPropertyId = ref(null)

// Generate a unique ID for new properties
const generateId = () => {
  return 'prop_' + Date.now() + '_' + Math.floor(Math.random() * 1000)
}

// Get color for type tag
const getTypeColor = (type) => {
  const typeColors = {
    'string': 'success',
    'number': 'warning',
    'boolean': 'info',
    'object': 'primary',
    'array': 'danger'
  }
  return typeColors[type] || ''
}

// Edit a property
const editProperty = (property) => {
  currentPropertyId.value = property.id
  currentProperty.value = { ...property }
  editDialogVisible.value = true
}

// Save property changes
const saveProperty = () => {
  emit('update:property', {
    id: currentPropertyId.value,
    ...currentProperty.value
  })
  editDialogVisible.value = false
}

// Emit event to add a child property
const emitAddChild = (parent) => {
  emit('add:child', parent)
}

// Delete a property
const deleteProperty = (property) => {
  emit('delete:property', property.id)
}
</script>

<style scoped>
.schema-property-tree {
  width: 100%;
}

.property-node {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 5px 0;
  width: 100%;
}

.property-info {
  display: flex;
  align-items: center;
  gap: 10px;
}

.property-name {
  font-weight: bold;
}

.property-actions {
  margin-left: 20px;
}
</style>
