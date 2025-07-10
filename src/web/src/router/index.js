import { createRouter, createWebHistory } from 'vue-router'
import SchemaEditor from '../views/SchemaEditor.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      redirect: '/schema-editor'
    },
    {
      path: '/schema-editor',
      name: 'schema-editor',
      component: SchemaEditor
    }
  ]
})

export default router
