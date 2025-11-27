'use client'

import AuthLayout from '@/layouts/common/auth'

// ----------------------------------------------------------------------

type Props = {
  children: React.ReactNode
}

export default function ProductLayout({ children }: Props) {
  return <AuthLayout>{children}</AuthLayout>
}
