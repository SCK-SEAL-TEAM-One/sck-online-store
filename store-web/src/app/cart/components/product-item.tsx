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
import { receiptPoint } from '@/utils/point'
import { useEffect, useState } from 'react'

// ----------------------------------------------------------------------

type ProductItemProps = ProductDetailInCart & {
  isHiddenLable?: boolean
}

const ProductItem = ({
  product_id,
  quantity,
  product_name,
  product_price_thb,
  product_image,
  stock,
  isHiddenLable = false
}: ProductItemProps) => {
  const [newQuantity, setNewQuantity] = useState(quantity)
  const { getProductListInCart } = useOrderStore()

  const handleQuantityChange = (e: { target: { value: string } }) => {
    if (isNumber(e.target.value)) {
      setNewQuantity(Number(e.target.value))
    }
  }

  const handleQuantityOnBlur = (e: { target: { value: string } }) => {
    const value = Number(e.target.value)

    if (value !== quantity) {
      if (value > 0 && value <= stock) {
        setNewQuantity(value)
        updateQuantity(value)
      } else {
        setNewQuantity(quantity)
      }
    }
  }

  const incrementQuantity = () => {
    const incrementQty = quantity + 1
    if (incrementQty <= stock) {
      setNewQuantity(incrementQty)
      updateQuantity(incrementQty)
    }
  }

  const decrementQuantity = () => {
    const decrementQty = quantity - 1
    if (decrementQty >= 1) {
      setNewQuantity(decrementQty)
      updateQuantity(decrementQty)
    }
  }

  const updateQuantity = async (qty: number) => {
    setNewQuantity(qty)

    const result = await updateProductInCartService({
      productId: product_id,
      quantity: qty
    })

    if (result.data) {
      // Update Cart
      getProductListInCart()
    } else if (result.message) {
      setNewQuantity(quantity)
      alert('Cannot update quantity, ' + result.message)
    }
  }

  const handleRemoveItem = async () => {
    if (window.confirm('Do you want to remove "' + product_name + '"?,')) {
      const result = await updateProductInCartService({
        productId: product_id,
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
            <div className="flex flex-col items-end">
              <p id={`product-${product_id}-price`} className="ml-4">
                {convertCurrency(product_price_thb * quantity, 'THB')}
              </p>
              <p
                id={`product-${product_id}-point`}
                className="ml-4 text-sm text-gray-600"
              >
                {`${converNumber(
                  receiptPoint(product_price_thb * quantity)
                )} Points`}
              </p>
            </div>
          </div>
          <Text
            id={`product-${product_id}-stock`}
            className="mt-1 text-sm text-gray-500"
          >
            {`Stock ${converNumber(stock)} items`}
          </Text>
        </div>

        <div className="flex flex-1 items-end justify-between">
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
