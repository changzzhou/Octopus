import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      name: 'home',
      component: () => import('../views/HomeView.vue'),
    },
    {
      path: '/workflows',
      name: 'workflows',
      component: () => import('../views/WorkflowsView.vue'),
    },
    {
      path: '/workflows/:id',
      name: 'workflow-detail',
      component: () => import('../views/WorkflowDetailView.vue'),
    },
    {
      path: '/workflows/:id/design',
      name: 'workflow-design',
      component: () => import('../views/WorkflowDesignView.vue'),
    },
    {
      path: '/workflows/:id/executions/:executionId',
      name: 'workflow-execution',
      component: () => import('../views/WorkflowExecutionView.vue'),
    },
  ],
})

export default router
