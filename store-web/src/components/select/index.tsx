type InputFieldProps = {
  label?: string
  id?: string
  value?: string
  placeholder?: string
  required?: boolean
  maxLength?: number
  minLength?: number
  onChange?: (event: React.ChangeEvent<HTMLInputElement>) => void
  children?: React.ReactNode
}

const Select = ({
  label,
  id,
  value,
  onChange,
  placeholder,
  required,
  maxLength,
  minLength,
  children
}: InputFieldProps) => {
  return (
    <div className={label ? 'mb-2' : ''}>
      {label ? (
        <label
          htmlFor={id}
          className="block mb-2 text-sm font-medium text-gray-900"
        >
          {label}
        </label>
      ) : null}

      <select
        id={id}
        className="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-indigo-500 focus:border-indigo-500 block w-full p-2.5"
      >
        {children}
      </select>
    </div>
  )
}

export default Select
