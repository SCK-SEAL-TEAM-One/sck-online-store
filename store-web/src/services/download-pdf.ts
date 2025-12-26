import axios from 'axios'

const downloadFileAxiosInstance = axios.create({
  baseURL: process.env.storeServiceURL || 'http://localhost:3000',
  headers: {
    Accept: 'application/pdf'
  },
  withCredentials: true
})

downloadFileAxiosInstance.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('accessToken')
    if (token) {
      config.headers['Authorization'] = `Bearer ${token}`
    }

    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

export const getOrderSummary = async (orderId: number) => {
  try {
    const { data, headers } = await downloadFileAxiosInstance.post(
      `/api/v1/order/${orderId}/summary`,
      {
        responseType: 'blob'
      }
    )

    const disposition = headers['content-disposition']
    let fileName = `Order_Summary_${orderId}.pdf`

    if (disposition) {
      const match = disposition.match(/filename="?([^";]+)"?/)
      if (match && match[1]) {
        fileName = match[1]
      }
    }

    const blob = new Blob([data], { type: 'application/pdf' })
    return { blob, fileName }
  } catch (error) {
    return null
  }
}
