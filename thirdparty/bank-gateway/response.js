function () {
  const txn = (Math.random() + 1).toString(36).substring(2);
  return { 
    statusCode: 200,
    headers: {
      "Content-Type": "application/json; charset=utf-8"
    },
    body: {
      status: "completed",
      payment_date: new Date,
      transaction_id: txn
    }
  }
}