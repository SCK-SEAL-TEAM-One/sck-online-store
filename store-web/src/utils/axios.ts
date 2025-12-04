import { useUserStore } from '@/hooks/use-user-store'
import { RefreshToken } from '@/services/auth'
import axios from 'axios'

// ----------------------------------------------------------------------------

const axiosShoppingMallApi = axios.create({
  baseURL: process.env.storeServiceURL || 'http://localhost:3000',
  headers: {
    'Accept-Language': 'en'
  },
  withCredentials: true
})

axiosShoppingMallApi.interceptors.request.use(
  (config) => {
    const userId = useUserStore.getState().userId
    const token = localStorage.getItem('accessToken')
    if (token) {
      config.headers['Authorization'] = `Bearer ${token}`
    }

    if (userId) {
      config.headers['uid'] = userId
    }

    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

let isRefreshing = false
let queue: Array<(token: string) => void> = []

function addToQueue(callback: (token: string) => void) {
  queue.push(callback)
}

function processQueue(newToken: string) {
  queue.forEach((cb) => cb(newToken))
  queue = []
}

axiosShoppingMallApi.interceptors.response.use(
  (response) => {
    return response
  },
  async (error) => {
    const originalRequest = error.config
    if (!originalRequest) return Promise.reject(error)

    if (error.response?.status === 401 && !originalRequest._retry) {
      // Wait to retry if the token is already refreshed
      if (isRefreshing) {
        return new Promise((resolve) => {
          addToQueue((token) => {
            originalRequest.headers['Authorization'] = `Bearer ${token}`
            resolve(axiosShoppingMallApi(originalRequest))
          })
        })
      }
      originalRequest._retry = true
      isRefreshing = true

      try {
        const response = await RefreshToken()
        if (!response?.data) {
          window.location.href = '/auth/login'
          return
        }
        const { accessToken } = response.data
        localStorage.setItem('accessToken', accessToken)

        axiosShoppingMallApi.defaults.headers.common['Authorization'] =
          `Bearer ${accessToken}`

        originalRequest.headers = {
          ...originalRequest.headers,
          Authorization: `Bearer ${accessToken}`
        }

        // Run the queue of previously unauthorized requests
        processQueue(accessToken)

        return axiosShoppingMallApi(originalRequest)
      } catch (error) {
        window.location.href = '/auth/login'
        return Promise.reject(error)
      }
    }
    return Promise.reject(error)
  }
)

export default axiosShoppingMallApi
