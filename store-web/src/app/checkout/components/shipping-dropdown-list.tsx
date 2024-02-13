'use client'

import {
  DistrictType,
  ProvinceType,
  SubDistrictType
} from '@/app/checkout/components/shipping-infomation'

// ----------------------------------------------------------------------

type ShippingDropdownListProps = {
  id?: string
  label: string
  list: ProvinceType[] | DistrictType[] | SubDistrictType[]
  name: string
  setSelected: Function
}

const ShippingDropdownList = ({
  id,
  label,
  list,
  name,
  setSelected
}: ShippingDropdownListProps) => {
  const handleSelectChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    setSelected({
      id: Number(e.target.value),
      field: name
    })
  }

  return (
    <div className="mb-2">
      <label
        id={`${id}-label`}
        htmlFor={`${id}-select`}
        className="block mb-2 text-sm font-medium text-gray-900"
      >
        {label}
      </label>
      <select
        id={`${id}-select`}
        onChange={handleSelectChange}
        className="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-indigo-500 focus:border-indigo-500 block w-full p-2.5"
      >
        <option label="Select ..." value="0" />
        {list &&
          list.map((item: any) => (
            <option key={item.id} value={item.id}>
              {item.name_th}
            </option>
          ))}
      </select>
    </div>
  )
}

export default ShippingDropdownList
