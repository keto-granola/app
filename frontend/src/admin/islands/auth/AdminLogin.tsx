import { useNavigate } from 'react-router-dom'

function AdminLogin() {
  const navigate = useNavigate()

  const onHandleLogin = () => {
    void navigate('/products')
  }

  return (
    <div>
      <h1>Login</h1>
      <button onClick={onHandleLogin}>Login</button>
    </div>
  )
}

export default AdminLogin
