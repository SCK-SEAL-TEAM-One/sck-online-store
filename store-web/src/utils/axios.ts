import { useUserStore } from '@/hooks/use-user-store'
import axios from 'axios'

// ----------------------------------------------------------------------------

const axiosShoppingMallApi = axios.create({
  baseURL: process.env.storeServiceURL || 'http://localhost:3000',
  headers: {
    'Accept-Language': 'en'
  }
})

axiosShoppingMallApi.interceptors.request.use(
  (config) => {
    const userId = useUserStore.getState().userId

    if (userId) {
      config.headers['uid'] = userId
    }

    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

export default axiosShoppingMallApi
