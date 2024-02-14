'use client'

import Text from '@/components/typography/text'
import useOrderStore from '@/hooks/use-order-store'
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

      if (result.data) {
        setPoint(result.data.point)
      } else {
        setPoint(0)
      }
    }

    getPoint()
  }, [setPoint])

  return (
    <div className="flex justify-between mt-5">
      <div className="flex items-center mb-4">
        <input
          id="discount-use-point-input"
          type="checkbox"
          className="w-4 h-4 text-blue-600 bg-gray-100 border-gray-300 rounded focus:ring-indigo-500"
          onChange={handleUsePointChange}
        />
        <label
          id="discount-use-point-label"
          htmlFor="discount-use-point-input"
          className="ms-2 text-md font-medium text-gray-900 cursor-pointer"
        >
          Use your points
        </label>
      </div>
      <div>
        <Text id='discount-use-point-total' size="md" className="font-medium text-gray-900">
          {`${converNumber(point.point)} Points`}
        </Text>
      </div>
    </div>
  )
}

export default DiscountPoint
