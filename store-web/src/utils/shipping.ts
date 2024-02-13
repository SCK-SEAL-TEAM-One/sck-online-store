import SHIPPING_METHODS from '@/assets/data/shipping_method.json'

// ----------------------------------------------------------------------------

export const getShippingMethodById = (value: number) => {
  return SHIPPING_METHODS.find((shipp) => shipp.id === value)
}
