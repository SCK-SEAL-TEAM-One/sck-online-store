'use client'

// ----------------------------------------------------------------------

type ButtonLinkProps = React.AnchorHTMLAttributes<HTMLAnchorElement> & {
  children: React.ReactNode
}

const ButtonLink = (props: ButtonLinkProps) => {
  return (
    <a
      {...props}
      className={`${props.className} flex items-center justify-center rounded-md border border-transparent bg-indigo-600 px-6 py-3 text-base font-medium text-white shadow-sm hover:bg-indigo-700`}
    >
      {props.children}
    </a>
  )
}

export default ButtonLink
