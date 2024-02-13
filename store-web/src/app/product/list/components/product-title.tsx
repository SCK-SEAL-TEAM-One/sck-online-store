'use client'

import Header2 from '@/components/typography/header2'

// ----------------------------------------------------------------------

type ProductTitleProps = {
  title: string
  id?: string
}

const ProductTitle = ({ title, id }: ProductTitleProps) => {
  return (
    <Header2 id={id} className="text-gray-600">
      {title}
    </Header2>
  )
}

export default ProductTitle
