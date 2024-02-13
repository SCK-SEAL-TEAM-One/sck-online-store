/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  images: {
    domains: ['tailwindui.com', 'localhost']
  },
  env: {
    storeServiceURL: 'https://localhost:8000',
  }
}

module.exports = nextConfig
