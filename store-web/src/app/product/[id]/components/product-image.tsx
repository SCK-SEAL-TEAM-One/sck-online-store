'use client'

import Image from '@/components/image'
import config from '@/config'
import { GetProductDetailServiceResponse } from '@/services/product-detail'

// ----------------------------------------------------------------------

const ProductImage = (product: GetProductDetailServiceResponse) => {
  return (
    <div className="lg:col-span-2 lg:border-r lg:border-gray-200 lg:pr-8">
      <div className="aspect-h-2 aspect-w-3 hidden overflow-hidden rounded-lg lg:block">
        <Image
          src={`${config.imageUrl}/${product.product_image}`}
          alt={product.product_name}
          width={767}
          height={575}
          className="h-full w-full object-contain object-center"
        />
      </div>
    </div>
  )
}

export default ProductImage
