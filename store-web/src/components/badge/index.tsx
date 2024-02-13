'use client'

// ----------------------------------------------------------------------

type BadgeProps = {
  total: number
}

const Badge = ({ total }: BadgeProps) => {
  return (
    <span className="relative bg-red-500 text-xs text-white rounded-md p-1.5 left-4 -top-[10px]">
      {total}
    </span>
  )
}

export default Badge
