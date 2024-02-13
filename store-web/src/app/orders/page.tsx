import { redirect } from 'next/navigation'

export const metadata = {
  title: 'Order'
}

const OrderPage = () => redirect('/order/completed')

export default OrderPage
