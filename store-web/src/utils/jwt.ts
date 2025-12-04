export const decodeJWT = (token: string) => {
  const payloadB64 = token.split('.')[1]
  const payloadJson = JSON.parse(
    atob(payloadB64.replace(/-/g, '+').replace(/_/g, '/'))
  )
  return payloadJson
}
