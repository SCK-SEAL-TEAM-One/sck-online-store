import { AxiosError, isAxiosError } from 'axios'

export const handleServiceError = (error: AxiosError | unknown) => {
  let response = null
  if (isAxiosError(error)) {
    response = handleAxiosError(error)
  } else {
    response = handleUnexpectedError(error)
  }

  return response
}

export const handleAxiosError = (error: AxiosError) => {
  console.log('Service Error: ', error)
  return {
    status: 'error',
    message: error.message
  }
}

export const handleUnexpectedError = (error: unknown) => {
  console.log('Service Error: ', error)
  return {
    status: 'error',
    message: 'Unknown Error'
  }
}
