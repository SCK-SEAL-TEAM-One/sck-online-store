import config from '@/config'

export const receiptPoint = (totalPayment: number) => {
  // ปัดเศษลง
  return Math.floor(totalPayment / config.pointRate)
}
