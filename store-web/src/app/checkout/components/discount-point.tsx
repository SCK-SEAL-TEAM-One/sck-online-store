'use client'

import useOrderStore from '@/hooks/use-order-store'
import Text from '@/components/typography/text'
import getPointService from '@/services/point'
import { converNumber } from '@/utils/format'
import { useEffect } from 'react'

// ----------------------------------------------------------------------

const DiscountPoint = () => {
  const { point, setPoint, setIsUsePoint } = useOrderStore((state) => state)

  const handleUsePointChange = (e: { target: { checked: boolean } }) => {
    setIsUsePoint(e.target.checked)
  }

  useEffect(() => {
    const getPoint = async () => {
      const result = await getPointService()
      setPoint(result.point)
    }

    getPoint()
  }, [setPoint])

  return (
    <div className="flex justify-between mt-5">
      <div className="flex items-center mb-4">
        <input
          id="default-checkbox"
          type="checkbox"
          className="w-4 h-4 text-blue-600 bg-gray-100 border-gray-300 rounded focus:ring-indigo-500"
          onChange={handleUsePointChange}
        />
        <label
          htmlFor="default-checkbox"
          className="ms-2 text-md font-medium text-gray-900 cursor-pointer"
        >
          Use your points
        </label>
      </div>
      <div>
        <Text size="md" className="font-medium text-gray-900">
          {`${converNumber(point.point)} Points`}
        </Text>
      </div>
    </div>
  )
}

export default DiscountPoint
