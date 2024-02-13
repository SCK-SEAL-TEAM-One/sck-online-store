import axios from 'axios'

// ----------------------------------------------------------------------------

const axiosShoppingMallApi = axios.create({
  baseURL: 'https://localhost:3000',
  headers: {
    'Accept-Language': 'en'
  }
})

export default axiosShoppingMallApi
