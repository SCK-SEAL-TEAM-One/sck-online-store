import GetProductInCartService, {
  ProductDetailInCart
} from '@/services/cart/get-product-list'
import * as calculate from '@/utils/total-price'
import type {} from '@redux-devtools/extension' // required for devtools typing
import { create } from 'zustand'
import { devtools, persist } from 'zustand/middleware'

// ---------------------------------------------------

export type ProductToCartProps = {
  product_id: number
  quantity: number
}

type CartStoreType = {
  cart: ProductDetailInCart[]
  totalProduct: number
  subTotal: number
  getProductListInCart: () => void
  // addToCartLocal: (product: ProductToCartProps) => void
}

const useCartStore = create<CartStoreType>()(
  devtools(
    persist(
      (set, get) => ({
        cart: [],
        totalProduct: 0,
        subTotal: 0,
        getProductListInCart: async () => {
          // Mock userId
          const userId = 1

          const productInCart = await GetProductInCartService(userId)
          const price = productInCart?.map((item) => item.product_price)
          const total = calculate.subTotal(price)

          set({
            cart: productInCart,
            totalProduct: productInCart.length,
            subTotal: total
          })
        }
        // addToCartLocal: (product: ProductToCartProps) => {
        //   set((state) => {
        //     if (state.cart.length > 0) {
        //       const cart = state.cart
        //       // Find index of product in cart
        //       const isProductItem = cart.findIndex(
        //         (item) => item.product_id === product.product_id
        //       )
        //       // Update quantity or add new product
        //       if (isProductItem > -1) {
        //         cart[isProductItem].product_id = product.product_id
        //         cart[isProductItem].quantity = product.quantity
        //       } else {
        //         cart.push(product)
        //       }
        //       return {
        //         cart: cart
        //       }
        //     } else {
        //       return {
        //         cart: [product]
        //       }
        //     }
        //   })
        // }
      }),
      {
        name: 'cart-storage'
      }
    )
  )
)

export default useCartStore
