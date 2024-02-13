export const mockProductListResponse = {
  status: 200,
  body: {
    total: 2,
    products: [
      {
        id: 1,
        product_name: 'Balance Training Bicycle',
        product_price: 119.95,
        product_image: '/Balance_Training_Bicycle.png'
      },
      {
        id: 2,
        product_name: '43 Piece dinner Set',
        product_price: 12.95,
        product_image: '/43_Piece_dinner_Set.png'
      }
    ]
  }
}

export const mockProductDetailResponse = (id: string) => {
  if (id === '1') {
    return {
      status: 200,
      body: {
        id: 1,
        product_name: 'Balance Training Bicycle',
        product_price: 119.95,
        product_image: '/Balance_Training_Bicycle.png',
        stock: 5,
        product_brand: 'SportsFun'
      }
    }
  } else {
    return {
      status: 200,
      body: {
        id: 2,
        product_name: '43 Piece dinner Set',
        product_price: 12.95,
        product_image: '/43_Piece_dinner_Set.png',
        stock: 5,
        product_brand: 'SportsFun'
      }
    }
  }
}
