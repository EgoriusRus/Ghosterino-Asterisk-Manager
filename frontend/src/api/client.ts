import type {
  Profile,
  ProfileWithLocation,
  Device,
  Location,
  PaginationParams,
  PaginatedResult
} from '@/types/api'

// API base URL - can be configured via environment variables
const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:3000/api'

// Generic fetch wrapper with error handling
async function fetchAPI<T>(endpoint: string, options?: RequestInit): Promise<T> {
  const url = `${API_BASE_URL}${endpoint}`

  try {
    const response = await fetch(url, {
      ...options,
      headers: {
        'Content-Type': 'application/json',
        ...options?.headers,
      },
    })

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`)
    }

    return await response.json()
  } catch (error) {
    console.error('API request failed:', error)
    throw error
  }
}

// Profiles API
export const profilesAPI = {
  /**
   * Get all profiles with pagination
   */
  getAll(params: PaginationParams): Promise<PaginatedResult<ProfileWithLocation>> {
    const queryParams = new URLSearchParams({
      page: params.page.toString(),
      perPage: params.perPage.toString(),
    })
    return fetchAPI<PaginatedResult<ProfileWithLocation>>(`/profiles?${queryParams}`)
  },

  /**
   * Get a single profile by ID
   */
  getById(id: number): Promise<Profile> {
    return fetchAPI<Profile>(`/profiles/${id}`)
  },

  /**
   * Create a new profile
   */
  create(profile: Omit<Profile, 'id' | 'createdAt' | 'updatedAt'>): Promise<Profile> {
    return fetchAPI<Profile>('/profiles', {
      method: 'POST',
      body: JSON.stringify(profile),
    })
  },

  /**
   * Update an existing profile
   */
  update(id: number, profile: Partial<Profile>): Promise<Profile> {
    return fetchAPI<Profile>(`/profiles/${id}`, {
      method: 'PUT',
      body: JSON.stringify(profile),
    })
  },

  /**
   * Delete a profile
   */
  delete(id: number): Promise<void> {
    return fetchAPI<void>(`/profiles/${id}`, {
      method: 'DELETE',
    })
  },
}

// Devices API
export const devicesAPI = {
  /**
   * Get all devices (no pagination)
   */
  getAll(): Promise<Device[]> {
    return fetchAPI<Device[]>('/devices')
  },

  /**
   * Get a single device by MAC address
   */
  getByMac(mac: string): Promise<Device> {
    return fetchAPI<Device>(`/devices/${mac}`)
  },

  /**
   * Create a new device
   */
  create(device: Omit<Device, 'createdAt' | 'updatedAt'>): Promise<Device> {
    return fetchAPI<Device>('/devices', {
      method: 'POST',
      body: JSON.stringify(device),
    })
  },

  /**
   * Update an existing device
   */
  update(mac: string, device: Partial<Device>): Promise<Device> {
    return fetchAPI<Device>(`/devices/${mac}`, {
      method: 'PUT',
      body: JSON.stringify(device),
    })
  },

  /**
   * Delete a device
   */
  delete(mac: string): Promise<void> {
    return fetchAPI<void>(`/devices/${mac}`, {
      method: 'DELETE',
    })
  },
}

// Locations API
export const locationsAPI = {
  /**
   * Get all locations (no pagination)
   */
  getAll(): Promise<Location[]> {
    return fetchAPI<Location[]>('/locations')
  },

  /**
   * Get a single location by ID
   */
  getById(id: number): Promise<Location> {
    return fetchAPI<Location>(`/locations/${id}`)
  },

  /**
   * Create a new location
   */
  create(location: Omit<Location, 'id' | 'createdAt' | 'updatedAt'>): Promise<Location> {
    return fetchAPI<Location>('/locations', {
      method: 'POST',
      body: JSON.stringify(location),
    })
  },

  /**
   * Update an existing location
   */
  update(id: number, location: Partial<Location>): Promise<Location> {
    return fetchAPI<Location>(`/locations/${id}`, {
      method: 'PUT',
      body: JSON.stringify(location),
    })
  },

  /**
   * Delete a location
   */
  delete(id: number): Promise<void> {
    return fetchAPI<void>(`/locations/${id}`, {
      method: 'DELETE',
    })
  },
}
