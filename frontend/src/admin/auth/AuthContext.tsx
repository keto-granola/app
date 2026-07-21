import { onIdTokenChanged } from 'firebase/auth'
import type { PropsWithChildren } from 'react'
import { createContext, useContext, useEffect, useState } from 'react'

import { firebaseAuth } from '../firebase'

interface AuthContextType {
  token: string | null
}

export const AuthContext = createContext<AuthContextType | undefined>(undefined)

export function AuthProvider({ children }: PropsWithChildren) {
  const [token, setToken] = useState<string | null>(null)

  useEffect(() => {
    const unsubscribe = onIdTokenChanged(firebaseAuth, user => {
      const updateToken = async () => {
        setToken(user ? await user.getIdToken() : null)
      }

      void updateToken()
    })

    return unsubscribe
  }, [])

  return <AuthContext.Provider value={{ token }}>{children}</AuthContext.Provider>
}

export function useAuth(): AuthContextType {
  const context = useContext(AuthContext)

  if (!context) {
    throw new Error('useAuth must be used within AuthProvider')
  }

  return context
}
