import { BrowserRouter, Route, Routes } from 'react-router-dom'

import mountIsland from '../../../mount'
import { AuthProvider } from '../../auth/AuthContext'
import RequireAdmin from '../../auth/RequireAdmin'
import AdminRoutes from '../../routes/AdminRoutes'
import AdminLogin from '../auth/AdminLogin'

function AdminApp() {
  return (
    <BrowserRouter basename="/admin">
      <AuthProvider>
        <Routes>
          <Route path="/login" element={<AdminLogin />} />
          <Route
            path="*"
            element={
              <RequireAdmin>
                <AdminRoutes />
              </RequireAdmin>
            }
          />
        </Routes>
      </AuthProvider>
    </BrowserRouter>
  )
}

const el = document.getElementById('admin-dashboard')
if (!el) {
  throw new Error('admin-dashboard not found')
}

mountIsland({ el, Component: AdminApp, props: {} })
