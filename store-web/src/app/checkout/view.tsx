'use client'

import OrderList from '@/app/checkout/components/order-list'
import OrderSummary from '@/app/checkout/components/order-summary'
import PaymentMethod from '@/app/checkout/components/payment-method'
import ShippingInfomation from '@/app/checkout/components/shipping-infomation'
import ShippingMethod from '@/app/checkout/components/shipping-method'
import Button from '@/components/button/button'
import useOrderStore from '@/hooks/use-order-store'
import orderCheckoutService from '@/services/order-checkout'
import { useEffect } from 'react'

// ----------------------------------------------------------------------

type CartType = {
  product_id: number
  quantity: number
}

const CheckoutView = () => {
  const {
    getProductListInCart,
    cart,
    shipping,
    point,
    payment,
    summary,
    totalPayment
  } = useOrderStore()

  const submitPaymentOrder = async () => {
    const cartList: CartType[] = []

    cart.map((item) => {
      cartList.push({
        product_id: item.product_id,
        quantity: item.quantity
      })
    })

    const order = {
      cart: cartList,
      burn_point: point.burnPoint,
      sub_total_price: summary.total_price_thb,
      discount_price: 0,
      total_price: totalPayment,
      shipping_method_id: shipping.shippingMethod,
      shipping_address: shipping.shippingInformation.address,
      shipping_sub_district: shipping.shippingInformation.subDistrictName,
      shipping_district: shipping.shippingInformation.districtName,
      shipping_province: shipping.shippingInformation.provinceName,
      shipping_zip_code: shipping.shippingInformation.zipCode.toString(),
      recipient_first_name: shipping.shippingInformation.firstName,
      recipient_last_name: shipping.shippingInformation.lastName,
      recipient_phone_number: shipping.shippingInformation.mobileNumber,
      payment_method_id: Number(payment.paymentMethod),
      payment_information: {
        card_name: payment.paymentCreditInformation.name,
        card_number: payment.paymentCreditInformation.number,
        expire_date: payment.paymentCreditInformation.expiry,
        cvv: payment.paymentCreditInformation.cvv
      }
    }

    const result = await orderCheckoutService(order)

    if (result.data) {
      window.location.href = `/payment?id=${result.data.order_id}&total=${totalPayment}&card=${payment.paymentCreditInformation.number.slice(-4)}`
    } else {
      alert('Error Checkout, Please Try Again')
    }
  }

  useEffect(() => {
    getProductListInCart()
  }, [getProductListInCart])

  return (
    <div className="bg-white">
      <div className="mx-auto max-w-2xl px-4 sm:px-6 lg:max-w-7xl lg:px-8">
        <div className="w-full bg-white border-gray-200 px-5 py-10 text-gray-800">
          <div className="w-full">
            <div className="-mx-3 md:flex items-start">
              <div className="px-3 md:w-7/12 lg:pr-10">
                <OrderList />
                <ShippingInfomation />
                <ShippingMethod />
                <PaymentMethod />
              </div>
              <div className="px-3 md:w-5/12">
                {/* <Discount /> */}
                <OrderSummary />
                <Button
                  id="payment-now-btn"
                  type="button"
                  onClick={submitPaymentOrder}
                  isblock="true"
                >
                  PAY NOW
                </Button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}

export default CheckoutView
