import axios from 'axios';

// 创建axios实例
const api = axios.create({
  baseURL: 'http://localhost:8080/api',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json'
  }
});

// Schema API服务
export const schemaService = {
  // 保存Schema
  saveSchema(id, name, description, schemaData) {
    return api.post(`/schemas/${id}`, {
      metadata: {
        name: name,
        description: description
      },
      schema: schemaData
    });
  },

  // 获取Schema
  getSchema(id) {
    return api.get(`/schemas/${id}`);
  },

  // 列出所有Schema
  listSchemas() {
    return api.get('/schemas');
  },

  // 删除Schema
  deleteSchema(id) {
    return api.delete(`/schemas/${id}`);
  }
};

export default api;
