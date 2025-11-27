'use client'

import Button from '@/components/button/button'
import InputField from '@/components/input-field'
import { useState } from 'react'
import { GoogleIcon } from './google-icon'

const LoginForm = () => {
  const [form, setForm] = useState({
    username: '',
    password: ''
  })
  const [error, setError] = useState({
    username: '',
    password: ''
  })

  const handleLogin = (event: React.MouseEvent<HTMLButtonElement>) => {
    event.preventDefault()

    const isValidInputs = validateInputs()

    if (!isValidInputs) return
    console.log('form,', form)
  }

  const handleChange = ({
    target: { name, value }
  }: React.ChangeEvent<HTMLInputElement>) => {
    setForm((prev) => ({ ...prev, [name]: value }))
    setError((prev) => ({ ...prev, [name]: '' }))
  }

  const validateInputs = (): boolean => {
    const newError = { username: '', password: '' }
    const { username, password } = form

    if (!username.trim()) newError.username = 'Username is required.'
    if (!password.trim()) newError.password = 'Password is required.'

    setError(newError)

    return !newError.username && !newError.password
  }

  return (
    <div className="max-w-[424px] w-full px-3 flex flex-col gap-8">
      <h1 id="login-form-header" className="font-bold text-3xl">
        Login
      </h1>
      <div id="login-form-container" className="flex flex-col gap-8">
        <div id="login-main-form" className="flex flex-col gap-8">
          <div className="flex flex-col gap-[10px]">
            <div>
              <InputField
                id="login-username"
                label="Username"
                type="text"
                name="username"
                placeholder="Enter your username"
                required
                onChange={handleChange}
              />
              {error.username && (
                <span
                  id="login-username-input-error-txt"
                  className="text-[10px] font-light text-red-500"
                >
                  {error.username}
                </span>
              )}
            </div>
            <div className="flex flex-col">
              <InputField
                id="login-password"
                label="Password"
                type="password"
                name="password"
                placeholder="Enter your password"
                required
                onChange={handleChange}
              />
              {error.password && (
                <span
                  id="login-password-input-error-txt"
                  className="text-[10px] font-light text-red-500"
                >
                  {error.password}
                </span>
              )}
            </div>
            <div className="flex justify-end">
              <button
                id="forget-password-btn"
                type="button"
                className="text-[10px] font-medium text-indigo-600 hover:text-indigo-500 underline flex gap-1 items-center"
                // onClick={onClick}
              >
                Forget Password?
              </button>
            </div>
          </div>
          <Button
            id="login-btn"
            type="button"
            onClick={handleLogin}
            isblock="true"
            size="sm"
          >
            <span id="login-btn-txt" className="font-bold">
              Login
            </span>
          </Button>
        </div>
        <div className="relative w-full">
          <div className="bg-[#F5F5F5] w-full h-[2px]"></div>
          <span className="absolute -top-2 left-[45%] px-1.5 bg-white text-[10px]">
            Or
          </span>
        </div>
        <button
          id="login-with-google-btn"
          className="bg-white text-xs font-medium w-full px-3 py-2 text-gray-900 border border-gray-200 rounded-md hover:bg-slate-50"
        >
          <span className="flex justify-center items-center gap-2">
            <GoogleIcon />
            Log in with Google
          </span>
        </button>
        <div className="flex justify-center gap-1 text-xs">
          <span>Don&apos;t have an account?</span>
          <button
            id="register-btn"
            type="button"
            className="text-indigo-600 hover:text-indigo-500 flex gap-1 items-center"
            // onClick={onClick}
          >
            Register
          </button>
        </div>
      </div>
    </div>
  )
}

export default LoginForm
