import { convertCurrency, converNumber, isNumber } from '@/utils/format'

describe('Utils > format > convertCurrency', () => {
  it('ต้องการเห็น ข้อมูลหลัง format USD $1,234.50 จากราคาที่จ่าย 1234.50 บาท', () => {
    const price = 1234.50
    const currency = 'USD'
    const actual = '$1,234.50'

    const point = convertCurrency(price, currency)

    expect(point).to.equal(actual)
  })

  it('ต้องการเห็น ข้อมูลหลัง format THB $1,234.50 จากราคาที่จ่าย 1234.50 บาท', () => {
    const price = 1234.50
    const currency = 'THB'
    const actual = '฿1,234.50'

    const point = convertCurrency(price, currency)

    expect(point).to.equal(actual)
  })

  it('ต้องการเห็น ข้อมูลหลัง format THB $1,234.00 จากราคาที่จ่าย 1234 บาท', () => {
    const price = 1000
    const currency = 'THB'
    const actual = '฿1,000.00'

    const point = convertCurrency(price, currency)

    expect(point).to.equal(actual)
  })
})

describe('Utils > format > converNumber', () => {
  it('ต้องการเห็น ข้อมูลหลัง format 1,234.5 จากราคาที่จ่าย 1234.50 บาท', () => {
    const price = 1234.50
    const actual = '1,234.5'

    const point = converNumber(price)

    expect(point).to.equal(actual)
  })

  it('ต้องการเห็น ข้อมูลหลัง format 1,234 จากราคาที่จ่าย 1234 บาท', () => {
    const price = 1234
    const actual = '1,234'

    const point = converNumber(price)

    expect(point).to.equal(actual)
  })

  it('ต้องการเห็น ข้อมูลหลัง format 1,234.65 จากราคาที่จ่าย 1234.65 บาท', () => {
    const price = 1234.65
    const actual = '1,234.65'

    const point = converNumber(price)

    expect(point).to.equal(actual)
  })
})

describe('Utils > format > isNumber', () => {
  it('ต้องการเห็น ข้อมูลหลัง true จากราคาแบบ string "234.50" บาท', () => {
    const price = '234.50'
    const actual = true

    const point = isNumber(price)

    expect(point).to.equal(actual)
  })

  it('ต้องการเห็น ข้อมูลหลัง false จากราคาแบบ string "สองร้อย" บาท', () => {
    const price = 'สองร้อย'
    const actual = false

    const point = isNumber(price)

    expect(point).to.equal(actual)
  })
})