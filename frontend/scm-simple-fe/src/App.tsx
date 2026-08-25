import { Navigate, Outlet, Route, Routes } from 'react-router'
import { AppLayout } from './layouts/AppLayout'
import { SearchPage } from './pages/items/SearchPage'
import { ProductsPage } from './pages/items/ProductsPage'
import { SuppliersPage } from './pages/items/SuppliersPage'
import { WarehouseTransferPage } from './pages/transfer/WarehousePage'
import { TechnicianTransferPage } from './pages/transfer/TechnicianPage'
import { CustomerTransferPage } from './pages/transfer/CustomerPage'
import { WarehouseLocationsPage } from './pages/locations/WarehousePage'
import { TechnicianLocationsPage } from './pages/locations/TechnicianPage'
import { CustomerLocationsPage } from './pages/locations/CustomerPage'
import { Login } from './pages/login/Login'

const AUTH_STORAGE_KEY = 'scm-authenticated'

function ProtectedRoutes() {
  const isAuthenticated = sessionStorage.getItem(AUTH_STORAGE_KEY) === 'true'

  return isAuthenticated ? <Outlet /> : <Navigate to="/auth/login" replace />
}

// Declarative mode: routes are described as JSX (<Routes>/<Route>) rather than
// built as a route-object tree with createBrowserRouter (data mode).
function App() {
  return (
    <Routes>
      <Route path="auth">
        <Route path="login" element={<Login />} />
      </Route>

      <Route element={<ProtectedRoutes />}>
        <Route element={<AppLayout />}>
          <Route index element={<Navigate to="/items/search" replace />} />

          <Route path="items">
            <Route path="search" element={<SearchPage />} />
            <Route path="products" element={<ProductsPage />} />
            <Route path="suppliers" element={<SuppliersPage />} />
          </Route>

          <Route path="locations">
            <Route path="warehouse" element={<WarehouseLocationsPage />} />
            <Route path="technician" element={<TechnicianLocationsPage />} />
            <Route path="customer" element={<CustomerLocationsPage />} />
          </Route>

          <Route path="transfer">
            <Route path="warehouse" element={<WarehouseTransferPage />} />
            <Route path="technician" element={<TechnicianTransferPage />} />
            <Route path="customer" element={<CustomerTransferPage />} />
          </Route>

          <Route
            path="*"
            element={
              <div className="flex h-full items-center justify-center text-sm text-ink-soft">
                Page not found.
              </div>
            }
          />
        </Route>
      </Route>

      <Route path="*" element={<Navigate to="/auth/login" replace />} />
    </Routes>
  )
}

export default App
export { AUTH_STORAGE_KEY }
