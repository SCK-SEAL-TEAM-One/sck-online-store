'use client'

// ----------------------------------------------------------------------

type InputFieldProps = React.InputHTMLAttributes<HTMLInputElement> & {
  label?: string
  id?: string
}

const InputField = (props: InputFieldProps) => {
  return (
    <div className={props.label ? 'mb-2' : ''}>
      {props.label ? (
        <label
          htmlFor={props.id}
          className="block mb-2 text-sm font-medium text-gray-900"
        >
          {props.label}
        </label>
      ) : null}

      <input
        className="bg-white text-sm w-full px-3 py-2 border border-gray-200 rounded-md focus:outline-none focus:border-indigo-500 transition-colors"
        {...props}
      />
    </div>
  )
}

export default InputField
