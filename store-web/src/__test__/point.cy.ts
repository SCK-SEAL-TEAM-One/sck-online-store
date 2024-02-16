import { receiptPoint } from '@/utils/point'

describe('Utils > point > receiptPoint', () => {
  it('ต้องการเห็น จำนวนแต้มที่ได้ 10 แต้ม จากราคาที่จ่าย 1000 บาท', () => {
    const payment = 1000
    const actual = 10

    const point = receiptPoint(payment)

    expect(point).to.equal(actual)
  })

  it('ต้องการเห็น จำนวนแต้มที่ได้ 10 แต้ม จากราคาที่จ่าย 1060 บาท', () => {
    const payment = 1060
    const actual = 10

    const point = receiptPoint(payment)

    expect(point).to.equal(actual)
  })

  it('ต้องการเห็น จำนวนแต้มที่ได้ 0 แต้ม จากราคาที่จ่าย 60 บาท', () => {
    const payment = 60
    const actual = 0

    const point = receiptPoint(payment)

    expect(point).to.equal(actual)
  })
})
