'use client'

import RemoveItem from '@/app/cart/components/remove-item'
import Image from '@/components/image'
import InputQuantity from '@/components/input-quantity'
import Text from '@/components/typography/text'
import config from '@/config'
import useOrderStore from '@/hooks/use-order-store'
import { ProductDetailInCart } from '@/services/cart/get-product-list'
import updateProductInCartService from '@/services/cart/update-product'
import { converNumber, convertCurrency, isNumber } from '@/utils/format'
import { useEffect, useState } from 'react'

// ----------------------------------------------------------------------

type ProductItemProps = ProductDetailInCart & {
  isHiddenLable?: boolean
}

const ProductItem = ({
  id,
  product_id,
  quantity,
  product_name,
  product_price_thb,
  product_image,
  stock,
  isHiddenLable = false
}: ProductItemProps) => {
  const [newQuantity, setNewQuantity] = useState(1)
  const { getProductListInCart } = useOrderStore()

  const handleQuantityChange = (e: { target: { value: string } }) => {
    if (isNumber(e.target.value)) {
      setNewQuantity(Number(e.target.value))
    }
  }

  const handleQuantityOnBlur = (e: { target: { value: string } }) => {
    const value = Number(e.target.value)
    if (value > 0 && value <= stock) {
      updateQuantity(value)
    } else {
      updateQuantity(1)
    }
  }

  const incrementQuantity = () => {
    if (newQuantity < stock) {
      updateQuantity(newQuantity + 1)
    }
  }

  const decrementQuantity = () => {
    if (newQuantity > 1) {
      updateQuantity(newQuantity - 1)
    }
  }

  const updateQuantity = async (qt: number) => {
    setNewQuantity(qt)

    const result = await updateProductInCartService({
      productId: id,
      quantity: qt
    })

    if (result.data) {
      // Update Cart
      getProductListInCart()
    } else if (result.message) {
      alert('Cannot update quantity, ' + result.message)
    }
  }

  const handleRemoveItem = async () => {
    if (window.confirm('Do you want to remove this item?' + id)) {
      const result = await updateProductInCartService({
        productId: id,
        quantity: 0
      })

      if (result.data) {
        // Get Cart Service
        alert('Remove item success')

        // Get Services for update Product List in cart
        getProductListInCart()
      } else if (result.message) {
        alert('Cannot remove item, ' + result.message)
      }
    }
  }

  useEffect(() => {
    setNewQuantity(quantity)
  }, [quantity])

  return (
    <li className="flex py-6">
      <div className="h-32 w-32 flex-shrink-0 overflow-hidden rounded-md border border-gray-200">
        <Image
          id={`product-${product_id}-image`}
          src={`${config.imageUrl}/${product_image}`}
          alt={product_name}
          width={94}
          height={94}
          className="h-full w-full object-contain object-center bg-white"
        />
      </div>

      <div className="ml-4 flex flex-1 flex-col">
        <div>
          <div className="flex justify-between text-base font-medium text-gray-900">
            <h3 id={`product-${product_id}-name`}>
              <a href={`/product/${product_id}`}>{product_name}</a>
            </h3>
            <p id={`product-${product_id}-price`} className="ml-4">
              {convertCurrency(product_price_thb, 'THB')}
            </p>
          </div>
          <Text
            id={`product-${product_id}-stock`}
            className="mt-1 text-sm text-gray-500"
          >
            {`Stock ${converNumber(stock)} items`}
          </Text>
        </div>

        <div className="flex flex-1 items-end justify-between text-sm mt-4">
          <InputQuantity
            id={`product-${product_id}-quantity`}
            placeholder="999"
            increment={incrementQuantity}
            decrement={decrementQuantity}
            onChange={handleQuantityChange}
            onBlur={handleQuantityOnBlur}
            value={newQuantity}
            isHiddenLable={isHiddenLable}
            required
          />

          <RemoveItem
            id={`product-${product_id}-remove-btn`}
            onClick={handleRemoveItem}
          />
        </div>
      </div>
    </li>
  )
}

export default ProductItem
