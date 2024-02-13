'use client'

import ProductCard from '@/app/product/list/components/product-card'
import { GetProductListServiceResponse } from '@/services/product-lists'

// const products = [
//   {
//     id: 1,
//     name: '43 Piece dinner set',
//     meta: '43-piece-dinner-set',
//     imageSrc:
//       'https://tailwindui.com/img/ecommerce-images/product-page-01-related-product-01.jpg',
//     imageAlt: '43 Piece dinner set.',
//     price: 35.0,
//     stock: 10
//   },
//   {
//     id: 2,
//     name: 'Steel Office Table',
//     meta: 'steel-office-table',
//     imageSrc:
//       'https://tailwindui.com/img/ecommerce-images/product-page-01-related-product-01.jpg',
//     imageAlt: 'Steel Office Table.',
//     price: 52,
//     stock: 8
//   },
//   {
//     id: 3,
//     name: 'Apple AirPods Pro',
//     meta: 'apple-airpods-pro',
//     imageSrc:
//       'https://tailwindui.com/img/ecommerce-images/product-page-01-related-product-01.jpg',
//     imageAlt: 'Apple AirPods Pro.',
//     price: 252.65,
//     stock: 10
//   }
// ]

type ProductListType = {
  products: GetProductListServiceResponse
}

// ----------------------------------------------------------------------

const ProductList = ({ products }: ProductListType) => {
  return (
    <div className="mt-6 grid grid-cols-1 gap-x-6 gap-y-10 sm:grid-cols-2 lg:grid-cols-4 xl:gap-x-8">
      {products?.products.map((product) => (
        <ProductCard key={`product-${product.id}`} data={product} />
      ))}
    </div>
  )
}

export default ProductList
