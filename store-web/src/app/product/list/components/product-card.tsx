'use client'

import Image from '@/components/image'
import Header4 from '@/components/typography/header4'
import Text from '@/components/typography/text'
import config from '@/config'
import { ProductDetailType } from '@/services/product-lists'
import { convertCurrency } from '@/utils/format'

// ----------------------------------------------------------------------

type ProductCardProps = {
  data: ProductDetailType
}

const ProductCard = ({ data }: ProductCardProps) => {
  return (
    <div id={`product-card-${data.id}`} className="group relative">
      <div className="aspect-h-1 aspect-w-1 w-full overflow-hidden rounded-md bg-gray-200 lg:aspect-none group-hover:opacity-75 lg:h-80">
        <Image
          id={`product-card-image-${data.id}`}
          src={`${config.imageUrl}${data.product_image}`}
          alt={data.product_name}
          width={280}
          height={320}
          className="h-full w-full object-contain bg-white object-center lg:h-full lg:w-full"
        />
      </div>
      <div className="mt-4 mb-1 flex align-top justify-between">
        <Header4 id={`product-card-name-${data.id}`} className="text-gray-700">
          <a href={`/product/${data.id}`}>
            <span aria-hidden="true" className="absolute inset-0" />
            {data.product_name}
          </a>
        </Header4>

        <Text
          id={`product-card-price-${data.id}`}
          size="md"
          className="font-medium text-gray-900"
        >
          {convertCurrency(data.product_price_thb, 'THB')}
        </Text>
      </div>
    </div>
  )
}

export default ProductCard
