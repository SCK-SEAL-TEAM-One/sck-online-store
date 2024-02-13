'use client'

import ProductItem from '@/app/cart/components/product-item'
import { ProductDetailInCart } from '@/services/cart/get-product-list'

// ----------------------------------------------------------------------

type ProductListProps = {
  list: ProductDetailInCart[]
}

const ProductList = ({ list }: ProductListProps) => {
  return (
    <div className="mt-8">
      <div className="flow-root">
        <ul role="list" className="-my-6 divide-y divide-gray-200">
          {list.map((product) => (
            <ProductItem
              key={`product-item-${product.id}`}
              isHiddenLable
              {...product}
            />
          ))}
        </ul>
      </div>
    </div>
  )
}

export default ProductList
