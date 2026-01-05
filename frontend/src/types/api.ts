// API Types matching backend Go models

// Device Models (enum from backend)
export type DeviceModel = 'Yealink T27G' | 'Yealink T23G' | 'Fanvil' | 'Cisco'

// Device represents an IP phone or device
export interface Device {
  mac: string
  deviceModel: DeviceModel
  createdAt: string
  updatedAt: string
}

// Location represents a location (building/address) with network settings
export interface Location {
  id: number
  name: string
  server: string // IP address
  subnet: string // CIDR notation
  voipVlan: number
  vlan: number
  createdAt: string
  updatedAt: string
}

// Profile represents an employee profile with SIP settings
export interface Profile {
  id: number
  name: string
  email: string
  device: string | null // MAC address
  locationId: number | null
  internalNumber: number
  externalNumber: string
  ringGroup: number | null
  pickupGroup: number | null
  isActive: boolean
  createdAt: string
  updatedAt: string
}

// ProfileWithLocation represents a profile with location data
export interface ProfileWithLocation extends Profile {
  locationName: string | null
  server: string | null
  subnet: string | null
  voipVlan: number | null
  vlan: number | null
}

// Pagination types
export interface PaginationParams {
  page: number
  perPage: number
}

export interface PaginationResponse {
  total: number
  page: number
  perPage: number
  pages: number
}

export interface PaginatedResult<T> {
  data: T[]
  pagination: PaginationResponse
}
