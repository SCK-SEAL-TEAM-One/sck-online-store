import { useUserStore } from '@/hooks/use-user-store'
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

axiosShoppingMallApi.interceptors.response.use(
  (response) => {
    return response
  },
  (error) => {
    if (error.response?.status === 401) {
      window.location.href = '/auth/login'
    }
    return Promise.reject(error)
  }
)

export default axiosShoppingMallApi
