'use client'

import RemoveItem from '@/app/cart/components/remove-item'
import useOrderStore from '@/hooks/use-order-store'
import Image from '@/components/image'
import InputQuantity from '@/components/input-quantity'
import Text from '@/components/typography/text'
import config from '@/config'
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
  // user_id,
  product_id,
  quantity,
  product_name,
  product_price,
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
    if (value > 0 && value < stock) {
      setNewQuantity(value)
    } else {
      setNewQuantity(1)
    }
  }

  const incrementQuantity = () => {
    if (newQuantity < stock) {
      setNewQuantity(newQuantity + 1)
    }
  }

  const decrementQuantity = () => {
    if (newQuantity > 1) {
      setNewQuantity(newQuantity - 1)
    }
  }

  const handleRemoveItem = async () => {
    const result = await updateProductInCartService({
      product_id: id,
      quantity: 0
    })

    if (result) {
      // Get Cart Service
      alert('Remove item success')

      // Get Services for update Product List in cart
      getProductListInCart()
    }
  }

  useEffect(() => {
    setNewQuantity(quantity)
  }, [quantity])

  return (
    <li className="flex py-6">
      <div className="h-32 w-32 flex-shrink-0 overflow-hidden rounded-md border border-gray-200">
        <Image
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
            <h3>
              <a href={`/product/${product_id}`}>{product_name}</a>
            </h3>
            <p className="ml-4">{convertCurrency(product_price)}</p>
          </div>
          <Text className="mt-1 text-sm text-gray-500">
            {`Stock ${converNumber(stock)} items`}
          </Text>
        </div>

        <div className="flex flex-1 items-end justify-between text-sm mt-4">
          <InputQuantity
            id="quantity"
            placeholder="999"
            increment={incrementQuantity}
            decrement={decrementQuantity}
            onChange={handleQuantityChange}
            onBlur={handleQuantityOnBlur}
            value={newQuantity}
            isHiddenLable={isHiddenLable}
            required
          />

          <RemoveItem onClick={handleRemoveItem} />
        </div>
      </div>
    </li>
  )
}

export default ProductItem
