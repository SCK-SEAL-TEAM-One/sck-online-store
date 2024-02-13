# Store-web

Web Application

## Getting Started

เริ่มต้น run development server:

```bash
npm run dev
# หรือ
yarn dev
# หรือ
pnpm dev
# หรือ
bun dev
```

หลังจาก run ให้เข้าไปที่ [http://localhost:3000](http://localhost:3000)

## Environment

อยู่ที่ `next.config.js` ในส่วนของ `env`

```js
const nextConfig = {
  // ...
  env: {
    storeServiceURL: 'https://localhost:3000'
  }
}

module.exports = nextConfig
```
