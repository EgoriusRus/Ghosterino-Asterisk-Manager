# Admin Layout Documentation

## Overview

This document describes the admin layout implementation for the Vue 3 + TypeScript + Tailwind CSS application. The layout provides a professional, responsive admin panel interface with sidebar navigation, header, and content area.

---

## Architecture

### Component Hierarchy

```
App.vue (router-view)
└── AdminLayout.vue
    ├── Sidebar.vue (layouts/components/)
    ├── Header.vue (layouts/components/)
    └── router-view (content area)
        ├── DashboardView.vue
        ├── UsersView.vue
        └── SettingsView.vue
```

### Technology Stack

- **Vue 3.4.21** - Progressive JavaScript framework
- **TypeScript 5.6.3** - Static type checking
- **Tailwind CSS 3.4.1** - Utility-first CSS framework
- **Vue Router 4.6.4** - Official router for Vue.js
- **Vite 5.1.4** - Next generation frontend tooling

---

## Components

### 1. AdminLayout.vue

**Location:** `/src/layouts/AdminLayout.vue`

**Purpose:** Main layout wrapper that orchestrates sidebar, header, and content area.

**Features:**
- Manages sidebar state (collapsed/expanded, mobile open/closed)
- Responsive layout with Tailwind CSS
- Smooth transitions between states
- Footer with links

**Props:** None (root layout)

**Events Emitted:** None

**State:**
```typescript
interface SidebarState {
  isOpen: boolean      // Mobile sidebar overlay state
  isCollapsed: boolean // Desktop sidebar collapse state
}
```

**Methods:**
- `toggleSidebar()` - Toggles desktop sidebar collapse
- `toggleMobileSidebar()` - Toggles mobile sidebar overlay
- `closeMobileSidebar()` - Closes mobile sidebar
- `handleLogout()` - Handles logout action (placeholder)

**Usage:**
```vue
<!-- Automatically rendered by router -->
<router-view />
```

---

### 2. Header.vue

**Location:** `/src/layouts/components/Header.vue`

**Purpose:** Top navigation bar with user info and controls.

**Features:**
- Mobile hamburger menu button
- Desktop sidebar toggle button
- Page title/breadcrumb area
- Notification icon with badge
- User profile display (name, role, avatar)
- Responsive design

**Props:**
```typescript
interface Props {
  isCollapsed: boolean // Sidebar collapse state
}
```

**Events Emitted:**
```typescript
{
  toggleSidebar: []        // Desktop toggle clicked
  toggleMobileSidebar: []  // Mobile hamburger clicked
}
```

**User Data:**
Currently uses mock data. In production, integrate with auth store:
```typescript
interface UserInfo {
  id: string
  name: string
  email: string
  avatar?: string
  role: string
}
```

---

### 3. Sidebar.vue

**Location:** `/src/layouts/components/Sidebar.vue`

**Purpose:** Navigation sidebar with menu items.

**Features:**
- Navigation menu with icons
- Active route highlighting
- Logout button
- Smooth animations and transitions
- Mobile overlay backdrop
- Responsive width adjustment

**Props:**
```typescript
interface Props {
  isOpen: boolean      // Mobile overlay state
  isCollapsed: boolean // Desktop collapse state
}
```

**Events Emitted:**
```typescript
{
  close: []   // Close mobile sidebar
  logout: []  // Logout clicked
}
```

**Menu Items:**
Defined in component as array of `MenuItem`:
```typescript
interface MenuItem {
  id: string
  label: string
  icon: string     // SVG path data
  path: string     // Route path
}
```

Current menu:
- Dashboard (`/admin/dashboard`)
- Users (`/admin/users`)
- Settings (`/admin/settings`)

**Customization:**
To add new menu items, modify the `menuItems` array:
```typescript
const menuItems: MenuItem[] = [
  {
    id: 'new-section',
    label: 'New Section',
    icon: 'M4 6h16M4 12h16M4 18h16', // SVG path
    path: '/admin/new-section'
  },
  // ... existing items
]
```

---

### 4. View Components

#### DashboardView.vue

**Location:** `/src/views/DashboardView.vue`

**Purpose:** Admin dashboard with statistics and activity.

**Features:**
- 4 stat cards with metrics (users, sessions, revenue, orders)
- Recent users list
- Activity log
- Responsive grid layout

**Customization:**
Modify the `stats` array to change displayed metrics:
```typescript
interface StatCard {
  title: string
  value: string
  change: string
  changeType: 'increase' | 'decrease'
  icon: string
}
```

#### UsersView.vue

**Location:** `/src/views/UsersView.vue`

**Purpose:** User management interface.

**Features:**
- User table with sample data
- Search bar
- Role and status filters
- Add user button
- Edit/Delete actions
- Responsive table

**Data Structure:**
```typescript
interface User {
  id: string
  name: string
  email: string
  role: string
  status: 'active' | 'inactive'
  createdAt: string
}
```

#### SettingsView.vue

**Location:** `/src/views/SettingsView.vue`

**Purpose:** Application settings management.

**Features:**
- Settings navigation sidebar
- General settings form (site name, email, timezone)
- Notification preferences (email, push, SMS)
- Maintenance mode toggle
- Save changes button

**Data Structure:**
```typescript
interface FormData {
  siteName: string
  siteEmail: string
  notifications: {
    email: boolean
    push: boolean
    sms: boolean
  }
  maintenance: boolean
  timezone: string
}
```

---

## Router Configuration

**Location:** `/src/router/index.ts`

### Route Structure

```typescript
/ → /admin/dashboard (redirect)
/admin
  ├── /dashboard → DashboardView
  ├── /users → UsersView
  └── /settings → SettingsView
/:pathMatch(.*) → /admin/dashboard (catch-all)
```

### Route Meta

Each route includes metadata:
```typescript
meta: {
  title: 'Page Title',
  requiresAuth: true
}
```

### Navigation Guard

Authentication guard is implemented (currently mocked):
```typescript
router.beforeEach((to, _from, next) => {
  const requiresAuth = to.matched.some(record => record.meta.requiresAuth)
  const isAuthenticated = true // Replace with actual auth check

  if (requiresAuth && !isAuthenticated) {
    // Redirect to login
    next('/login')
  } else {
    next()
  }
})
```

### Adding New Routes

1. Create view component in `/src/views/`
2. Add route definition:
```typescript
{
  path: 'new-page',
  name: 'new-page',
  component: () => import('@/views/NewPageView.vue'),
  meta: {
    title: 'New Page',
    requiresAuth: true
  }
}
```
3. Add menu item to Sidebar.vue

---

## Styling Guide

### Tailwind CSS Classes

#### Layout
- Sidebar: `w-64` (expanded) / `w-16` (collapsed)
- Header: `h-16` fixed height
- Content: `flex-1` to fill remaining space

#### Colors
- Primary: `indigo-600`
- Background: `gray-50`, `gray-800`
- Text: `gray-800`, `gray-600`, `gray-500`
- Borders: `gray-200`, `gray-700`

#### Responsive Breakpoints
- Mobile: `< 1024px`
- Desktop: `≥ 1024px` (`lg:` prefix)

#### Transitions
```css
transition-all duration-300 ease-in-out
```

### Custom Styling

To customize colors, edit Tailwind config or use CSS variables:
```css
/* In assets/main.css */
:root {
  --primary-color: #4f46e5; /* indigo-600 */
  --sidebar-bg: #1f2937;     /* gray-800 */
}
```

---

## State Management

### Current Implementation

State is managed locally in components using Vue 3 Composition API:
- Sidebar state: `AdminLayout.vue`
- User data: Mock data in `Header.vue`
- Form data: Local refs in view components

### Recommended Enhancement

For larger applications, consider centralizing state with Pinia:

```typescript
// stores/auth.ts
import { defineStore } from 'pinia'

export const useAuthStore = defineStore('auth', () => {
  const user = ref<UserInfo | null>(null)
  const isAuthenticated = computed(() => !!user.value)

  const login = async (credentials: LoginCredentials) => {
    // API call
  }

  const logout = async () => {
    // API call
    user.value = null
  }

  return { user, isAuthenticated, login, logout }
})
```

---

## Responsive Design

### Desktop (≥1024px)

- Sidebar always visible
- Can collapse to icon-only (64px)
- Content area adjusts margin: `ml-64` or `ml-16`
- Desktop toggle button in header

### Mobile (<1024px)

- Sidebar hidden by default
- Hamburger menu in header
- Sidebar opens as fixed overlay
- Dark backdrop (50% opacity)
- Closes on outside click or navigation

### Breakpoint Reference

```
sm:  640px
md:  768px
lg:  1024px ← Main breakpoint
xl:  1280px
2xl: 1536px
```

---

## TypeScript Types

### Navigation Types

**Location:** `/src/types/navigation.ts`

```typescript
export interface MenuItem {
  id: string
  label: string
  icon: string
  path: string
  requiresAuth?: boolean
  isAction?: boolean
}

export interface UserInfo {
  id: string
  name: string
  email: string
  avatar?: string
  role: string
}

export interface SidebarState {
  isOpen: boolean
  isCollapsed: boolean
}
```

### Adding Custom Types

Create additional type files in `/src/types/`:
```typescript
// types/settings.ts
export interface AppSettings {
  // ...
}
```

Import with path alias:
```typescript
import type { AppSettings } from '@/types/settings'
```

---

## Development

### Setup

```bash
cd frontend
npm install
```

### Development Server

```bash
npm run dev
```
Runs on `http://localhost:3000/`

### Type Checking

```bash
npm run type-check
```

### Production Build

```bash
npm run build
```
Output: `dist/` directory

### Preview Production Build

```bash
npm run preview
```

---

## Integration Guide

### Adding Authentication

1. **Create auth store** (Pinia recommended):
```typescript
// stores/auth.ts
export const useAuthStore = defineStore('auth', {
  state: () => ({
    user: null,
    token: null
  }),
  actions: {
    async login(credentials) { /* ... */ },
    async logout() { /* ... */ }
  }
})
```

2. **Update router guard**:
```typescript
router.beforeEach((to, _from, next) => {
  const authStore = useAuthStore()
  if (to.meta.requiresAuth && !authStore.isAuthenticated) {
    next('/login')
  } else {
    next()
  }
})
```

3. **Update Header.vue**:
```typescript
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()
const user = computed(() => authStore.user)
```

4. **Update AdminLayout logout**:
```typescript
const handleLogout = async () => {
  const authStore = useAuthStore()
  await authStore.logout()
  router.push('/login')
}
```

### Adding API Integration

1. **Create API client** (axios/fetch):
```typescript
// api/client.ts
import axios from 'axios'

export const apiClient = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL,
  headers: {
    'Content-Type': 'application/json'
  }
})
```

2. **Create API services**:
```typescript
// api/users.ts
export const usersApi = {
  getAll: () => apiClient.get('/users'),
  getById: (id: string) => apiClient.get(`/users/${id}`),
  create: (data: CreateUserDto) => apiClient.post('/users', data)
}
```

3. **Use in components**:
```typescript
const fetchUsers = async () => {
  try {
    const { data } = await usersApi.getAll()
    users.value = data
  } catch (error) {
    console.error('Failed to fetch users:', error)
  }
}
```

### Adding Notifications

Recommended: **Vue Toastification**

```bash
npm install vue-toastification@next
```

```typescript
// main.ts
import Toast from 'vue-toastification'
import 'vue-toastification/dist/index.css'

app.use(Toast)
```

```typescript
// In components
import { useToast } from 'vue-toastification'

const toast = useToast()
toast.success('Settings saved successfully!')
toast.error('Failed to save settings')
```

---

## Customization Examples

### Changing Sidebar Width

```vue
<!-- Sidebar.vue -->
<template>
  <aside :class="[
    // Change these values:
    props.isCollapsed ? 'w-20' : 'w-72',
  ]">
```

```vue
<!-- AdminLayout.vue -->
<div :class="[
  // Match sidebar widths:
  sidebarState.isCollapsed ? 'lg:ml-20' : 'lg:ml-72'
]">
```

### Adding User Dropdown Menu

```vue
<!-- Header.vue -->
<div class="relative">
  <button @click="showDropdown = !showDropdown">
    <!-- User avatar -->
  </button>
  <div v-if="showDropdown" class="absolute right-0 mt-2 w-48 bg-white rounded-lg shadow-lg">
    <a href="#" class="block px-4 py-2 hover:bg-gray-100">Profile</a>
    <a href="#" class="block px-4 py-2 hover:bg-gray-100">Settings</a>
    <a href="#" class="block px-4 py-2 hover:bg-gray-100">Logout</a>
  </div>
</div>
```

### Adding Breadcrumbs

```vue
<!-- Header.vue -->
<nav class="flex" aria-label="Breadcrumb">
  <ol class="flex items-center space-x-2">
    <li v-for="(crumb, i) in breadcrumbs" :key="i">
      <router-link :to="crumb.path">{{ crumb.label }}</router-link>
      <span v-if="i < breadcrumbs.length - 1">/</span>
    </li>
  </ol>
</nav>
```

---

## Performance Optimization

### Lazy Loading Routes

```typescript
// router/index.ts
const routes = [
  {
    path: 'dashboard',
    component: () => import('@/views/DashboardView.vue')
  }
]
```

### Component Code Splitting

```typescript
// AdminLayout.vue
const Header = defineAsyncComponent(() => import('./components/Header.vue'))
const Sidebar = defineAsyncComponent(() => import('./components/Sidebar.vue'))
```

### Tailwind Purging

Already configured in `tailwind.config.js`:
```javascript
content: [
  './index.html',
  './src/**/*.{vue,js,ts,jsx,tsx}'
]
```

---

## Troubleshooting

### Sidebar not appearing on mobile

**Issue:** Sidebar hidden behind content
**Solution:** Check z-index values:
```css
sidebar: z-30
overlay: z-20
content: z-10
```

### Routes not working

**Issue:** Router not initialized
**Solution:** Verify `main.ts`:
```typescript
import router from './router'
app.use(router)
```

### TypeScript errors with @/ imports

**Issue:** Path alias not recognized
**Solution:** Check `tsconfig.json` and `vite.config.ts`:
```json
// tsconfig.json
"paths": {
  "@/*": ["./src/*"]
}
```
```typescript
// vite.config.ts
resolve: {
  alias: {
    '@': fileURLToPath(new URL('./src', import.meta.url))
  }
}
```

### Tailwind styles not applied

**Issue:** PostCSS not configured
**Solution:** Verify `assets/main.css`:
```css
@tailwind base;
@tailwind components;
@tailwind utilities;
```

---

## Best Practices

### Component Structure

1. Group related components in subdirectories
2. Use clear, descriptive names
3. Keep components focused and single-purpose
4. Extract reusable logic to composables

### TypeScript

1. Always define interfaces for props and state
2. Use type inference where possible
3. Avoid `any` types
4. Enable strict mode

### Styling

1. Use Tailwind utilities first
2. Create custom components for repeated patterns
3. Follow mobile-first approach
4. Maintain consistent spacing

### State Management

1. Keep state as local as possible
2. Use stores for global/shared state
3. Avoid prop drilling (use provide/inject or stores)
4. Separate business logic from UI logic

---

## Future Enhancements

### Recommended Features

1. **Dark Mode**
   - Add theme toggle
   - Store preference in localStorage
   - Use Tailwind dark: variants

2. **Multi-language Support**
   - Integrate vue-i18n
   - Add language switcher
   - Externalize strings

3. **Advanced Permissions**
   - Role-based access control
   - Route-level permissions
   - Component-level visibility

4. **Search Functionality**
   - Global search in header
   - Command palette (Cmd+K)
   - Fuzzy search

5. **Animations**
   - Page transitions
   - Loading states
   - Skeleton screens

6. **Accessibility**
   - ARIA labels
   - Keyboard navigation
   - Focus management
   - Screen reader support

---

## Resources

### Documentation

- [Vue 3 Documentation](https://vuejs.org/)
- [TypeScript Documentation](https://www.typescriptlang.org/)
- [Tailwind CSS Documentation](https://tailwindcss.com/)
- [Vue Router Documentation](https://router.vuejs.org/)

### Icons

Current implementation uses inline SVG paths. Consider:
- [Heroicons](https://heroicons.com/) - Used in examples
- [Font Awesome](https://fontawesome.com/)
- [Material Icons](https://fonts.google.com/icons)

### Libraries

Recommended for production:
- **Pinia** - State management
- **VueUse** - Composition utilities
- **Vue Toastification** - Notifications
- **VeeValidate** - Form validation
- **Axios** - HTTP client

---

## Support

For questions or issues:
1. Check this documentation
2. Review component source code
3. Consult Vue 3 / TypeScript documentation
4. Check project README.md

---

## Changelog

### Version 1.0.0 (2026-01-05)

**Initial Release**
- AdminLayout component with sidebar and header
- Three view components (Dashboard, Users, Settings)
- Vue Router configuration
- Full TypeScript typing
- Responsive design (mobile + desktop)
- Sidebar toggle and mobile overlay
- Tailwind CSS styling

---

**Last Updated:** 2026-01-05
**Version:** 1.0.0
**Author:** Admin Layout Implementation Team
