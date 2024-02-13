import config from '@/config'

// ----------------------------------------------------------------------------

export const subTotal = (priceList: number[]) => {
  let total = 0

  for (let i = 0; i < priceList.length; i++) {
    total += priceList[i]
  }

  return total
}

export const pointBurn = (
  point: number,
  subTotal: number
) => {
  let pointsUsed = 0

  if (point <= subTotal) {
    pointsUsed = point
  } else {
    // ปัดเศษขึ้น
    pointsUsed = Math.ceil(subTotal)
  }

  return pointsUsed
}

export const totalPayment = (
  isUsePoint: boolean,
  pointsUsed: number,
  subTotal: number,
  shippingFee: number
) => {
  let totalPayment = 0

  if (isUsePoint) {
    if (subTotal <= pointsUsed) {
      totalPayment = shippingFee
    } else {
      totalPayment = subTotal - pointsUsed + shippingFee
    }
  } else {
    totalPayment = subTotal + shippingFee
  }

  return totalPayment
}

export const receiptPoint = (totalPayment: number) => {
  // ปัดเศษลง
  return Math.floor(totalPayment / config.pointRate)
}
