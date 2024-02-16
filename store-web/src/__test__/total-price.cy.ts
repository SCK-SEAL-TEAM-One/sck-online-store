import { pointBurn, subTotal, totalPayment } from '@/utils/total-price'

describe('Utils > total price > Sub Total', () => {
  it('ต้องการเห็น ผลรวมสินค้า 620 บาท จากรายการสินค้าทั้งหมด', () => {
    const priceList = [
      {
        price: 100,
        quantity: 1
      },
      {
        price: 200,
        quantity: 2
      },
      {
        price: 120,
        quantity: 1
      }
    ]
    const actual = 620

    const pointUsed = subTotal(priceList)

    expect(pointUsed).to.equal(actual)
  })
})

describe('Utils > total price > Total Price', () => {
  it('ต้องการเห็น เงินทั้งหมด 550 บาท จากราคารวมสินค้าทั้งหมด 500 บาท และค่าขนส่ง 50 บาท โดยไม่ใช้ส่วนลดจากแต้ม', () => {
    const isUsePoint = false
    const pointsUsed = 0
    const subTotal = 500
    const shippingFee = 50
    const actual = 550

    const total = totalPayment(isUsePoint, pointsUsed, subTotal, shippingFee)

    expect(total).to.equal(actual)
  })

  it('ต้องการเห็น เงินทั้งหมด 450 บาท จากราคารวมสินค้าทั้งหมด 500 บาท และค่าขนส่ง 50 บาท โดยใช้ส่วนลดจากแต้ม 100 แต้ม', () => {
    const isUsePoint = true
    const pointsUsed = 100
    const subTotal = 500
    const shippingFee = 50
    const actual = 450

    const total = totalPayment(isUsePoint, pointsUsed, subTotal, shippingFee)

    expect(total).to.equal(actual)
  })

  it('ต้องการเห็น เงินทั้งหมด 50 บาท จากราคารวมสินค้าทั้งหมด 100 บาท และค่าขนส่ง 50 บาท โดยใช้ส่วนลดจากแต้ม 100 แต้ม', () => {
    const isUsePoint = true
    const pointsUsed = 100
    const subTotal = 100
    const shippingFee = 50
    const actual = 50

    const total = totalPayment(isUsePoint, pointsUsed, subTotal, shippingFee)

    expect(total).to.equal(actual)
  })

  it('ต้องการเห็น เงินทั้งหมด 50 บาท จากราคารวมสินค้าทั้งหมด 100 บาท และค่าขนส่ง 50 บาท โดยใช้ส่วนลดจากแต้ม 150 แต้ม', () => {
    const isUsePoint = true
    const pointsUsed = 150
    const subTotal = 100
    const shippingFee = 50
    const actual = 50

    const total = totalPayment(isUsePoint, pointsUsed, subTotal, shippingFee)

    expect(total).to.equal(actual)
  })
})

describe('Utils > total price > Point Burn', () => {
  it('ต้องการเห็น แต้มที่ใช้ 100 แต้ม จากราคารวมสินค้าทั้งหมด 100 บาท โดยมีแต้มทั้งหมด 100 แต้ม', () => {
    const point = 100
    const subTotal = 100
    const actual = 100

    const pointUsed = pointBurn(point, subTotal)

    expect(pointUsed).to.equal(actual)
  })

  it('ต้องการเห็น แต้มที่ใช้ 100 แต้ม จากราคารวมสินค้าทั้งหมด 100 บาท โดยมีแต้มทั้งหมด 200 แต้ม', () => {
    const point = 100
    const subTotal = 100
    const actual = 100

    const pointUsed = pointBurn(point, subTotal)

    expect(pointUsed).to.equal(actual)
  })

  it('ต้องการเห็น แต้มที่ใช้ 50 แต้ม จากราคารวมสินค้าทั้งหมด 100 บาท โดยมีแต้มทั้งหมด 50 แต้ม', () => {
    const point = 50
    const subTotal = 100
    const actual = 50

    const pointUsed = pointBurn(point, subTotal)

    expect(pointUsed).to.equal(actual)
  })
})
