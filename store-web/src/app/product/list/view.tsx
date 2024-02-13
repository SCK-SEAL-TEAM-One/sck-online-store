'use client'

import ProductList from '@/app/product/list/components/product-list'
import ProductNotFound from '@/app/product/list/components/product-not-found'
import ProductTitle from '@/app/product/list/components/product-title'
import SearchForm from '@/app/product/list/components/search-form'
import getProductListService, {
  GetProductListServiceResponse
} from '@/services/product-lists'
import axios from 'axios'
import { useSearchParams } from 'next/navigation'
import { useEffect, useState } from 'react'

// ----------------------------------------------------------------------

const ProductView = () => {
  const searchParams = useSearchParams()
  const searchKeyword = searchParams.get('keyword') || ''

  const [products, setProducts] =
    useState<GetProductListServiceResponse | null>(null)
  const [keyword, setKeword] = useState(searchKeyword)

  const getProductList = async (keyword: string) => {
    const productList = await getProductListService({
      keyword,
      limit: 20,
      offset: 0
    })

    if (productList.data) {
      setProducts(productList)
    } else {
      setProducts(null)
    }
  }

  const onSubmitSearch = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault()
    // getProductList(keyword)
    window.location.href = '/product/list?keyword=' + keyword
  }

  const onChangeSearchKeyword = (e: React.ChangeEvent<HTMLInputElement>) => {
    setKeword(e.target.value)
  }

  useEffect(() => {
    getProductList(searchKeyword)
  }, [searchKeyword])

  return (
    <div className="bg-white">
      <div className="min-h-[calc(100vh-88px)] mx-auto max-w-2xl px-4 py-16 sm:px-6 sm:py-6 lg:max-w-7xl lg:px-8">
        <SearchForm
          keyword={keyword}
          onChangeSearchKeyword={onChangeSearchKeyword}
          onSubmitSearch={onSubmitSearch}
        />

        <ProductTitle id='product-title' title="All Products" />

        {products && products.data && products.data?.total > 0 ? (
          <ProductList products={products} />
        ) : (
          <ProductNotFound />
        )}
      </div>
    </div>
  )
}

export default ProductView
