import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'
import AdminLayout from '@/layouts/AdminLayout.vue'
import ProfilesView from '@/views/ProfilesView.vue'
import DevicesView from '@/views/DevicesView.vue'
import LocationsView from '@/views/LocationsView.vue'

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    redirect: '/admin/profiles'
  },
  {
    path: '/admin',
    component: AdminLayout,
    children: [
      {
        path: '',
        redirect: '/admin/profiles'
      },
      {
        path: 'profiles',
        name: 'profiles',
        component: ProfilesView,
        meta: {
          title: 'Сотрудники',
          requiresAuth: true
        }
      },
      {
        path: 'devices',
        name: 'devices',
        component: DevicesView,
        meta: {
          title: 'Устройства',
          requiresAuth: true
        }
      },
      {
        path: 'locations',
        name: 'locations',
        component: LocationsView,
        meta: {
          title: 'Локации',
          requiresAuth: true
        }
      }
    ]
  },
  {
    // Catch-all route for 404
    path: '/:pathMatch(.*)*',
    redirect: '/admin/profiles'
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
