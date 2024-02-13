export const convertCurrency = (value: number, currency?: string) => {
  if (currency?.toLocaleLowerCase() === 'thb') {
    return new Intl.NumberFormat('th-TH', {
      style: 'currency',
      currency: 'THB'
    }).format(value)
  }

  return new Intl.NumberFormat('en-US', {
    style: 'currency',
    currency: 'USD'
  }).format(value)
}

export const converNumber = (value: number) => {
  return value.toLocaleString('en-US')
}

export const isNumber = (value: string) => {
  return /^\d*\.?\d*$/.test(value)
}
