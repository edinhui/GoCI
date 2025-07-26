<template>
  <el-config-provider>
    <div class="app-container">
      <el-container direction="vertical">
        <el-header>
          <div class="header-content">
            <h1>Configuration Management</h1>
          </div>
        </el-header>
        <el-container class="main-container">
          <el-aside width="200px">
            <el-menu router :default-active="$route.path">
              <el-menu-item index="/schemas">
                <el-icon><Document /></el-icon>
                <span>{{ $t('sidebar.schemaManagement') }}</span>
              </el-menu-item>
              <el-menu-item index="/config-editor" disabled>
                <el-icon><Edit /></el-icon>
                <span>{{ $t('sidebar.configEditor') }}</span>
              </el-menu-item>
            </el-menu>
          </el-aside>
          <el-main>
            <router-view />
          </el-main>
        </el-container>
      </el-container>
    </div>
  </el-config-provider>
</template>

<script setup>
import { Document, Edit } from '@element-plus/icons-vue'
</script>

<style>
.app-container {
  height: 100vh;
  width: 100%;
  display: flex;
  flex-direction: column;
}

.el-header {
  background-color: #409EFF;
  color: white;
  padding: 0;
  height: auto !important;
  min-height: 80px;
}

.header-content {
  display: flex;
  align-items: center;
  height: 100%;
  padding: 0 20px;
}

.header-content h1 {
  margin: 0;
  padding: 25px 0;
  white-space: nowrap;
  overflow: visible;
  font-size: 1.7rem;
}

.main-container {
  margin-top: 10px;
  height: calc(100vh - 90px);
}

.el-aside {
  background-color: #f8f9fa;
  border-right: 1px solid #e6e6e6;
}

.el-main {
  padding: 20px;
}
</style>
