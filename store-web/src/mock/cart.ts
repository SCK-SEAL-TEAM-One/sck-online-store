export const mockCartListResponse = {
  status: 200,
  body: [
    {
      id: 1,
      user_id: 1,
      product_id: 1,
      quantity: 1,
      product_name: 'Balance Training Bicycle',
      product_price: 119.95,
      product_image: '/Balance_Training_Bicycle.png',
      stock: 5,
      product_brand: 'SportsFun'
    },
    {
      id: 4,
      user_id: 1,
      product_id: 3,
      quantity: 2,
      product_name: 'Horses and Unicorns Set',
      product_price: 24.95,
      product_image: '/Horses_and_Unicorns_Set.png',
      stock: 3,
      product_brand: 'CoolKidZ'
    }
  ]
}

export const mockAddToCartResponse = {
  status: 200,
  body: {
    status: 'updated' // added, updated
  }
}

export const mockUpdateCartResponse = {
  status: 200,
  body: {
    status: 'updated' // deleted, updated
  }
}
