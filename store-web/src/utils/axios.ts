import axios from 'axios'

// ----------------------------------------------------------------------------

const axiosShoppingMallApi = axios.create({
  baseURL: process.env.storeServiceURL || 'http://localhost:3000',
  headers: {
    'Accept-Language': 'en'
  }
})

export default axiosShoppingMallApi
