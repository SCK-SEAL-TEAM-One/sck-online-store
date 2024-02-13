'use client'

import ProductCard from '@/app/product/list/components/product-card'
import { GetProductListServiceResponse } from '@/services/product-lists'

type ProductListType = {
  products: GetProductListServiceResponse
}

// ----------------------------------------------------------------------

const ProductList = ({ products }: ProductListType) => {
  return (
  <div id='product-list' className="mt-6 grid grid-cols-1 gap-x-6 gap-y-10 sm:grid-cols-2 lg:grid-cols-4 xl:gap-x-8">
      {products.data?.products.map((product) => (
        <ProductCard key={`product-${product.id}`} data={product} />
      ))}
    </div>
  )
}

export default ProductList
