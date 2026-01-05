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
