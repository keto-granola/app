import type { PropsWithChildren } from 'react'
import { Navigate } from 'react-router-dom'

import { useAuth } from './AuthContext'

function RequireAdmin({ children }: PropsWithChildren) {
  const { token } = useAuth()

  if (!token) {
    return <Navigate to="/login" />
  }

  return children
}

export default RequireAdmin
