'use client'

import ProductContent from '@/app/product/[id]/components/product-content'
import ProductImage from '@/app/product/[id]/components/product-image'
import getProductDetailService, {
  GetProductDetailServiceResponse
} from '@/services/product-detail'
import { useParams } from 'next/navigation'
import { useEffect, useState } from 'react'

// ----------------------------------------------------------------------------

const ProductDetailView = () => {
  const { id } = useParams<{ id: string }>()
  const [productDetail, setProductDetail] =
    useState<GetProductDetailServiceResponse | null>(null)

  useEffect(() => {
    const getProductDetail = async () => {
      const result = await getProductDetailService(id)
      setProductDetail(result)
    }

    getProductDetail()
  }, [id])

  return (
    <div className="bg-white min-h-[calc(100vh-88px)]">
      {/* Product info */}
      {productDetail && productDetail.data ? (
        <div className="mx-auto max-w-2xl px-4 pb-16 pt-10 sm:px-6 lg:grid lg:max-w-7xl lg:grid-cols-3 lg:grid-rows-[auto,auto,1fr] lg:gap-x-8 lg:px-8 lg:pb-24 lg:pt-16">
          {productDetail ? (
            <>
              {/* Images */}
              <ProductImage {...productDetail.data} />
              {/* ProductContent */}
              <ProductContent {...productDetail.data} />
            </>
          ) : null}
        </div>
      ) : null}
    </div>
  )
}

export default ProductDetailView
