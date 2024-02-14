import config from '@/config'

export const receiptPointWithQuantity = (
  totalPrice: number,
  quantity: number
) => {
  const totalPayment = totalPrice * quantity
  // ปัดเศษลง
  return Math.floor(totalPayment / config.pointRate)
}

export const receiptPoint = (totalPayment: number) => {
  // ปัดเศษลง
  return Math.floor(totalPayment / config.pointRate)
}
