'use client'

import ProductItem from '@/app/cart/components/product-item'
import Header3 from '@/components/typography/header3'
import useOrderStore from '@/hooks/use-order-store'

// ----------------------------------------------------------------------

const OrderList = () => {
  const { cart } = useOrderStore()

  return (
    <>
      <div className="w-full mx-auto text-gray-800 font-light mb-6 border-b border-gray-200 pb-6">
        <Header3>Orders</Header3>

        <ul
          id="order-product-list"
          role="list"
          className="-my-6 divide-y divide-gray-200"
        >
          {cart.map((product) => (
            <ProductItem
              key={`product-item-${product.id}`}
              isHiddenLable
              {...product}
            />
          ))}
        </ul>
      </div>
    </>
  )
}

export default OrderList
