import { createRouter, createWebHistory } from 'vue-router'
import SchemaEditor from '../views/SchemaEditor.vue'
import SchemaList from '../views/SchemaList.vue'
import SchemaViewer from '../views/SchemaViewer.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      redirect: '/schemas'
    },
    {
      path: '/schemas',
      name: 'SchemaList',
      component: SchemaList
    },
    {
      path: '/schemas/:id/view',
      name: 'SchemaViewer',
      component: SchemaViewer,
      props: true
    },
    {
      path: '/schema-editor',
      name: 'SchemaEditor',
      component: SchemaEditor
    },
    {
      path: '/schema-editor/:id',
      name: 'SchemaEditorEdit',
      component: SchemaEditor,
      props: true
    }
  ]
})

export default router
