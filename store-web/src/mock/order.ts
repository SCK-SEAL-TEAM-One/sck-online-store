export const mockOrderCheckoutResponse = {
  status: 200,
  body: {
    order_id: 1,
  }
}

export const mockOrderUpdateStatusResponse = {
  status: 200,
  body: {
    order_id: 1,
    payment_date: '2023-01-31 10:00:00',
    shipping_method_id: 1,
    tracking_id: '51547878755545848512'
  }
}

export const mockOrderCompletedResponse = {
  status: 200,
  body: {
    date: '2024-01-30 10:00:00',
    status: 'success',
    orderId: 1,
    trackingNumber: '1234567890123TH'
  }
}
