import { getShippingMethodById } from '@/utils/shipping'

describe('Utils > shipping > getShippingMethodById', () => {
  it('ต้องการข้อมูล kerry จาก id: 1', () => {
    const shippingMethodId = 1
    const actual = {
      id: 1,
      name: 'kerry',
      shippingTime: '4–10 business days',
      condition: '',
      price: 50
    }

    const pointUsed = getShippingMethodById(shippingMethodId)

    expect(pointUsed).to.deep.equal(actual)
  })

  it('ต้องการข้อมูล thai post จาก id: 2', () => {
    const shippingMethodId = 2
    const actual = {
      id: 2,
      name: 'thai post',
      shippingTime: '5–15 business days',
      condition: '',
      price: 50
    }

    const pointUsed = getShippingMethodById(shippingMethodId)

    expect(pointUsed).to.deep.equal(actual)
  })

  it('ต้องการข้อมูล lineman จาก id: 3', () => {
    const shippingMethodId = 3
    const actual = {
      id: 3,
      name: 'lineman',
      shippingTime: '1-2 business days',
      condition: '*Bangkok and perimeter only',
      price: 100
    }

    const pointUsed = getShippingMethodById(shippingMethodId)

    expect(pointUsed).to.deep.equal(actual)
  })
})
