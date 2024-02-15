'use client'

import Button from '@/components/button/button'
import InputQuantity from '@/components/input-quantity'
import Header1 from '@/components/typography/header1'
import Text from '@/components/typography/text'
import useOrderStore from '@/hooks/use-order-store'
import addToCartService from '@/services/cart/add-to-cart'
import { ProductDetailType } from '@/services/product-detail'
import { converNumber, convertCurrency, isNumber } from '@/utils/format'
import { receiptPoint } from '@/utils/point'
import { useState } from 'react'

// ----------------------------------------------------------------------

const ProductContent = (product: ProductDetailType) => {
  const [quantity, setQuantity] = useState(1)
  const { getProductListInCart } = useOrderStore()

  const handleQuantityChange = (e: { target: { value: string } }) => {
    if (isNumber(e.target.value)) {
      setQuantity(Number(e.target.value))
    }
  }

  const handleQuantityOnBlur = (e: { target: { value: string } }) => {
    const value = Number(e.target.value)
    if (value > 0 && value <= product.stock) {
      setQuantity(value)
    } else {
      setQuantity(1)
    }
  }

  const incrementQuantity = () => {
    if (quantity < product.stock) {
      setQuantity(quantity + 1)
    }
  }

  const decrementQuantity = () => {
    if (quantity > 1) {
      setQuantity(quantity - 1)
    }
  }

  const addToCart = async () => {
    if (product.stock > 0) {
      const result = await addToCartService({
        productId: product.id,
        quantity
      })

      // Add to cart is Success
      if (result.data) {
        getProductListInCart()
      }
    } else {
      alert('Cannot add to cart, this product out of stock.')
    }
  }

  return (
    <div className="mt-4 lg:row-span-3 lg:mt-0">
      <Header1
        id="product-detail-product-name"
        className="mb-4 tracking-tight sm:text-3xl"
      >
        {product.product_name}
      </Header1>

      <Text id="product-detail-brand" size="sm" className="mb-4">
        {product.product_brand}
      </Text>

      <Text
        id="product-detail-price-thb"
        size="xl"
        className="font-medium tracking-tight text-gray-900 mb-2"
      >
        {convertCurrency(product.product_price_thb, 'THB')}
      </Text>
      <Text
        id="product-detail-point"
        size="md"
        className="text-sm font-medium tracking-tight text-gray-400"
      >
        {`${converNumber(receiptPoint(product.product_price_thb))} Points`}
      </Text>

      <form className="mt-6">
        <InputQuantity
          label="Quantity:"
          id="product-detail-quantity"
          placeholder="999"
          increment={incrementQuantity}
          decrement={decrementQuantity}
          onChange={handleQuantityChange}
          onBlur={handleQuantityOnBlur}
          value={quantity}
          required
        />

        <Text id="product-detail-stock" size="sm" className="mt-3">
          {`Stock ${converNumber(product.stock)} items`}
        </Text>

        <Button
          id="product-detail-add-to-cart-btn"
          className="mt-6 disabled:opacity-50 disabled:cursor-not-allowed"
          type="button"
          isblock="true"
          onClick={addToCart}
          disabled={product.stock === 0}
        >
          Add to cart
        </Button>
      </form>
    </div>
  )
}

export default ProductContent
