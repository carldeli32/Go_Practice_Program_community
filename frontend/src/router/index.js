import { createRouter, createWebHistory } from 'vue-router'
import { isLoggedIn } from '../stores/auth'

const routes = [
  { path: '/', name: 'Home', component: () => import('../pages/Home.vue') },
  { path: '/login', name: 'Login', component: () => import('../pages/Login.vue') },
  { path: '/post/:id', name: 'PostDetail', component: () => import('../pages/PostDetail.vue') },
  { path: '/user/:id', name: 'UserProfile', component: () => import('../pages/UserProfile.vue') },
  { path: '/create', name: 'CreatePost', component: () => import('../pages/CreatePost.vue'), meta: { requiresAuth: true } },
  { path: '/messages', name: 'Messages', component: () => import('../pages/Messages.vue'), meta: { requiresAuth: true } },
  { path: '/messages/:id', name: 'Conversation', component: () => import('../pages/Conversation.vue'), meta: { requiresAuth: true } },
  { path: '/friends', name: 'Friends', component: () => import('../pages/Friends.vue'), meta: { requiresAuth: true } },
  { path: '/admin', name: 'Admin', component: () => import('../pages/AdminPanel.vue'), meta: { requiresAuth: true } },
]

const router = createRouter({ history: createWebHistory(), routes })

router.beforeEach((to, from, next) => {
  if (to.meta.requiresAuth && !isLoggedIn()) next('/login')
  else next()
})

export default router
