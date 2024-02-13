import SHIPPING_METHOD from '@/assets/data/shipping_method.json'
import GetProductInCartService, {
  ProductDetailInCart
} from '@/services/cart/get-product-list'
import * as calculate from '@/utils/total-price'
import type {} from '@redux-devtools/extension' // required for devtools typing
import { produce } from 'immer'
import { create } from 'zustand'
import { devtools, persist } from 'zustand/middleware'

// ---------------------------------------------------

type PaymentCreditInformationType = {
  name: string
  number: string
  expiry: string
  cvv: string
  issuer: string
  focused: string
}

type ShippingInformationType = {
  firstName: string
  lastName: string
  address: string
  mobileNumber: string
  provinceId: number
  districtId: number
  subDistrictId: number
  provinceName: string
  districtName: string
  subDistrictName: string
  zipCode: number
  focused: string
}

type ShippingType = {
  shippingMethod: number
  shippingFee: number
  shippingInformation: ShippingInformationType
}

type PointType = {
  point: number // Current Point
  isUsePoint: boolean
  burnPoint: number
}

type PaymentType = {
  paymentMethod: string
  paymentCreditInformation: PaymentCreditInformationType
}

type OrderStoreType = {
  cart: ProductDetailInCart[]
  totalProduct: number
  subTotal: number
  totalPayment: number
  receivePoint: number
  shipping: ShippingType
  point: PointType
  payment: PaymentType
  getProductListInCart: () => void
  setPoint: (point: number) => void
  setIsUsePoint: (isUsePoint: boolean) => void
  setPaymentMethod: (paymentMethod: string) => void
  setPaymentInformation: (
    paymentCreditInformation: PaymentCreditInformationType
  ) => void
  setShippingMethod: (shippingMethod: number, shippingFee: number) => void
  setShippingInformation: (shippingInformation: ShippingInformationType) => void
  updateSummary: () => void
}

const useOrderStore = create<OrderStoreType>()(
  devtools(
    persist(
      (set, get) => ({
        cart: [],
        totalProduct: 0,
        subTotal: 0,
        totalPayment: 0,
        receivePoint: 0,
        shipping: {
          shippingMethod: SHIPPING_METHOD[0].id,
          shippingFee: SHIPPING_METHOD[0].price,
          shippingInformation: {
            firstName: '',
            lastName: '',
            address: '',
            mobileNumber: '',
            provinceId: 0,
            districtId: 0,
            subDistrictId: 0,
            provinceName: '',
            districtName: '',
            subDistrictName: '',
            zipCode: 0,
            focused: ''
          }
        },
        point: {
          point: 0,
          burnPoint: 0,
          isUsePoint: false
        },
        payment: {
          paymentMethod: 'credit',
          paymentCreditInformation: {
            number: '',
            name: '',
            expiry: '',
            cvv: '',
            issuer: '',
            focused: ''
          }
        },
        getProductListInCart: async () => {
          // Mock userId
          const userId = 1

          const productInCart = await GetProductInCartService(userId)
          const price = productInCart?.map((item) => item.product_price)
          const total = calculate.subTotal(price)

          set(
            produce((state) => {
              state.totalProduct = productInCart.length
              state.cart = productInCart
              state.subTotal = total
              state.totalPayment = total
            })
          )

          get().updateSummary()
        },
        setPoint(point: number) {
          set(
            produce((state) => {
              state.point.point = point
            })
          )
        },
        setIsUsePoint: (isUsePoint: boolean) => {
          set(
            produce((state) => {
              state.point.isUsePoint = isUsePoint
            })
          )

          get().updateSummary()
        },
        setPaymentMethod: (paymentMethod: string) => {
          set(
            produce((state) => {
              state.payment.paymentMethod = paymentMethod
            })
          )
        },
        setPaymentInformation: (
          paymentCreditInformation: PaymentCreditInformationType
        ) => {
          set(
            produce((state) => {
              state.payment.paymentCreditInformation = paymentCreditInformation
            })
          )
        },
        setShippingMethod: (shippingMethod: number, shippingFee: number) => {
          const newTotalPayment = get().subTotal + shippingFee

          set(
            produce((state) => {
              state.totalPayment = newTotalPayment
              state.shipping.shippingMethod = shippingMethod
              state.shipping.shippingFee = shippingFee
            })
          )

          get().updateSummary()
        },
        setShippingInformation: (
          shippingInformation: ShippingInformationType
        ) => {
          set(
            produce((state) => {
              state.shipping.shippingInformation = shippingInformation
            })
          )
        },
        updateSummary: async () => {
          const isUsePoint = get().point.isUsePoint
          const point = get().point.point

          const subTotal = get().subTotal
          const shippingFee = get().shipping.shippingFee

          // Calculate Point
          const pointsUsed = isUsePoint
            ? calculate.pointBurn(point, subTotal)
            : 0

          // Total Payment
          const totalPayment = calculate.totalPayment(
            isUsePoint,
            pointsUsed,
            subTotal,
            shippingFee
          )

          // Point Receive
          const receivePoint = calculate.receiptPoint(totalPayment)

          set(
            produce((state) => {
              state.totalPayment = totalPayment
              state.receivePoint = receivePoint
              state.point.burnPoint = pointsUsed
            })
          )
        }
      }),
      {
        name: 'order-storage'
      }
    )
  )
)

export default useOrderStore
