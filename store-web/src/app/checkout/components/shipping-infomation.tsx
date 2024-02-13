'use client'

import ShippingDropdownList from '@/app/checkout/components/shipping-dropdown-list'
import InputField from '@/components/input-field'
import Header3 from '@/components/typography/header3'
import React, { useState } from 'react'

import useOrderStore from '@/hooks/use-order-store'
import DISTRICT_LIST from '@/assets/data/api_district.json'
import PROVINCE_LIST from '@/assets/data/api_province.json'
import SUB_DISTRICT_LIST from '@/assets/data/api_sub_district.json'

// ----------------------------------------------------------------------

export type ProvinceType = {
  id: number
  name_th: string
  name_en: string
  geography_id: number
  created_at: string
  updated_at: string
  deleted_at: string | null
}

export type DistrictType = {
  id: number
  name_th: string
  name_en: string
  province_id: number
  created_at: string
  updated_at: string
  deleted_at: string | null
}

export type SubDistrictType = {
  id: number
  zip_code: number
  name_th: string
  name_en: string
  amphure_id: number
  created_at: string
  updated_at: string
  deleted_at: string | null
}

const ShippingInfomation = () => {
  const { setShippingInformation } = useOrderStore((state) => state)

  const [provinceList] = useState<ProvinceType[]>(PROVINCE_LIST)
  const [districtList, setDistrictList] = useState<DistrictType[]>([])
  const [subDistrictList, setSubDistrictList] = useState<SubDistrictType[]>([])
  const [addressInfo, setAddressInfo] = useState({
    firstName: '',
    lastName: '',
    address: '',
    mobileNumber: '',
    provinceId: 0,
    provinceName: '',
    districtId: 0,
    districtName: '',
    subDistrictId: 0,
    subDistrictName: '',
    zipCode: 0,
    focused: ''
  })

  const handleInputFocus = ({ target }: React.FocusEvent<HTMLInputElement>) => {
    setAddressInfo({
      ...addressInfo,
      focused: target.name
    })

    // Save Shipping Information on Checkout Store
    setShippingInformation(addressInfo)
  }

  const handleAddressInputChange = ({
    target
  }: React.ChangeEvent<HTMLInputElement>) => {
    if (target.name === 'firstName') {
      setAddressInfo({ ...addressInfo, firstName: target.value })
    } else if (target.name === 'lastName') {
      setAddressInfo({
        ...addressInfo,
        lastName: target.value
      })
    } else if (target.name === 'address') {
      setAddressInfo({
        ...addressInfo,
        address: target.value
      })
    } else if (target.name === 'mobileNumber') {
      setAddressInfo({
        ...addressInfo,
        mobileNumber: target.value
      })
    }

    // Save Shipping Information on Checkout Store
    setShippingInformation(addressInfo)
  }

  const handleAddressSelectChange = (selected: {
    id: number
    field: string
  }) => {
    if (selected.field === 'province') {
      const province = PROVINCE_LIST.filter((p: any) =>
        p.id === selected.id ? p : null
      )

      const newProvinceInformation = {
        ...addressInfo,
        provinceId: selected.id,
        provinceName: province[0].name_th
      }

      setAddressInfo(newProvinceInformation)
      setShippingInformation(newProvinceInformation)

      getDistrictList(selected.id)
    } else if (selected.field === 'district') {
      const district = DISTRICT_LIST.filter((d: any) =>
        d.id === selected.id ? d : null
      )

      const newDistrictInformation = {
        ...addressInfo,
        districtId: selected.id,
        districtName: district[0].name_th,
        zipCode: 0
      }

      setAddressInfo(newDistrictInformation)
      setShippingInformation(newDistrictInformation)

      getSubDistrictList(selected.id)
    } else if (selected.field === 'subDistrict') {
      const subDistrict = SUB_DISTRICT_LIST.filter((d: any) =>
        d.id === selected.id ? d : null
      )

      const newSubDistrictInformation = {
        ...addressInfo,
        subDistrictId: selected.id,
        subDistrictName: subDistrict[0].name_th,
        zipCode: subDistrict[0].zip_code
      }

      setAddressInfo(newSubDistrictInformation)
      setShippingInformation(newSubDistrictInformation)
    }
  }

  const getDistrictList = (provinceId: number) => {
    const district = DISTRICT_LIST.filter((d: DistrictType) =>
      d.province_id === provinceId ? d : null
    )
    setDistrictList(district)
  }

  const getSubDistrictList = (districtId: number) => {
    const subDistrict = SUB_DISTRICT_LIST.filter((d: any) =>
      d.amphure_id === districtId ? d : null
    )
    setSubDistrictList(subDistrict)
  }

  return (
    <div className="mb-6 border-b border-gray-200 pb-6">
      <Header3>Shipping information</Header3>

      <div className="grid gap-6 mb-2 md:grid-cols-2">
        <InputField
          id="shipping-form-first-name"
          label="First name"
          type="text"
          name="firstName"
          placeholder="first name"
          required
          onChange={handleAddressInputChange}
          onFocus={handleInputFocus}
        />

        <InputField
          id="shipping-form-last-name"
          label="Last name"
          type="text"
          name="lastName"
          placeholder="last name"
          required
          onChange={handleAddressInputChange}
          onFocus={handleInputFocus}
        />
      </div>

      <InputField
        id="shipping-form-address"
        label="Address (Building, Street, etc.)"
        type="text"
        name="address"
        placeholder="address"
        required
        maxLength={150}
        onChange={handleAddressInputChange}
        onFocus={handleInputFocus}
      />

      <ShippingDropdownList
        id="shipping-form-province"
        label="Province: "
        list={provinceList}
        name="province"
        setSelected={handleAddressSelectChange}
      />

      <ShippingDropdownList
        id="shipping-form-district"
        label="District: "
        list={districtList}
        name="district"
        setSelected={handleAddressSelectChange}
      />

      <ShippingDropdownList
        id="shipping-form-sub-district"
        label="Sub-district: "
        list={subDistrictList}
        name="subDistrict"
        setSelected={handleAddressSelectChange}
      />

      <InputField
        id="shipping-form-zipcode"
        label="Zipcode"
        type="text"
        name="zipCode"
        placeholder="zipcode"
        maxLength={5}
        value={addressInfo.zipCode ? addressInfo.zipCode.toString() : ''}
        readOnly
        disabled
      />

      <InputField
        id="shipping-form-mobile"
        label="Mobile number (For Contact)"
        type="tel"
        name="mobileNumber"
        placeholder="0923456789"
        maxLength={10}
        onChange={handleAddressInputChange}
        onFocus={handleInputFocus}
        required
      />
    </div>
  )
}

export default ShippingInfomation
