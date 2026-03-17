import http from "k6/http"
import { check, sleep } from "k6"

const BASE_URL = __ENV.BASE_URL || "http://localhost"
const PASSWORD = "P@ssw0rd"

export const options = {
  stages: [
    { duration: "30s", target: 20 },
    { duration: "2m30s", target: 20 },
    { duration: "30s", target: 0 },
  ],
  thresholds: {
    http_req_duration: ["p(99)<200"],
    http_req_failed: ["rate==0"],
  },
}

function randomInt(min, max) {
  return Math.floor(Math.random() * (max - min + 1)) + min
}

function thinkTime() {
  sleep(randomInt(1, 3))
}

function authHeaders(token) {
  return {
    headers: {
      Authorization: `Bearer ${token}`,
      "Content-Type": "application/json",
    },
  }
}

export default function () {
  const userId = (__VU % 84) + 1
  const username = `user_${userId}`

  // Step 1: Login
  const loginRes = http.post(
    `${BASE_URL}/api/v1/login`,
    JSON.stringify({ username, password: PASSWORD }),
    { headers: { "Content-Type": "application/json" } }
  )

  check(loginRes, {
    "login: status 200": (r) => r.status === 200,
    "login: has access_token": (r) => r.json("access_token") !== undefined,
  })

  if (loginRes.status !== 200) return

  const token = loginRes.json("access_token")
  const opts = authHeaders(token)

  thinkTime()

  // Step 2: Browse products
  const browseRes = http.get(
    `${BASE_URL}/api/v1/product?q=&offset=0&limit=20`,
    opts
  )

  check(browseRes, {
    "browse: status 200": (r) => r.status === 200,
  })

  thinkTime()

  // Step 3: View a random product (skip ID 7 — intentional fault injection)
  const availableProducts = [1, 2, 3, 4, 5, 6, 8, 9]
  const productId = availableProducts[randomInt(0, availableProducts.length - 1)]
  const productRes = http.get(
    `${BASE_URL}/api/v1/product/${productId}`,
    opts
  )

  check(productRes, {
    "view product: status 200": (r) => r.status === 200,
    "view product: has product_name": (r) => r.json("product_name") !== undefined,
  })

  thinkTime()

  // Step 4: Add to cart
  const addCartRes = http.put(
    `${BASE_URL}/api/v1/addCart`,
    JSON.stringify({ product_id: productId, quantity: 1 }),
    opts
  )

  check(addCartRes, {
    "add cart: status 200": (r) => r.status === 200,
  })

  thinkTime()

  // Step 5: Get cart
  const cartRes = http.get(`${BASE_URL}/api/v1/cart`, opts)

  check(cartRes, {
    "get cart: status 200": (r) => r.status === 200,
    "get cart: has carts": (r) => r.json("carts") !== undefined,
    "get cart: has summary": (r) => r.json("summary") !== undefined,
  })

  if (cartRes.status !== 200) return

  const cartData = cartRes.json()
  const cartItems = cartData.carts || []
  const summary = cartData.summary || {}

  if (cartItems.length === 0) return

  thinkTime()

  // Step 6: Submit order
  const shippingMethodId = randomInt(1, 3)
  const paymentMethodId = randomInt(1, 2)

  const orderBody = {
    cart: cartItems.map((item) => ({
      product_id: item.product_id,
      quantity: item.quantity,
    })),
    shipping_method_id: shippingMethodId,
    shipping_address: "123 Load Test Street",
    shipping_sub_disterict: "Klongtoey",
    shipping_district: "Klongtoey",
    shipping_province: "Bangkok",
    shipping_zip_code: "10110",
    recipient_first_name: "Load",
    recipient_last_name: `Test ${userId}`,
    recipient_phone_number: "0812345678",
    payment_method_id: paymentMethodId,
    sub_total_price: summary.total_price_thb || 0,
    discount_price: 0,
    total_price: summary.total_price_thb || 0,
    burn_point: 0,
  }

  const orderRes = http.post(
    `${BASE_URL}/api/v1/order`,
    JSON.stringify(orderBody),
    opts
  )

  check(orderRes, {
    "submit order: status 200": (r) => r.status === 200,
    "submit order: has order_number": (r) => r.json("order_number") !== undefined,
  })

  if (orderRes.status !== 200) return

  const orderNumber = orderRes.json("order_number")

  thinkTime()

  // Step 7: Confirm payment
  const paymentRes = http.post(
    `${BASE_URL}/api/v1/confirmPayment`,
    JSON.stringify({
      order_number: orderNumber,
      otp: 123456,
      ref_otp: "AXYZ",
    }),
    opts
  )

  check(paymentRes, {
    "confirm payment: status 200": (r) => r.status === 200,
    "confirm payment: has tracking_number": (r) => r.json("tracking_number") !== undefined,
  })

  thinkTime()

  // Step 8: Order summary
  const summaryOpts = {
    headers: {
      Authorization: `Bearer ${token}`,
      "Content-Type": "application/json",
      Accept: "application/json",
    },
  }

  const summaryRes = http.post(
    `${BASE_URL}/api/v1/order/${orderNumber}/summary`,
    null,
    summaryOpts
  )

  check(summaryRes, {
    "order summary: status 200": (r) => r.status === 200,
    "order summary: has order_number": (r) => r.json("order_number") !== undefined,
    "order summary: has tracking_no": (r) => r.json("tracking_no") !== undefined,
  })

  thinkTime()
}
