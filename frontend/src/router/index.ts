import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'
import AdminLayout from '@/layouts/AdminLayout.vue'
import DashboardView from '@/views/DashboardView.vue'
import UsersView from '@/views/UsersView.vue'
import SettingsView from '@/views/SettingsView.vue'

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    redirect: '/admin/dashboard'
  },
  {
    path: '/admin',
    component: AdminLayout,
    children: [
      {
        path: '',
        redirect: '/admin/dashboard'
      },
      {
        path: 'dashboard',
        name: 'dashboard',
        component: DashboardView,
        meta: {
          title: 'Dashboard',
          requiresAuth: true
        }
      },
      {
        path: 'users',
        name: 'users',
        component: UsersView,
        meta: {
          title: 'Users',
          requiresAuth: true
        }
      },
      {
        path: 'settings',
        name: 'settings',
        component: SettingsView,
        meta: {
          title: 'Settings',
          requiresAuth: true
        }
      }
    ]
  },
  {
    // Catch-all route for 404
    path: '/:pathMatch(.*)*',
    redirect: '/admin/dashboard'
  }
]

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes
})

// Navigation guard for authentication (placeholder)
router.beforeEach((to, _from, next) => {
  const requiresAuth = to.matched.some(record => record.meta.requiresAuth)

  // In production, check actual auth state from store/API
  // const isAuthenticated = useAuthStore().isAuthenticated
  const isAuthenticated = true // Mock: always authenticated for now

  if (requiresAuth && !isAuthenticated) {
    // Redirect to login page
    // next('/login')
    next()
  } else {
    next()
  }
})

export default router
