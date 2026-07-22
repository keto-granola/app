import { Route, Routes } from 'react-router-dom'

import AddProduct from '../islands/products/AddProduct'
import GetProducts from '../islands/products/GetProducts'

function AdminRoutes() {
  return (
    <Routes>
      <Route path="products" element={<GetProducts />} />
      <Route path="products/new" element={<AddProduct />} />
    </Routes>
  )
}

export default AdminRoutes
