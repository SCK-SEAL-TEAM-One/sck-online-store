'use client'

import Header2 from '@/components/typography/header2'

// ----------------------------------------------------------------------

type ProductTitleProps = {
  title: string
}

const ProductTitle = ({ title }: ProductTitleProps) => {
  return <Header2 className="text-gray-600">{title}</Header2>
}

export default ProductTitle
