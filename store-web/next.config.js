/** @type {import('next').NextConfig} */
const STORE_SERVICE_URL =
  process.env.STORE_SERVICE_URL || 'http://localhost:8000'

const nextConfig = {
  reactStrictMode: true,
  images: {
    domains: ['tailwindui.com', 'localhost']
  },
  env: {
    storeServiceURL: STORE_SERVICE_URL
  }
}

module.exports = nextConfig
