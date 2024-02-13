'use client'

import Image from '@/components/image'
import Header1 from '@/components/typography/header1'
import Text from '@/components/typography/text'

// ----------------------------------------------------------------------

const ProductNotFound = () => {
  return (
    <div className="flex flex-col justify-center items-center text-gray-500 my-10">
      <Image
        src="/search.png"
        width={200}
        height={200}
        alt="search product"
        className="opacity-60 my-5"
      />
      <Header1 className="text-gray-600">ไม่พบผลการค้นหา</Header1>
      <Text className="text-gray-400">
        ลองใช้คำอื่นที่แตกต่างหรือคำอื่นที่มีความหมายกว้างกว่านี้
      </Text>
    </div>
  )
}

export default ProductNotFound
