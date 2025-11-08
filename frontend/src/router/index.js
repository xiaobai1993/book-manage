import { createRouter, createWebHistory } from 'vue-router'
import { useUserStore } from '@/stores/user'

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/Login.vue'),
    meta: { requiresAuth: false }
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('@/views/Register.vue'),
    meta: { requiresAuth: false }
  },
  {
    path: '/forget-password',
    name: 'ForgetPassword',
    component: () => import('@/views/ForgetPassword.vue'),
    meta: { requiresAuth: false }
  },
  {
    path: '/',
    component: () => import('@/layouts/MainLayout.vue'),
    redirect: '/books',
    meta: { requiresAuth: true },
    children: [
      {
        path: 'books',
        name: 'Books',
        component: () => import('@/views/Books.vue'),
        meta: { requiresAuth: true }
      },
      {
        path: 'book/:id',
        name: 'BookDetail',
        component: () => import('@/views/BookDetail.vue'),
        meta: { requiresAuth: true }
      },
      {
        path: 'book-add',
        name: 'BookAdd',
        component: () => import('@/views/BookAdd.vue'),
        meta: { requiresAuth: true, requiresAdmin: true }
      },
      {
        path: 'book-edit/:id',
        name: 'BookEdit',
        component: () => import('@/views/BookEdit.vue'),
        meta: { requiresAuth: true, requiresAdmin: true }
      },
      {
        path: 'my-borrows',
        name: 'MyBorrows',
        component: () => import('@/views/MyBorrows.vue'),
        meta: { requiresAuth: true }
      },
      {
        path: 'all-borrows',
        name: 'AllBorrows',
        component: () => import('@/views/AllBorrows.vue'),
        meta: { requiresAuth: true, requiresAdmin: true }
      },
      {
        path: 'profile',
        name: 'Profile',
        component: () => import('@/views/Profile.vue'),
        meta: { requiresAuth: true }
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 路由守卫
router.beforeEach(async (to, from, next) => {
  const userStore = useUserStore()
  
  // 检查是否需要登录
  if (to.meta.requiresAuth && !userStore.isLoggedIn) {
    next('/login')
    return
  }
  
  // 检查是否需要管理员权限
  if (to.meta.requiresAdmin && !userStore.isAdmin) {
    const { ElMessage } = await import('element-plus')
    ElMessage.error('权限不足')
    next('/books')
    return
  }
  
  // 如果已登录，访问登录页则重定向到首页
  if (to.path === '/login' && userStore.isLoggedIn) {
    next('/books')
    return
  }
  
  next()
})

export default router
