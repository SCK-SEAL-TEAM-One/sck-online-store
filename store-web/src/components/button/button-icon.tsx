'use client'

// ----------------------------------------------------------------------

type InputFieldProps = React.ButtonHTMLAttributes<HTMLButtonElement> & {
  children: React.ReactNode
  isblock?: boolean
  id?: string
}

const ButtonIcon = (props: InputFieldProps) => {
  return <button {...props}>{props.children}</button>
}

export default ButtonIcon
