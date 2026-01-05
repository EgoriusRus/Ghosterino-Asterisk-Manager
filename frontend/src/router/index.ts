import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'
import { useAuth } from '@/stores/auth'
import AdminLayout from '@/layouts/AdminLayout.vue'
import LoginView from '@/views/LoginView.vue'
import ProfilesView from '@/views/ProfilesView.vue'
import DevicesView from '@/views/DevicesView.vue'
import LocationsView from '@/views/LocationsView.vue'

const routes: RouteRecordRaw[] = [
  {
    path: '/login',
    name: 'login',
    component: LoginView,
    meta: {
      title: 'Вход',
      requiresAuth: false
    }
  },
  {
    path: '/',
    redirect: '/admin/profiles'
  },
  {
    path: '/admin',
    component: AdminLayout,
    meta: {
      requiresAuth: true
    },
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

// Navigation guard for authentication
router.beforeEach((to, _from, next) => {
  const { isAuthenticated } = useAuth()
  const requiresAuth = to.matched.some(record => record.meta.requiresAuth)

  if (requiresAuth && !isAuthenticated.value) {
    // Redirect to login page
    next('/login')
  } else if (to.path === '/login' && isAuthenticated.value) {
    // Already logged in, redirect to admin
    next('/admin/profiles')
  } else {
    next()
  }
})

export default router
