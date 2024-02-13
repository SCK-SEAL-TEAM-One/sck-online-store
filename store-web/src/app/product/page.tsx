import { redirect } from 'next/navigation'

export const metadata = {
  title: 'Product Lists'
}

const ProductPage = () => redirect('/product/list')

export default ProductPage
