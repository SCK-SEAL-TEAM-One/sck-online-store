'use client'

// ----------------------------------------------------------------------

type BadgeProps = {
  id?: string
  total: number
}

const Badge = ({ id, total }: BadgeProps) => {
  return (
    <span
      id={id}
      className="relative bg-red-500 text-xs text-white rounded-md p-1.5 left-4 -top-[10px]"
    >
      {total}
    </span>
  )
}

export default Badge
