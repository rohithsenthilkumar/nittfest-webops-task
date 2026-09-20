import { useState, useEffect } from 'react'
import axios from 'axios'
import './App.css'

function App() {
  const [users, setUsers] = useState([])
  const [orders, setOrders] = useState([])
  const [activeTab, setActiveTab] = useState('users')
  
  const [userServiceStatus, setUserServiceStatus] = useState({ loading: true, status: 'unknown' })
  const [orderServiceStatus, setOrderServiceStatus] = useState({ loading: true, status: 'unknown' })
  
  const [loadingUsers, setLoadingUsers] = useState(false)
  const [loadingOrders, setLoadingOrders] = useState(false)
  const [errorUsers, setErrorUsers] = useState(null)
  const [errorOrders, setErrorOrders] = useState(null)

  const userServiceUrl = '/api/users'
  const orderServiceUrl = '/api/orders'

  const checkHealth = async () => {
    setUserServiceStatus(prev => ({ ...prev, loading: true }))
    setOrderServiceStatus(prev => ({ ...prev, loading: true }))
    
    try {
      await axios.get(`/api/user-health`, { timeout: 3000 })
      setUserServiceStatus({ loading: false, status: 'online' })
    } catch {
      setUserServiceStatus({ loading: false, status: 'offline' })
    }

    try {
      await axios.get(`/api/order-health`, { timeout: 3000 })
      setOrderServiceStatus({ loading: false, status: 'online' })
    } catch {
      setOrderServiceStatus({ loading: false, status: 'offline' })
    }
  }

  const fetchUsers = async () => {
    setLoadingUsers(true)
    setErrorUsers(null)
    try {
      const res = await axios.get(userServiceUrl)
      setUsers(res.data)
      setLoadingUsers(false)
    } catch (err) {
      setErrorUsers(err.message)
      setLoadingUsers(false)
    }
  }

  const fetchOrders = async () => {
    setLoadingOrders(true)
    setErrorOrders(null)
    try {
      const res = await axios.get(orderServiceUrl)
      setOrders(res.data)
      setLoadingOrders(false)
    } catch (err) {
      setErrorOrders(err.message)
      setLoadingOrders(false)
    }
  }

  useEffect(() => {
    checkHealth()
    fetchUsers()
    fetchOrders()
  }, [])

  const refreshAll = () => {
    checkHealth()
    fetchUsers()
    fetchOrders()
  }

  return (
    <div className="dashboard-container">
      <header className="dashboard-header">
        <div>
          <h1>Nittfest Webops</h1>
          <p className="subtitle">Service Integration Dashboard</p>
        </div>
        <button className="refresh-btn" onClick={refreshAll}>
          Refresh
        </button>
      </header>

      <section className="status-section">
        <div className="status-item">
          <span className="status-label">User Service Status:</span>
          <span className={`status-value ${userServiceStatus.status}`}>
            {userServiceStatus.loading ? 'Checking' : userServiceStatus.status.toUpperCase()}
          </span>
        </div>
        <div className="status-item">
          <span className="status-label">Order Service Status:</span>
          <span className={`status-value ${orderServiceStatus.status}`}>
            {orderServiceStatus.loading ? 'Checking' : orderServiceStatus.status.toUpperCase()}
          </span>
        </div>
      </section>

      <main className="main-content">
        <div className="tabs-header">
          <button 
            className={`tab-btn ${activeTab === 'users' ? 'active' : ''}`}
            onClick={() => setActiveTab('users')}
          >
            Users
          </button>
          <button 
            className={`tab-btn ${activeTab === 'orders' ? 'active' : ''}`}
            onClick={() => setActiveTab('orders')}
          >
            Orders
          </button>
        </div>

        <div className="tab-content">
          {activeTab === 'users' ? (
            <div className="table-wrapper">
              {loadingUsers && <p className="loading-text">Loading users...</p>}
              {errorUsers && (
                <div className="error-box">
                  <p>Failed to load Users from {userServiceUrl}</p>
                  <button className="retry-btn" onClick={fetchUsers}>Retry</button>
                </div>
              )}
              {!loadingUsers && !errorUsers && users.length === 0 && <p className="empty-text">No users found</p>}
              {!loadingUsers && !errorUsers && users.length > 0 && (
                <table className="data-table">
                  <thead>
                    <tr>
                      <th>ID</th>
                      <th>Name</th>
                      <th>Email</th>
                      <th>Role</th>
                    </tr>
                  </thead>
                  <tbody>
                    {users.map(user => (
                      <tr key={user.id}>
                        <td>{user.id}</td>
                        <td>{user.name}</td>
                        <td>{user.email}</td>
                        <td>{user.role}</td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              )}
            </div>
          ) : (
            <div className="table-wrapper">
              {loadingOrders && <p className="loading-text">Loading orders...</p>}
              {errorOrders && (
                <div className="error-box">
                  <p>Failed to load Orders from {orderServiceUrl}</p>
                  <button className="retry-btn" onClick={fetchOrders}>Retry</button>
                </div>
              )}
              {!loadingOrders && !errorOrders && orders.length === 0 && <p className="empty-text">No orders found</p>}
              {!loadingOrders && !errorOrders && orders.length > 0 && (
                <table className="data-table">
                  <thead>
                    <tr>
                      <th>Order ID</th>
                      <th>Item</th>
                      <th>Quantity</th>
                      <th>User ID</th>
                      <th>User Name</th>
                      <th>Status</th>
                    </tr>
                  </thead>
                  <tbody>
                    {orders.map(order => (
                      <tr key={order.id}>
                        <td>{order.id}</td>
                        <td>{order.item}</td>
                        <td>{order.quantity}</td>
                        <td>{order.user_id}</td>
                        <td>{order.user_name || 'Loading'}</td>
                        <td>{order.status}</td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              )}
            </div>
          )}
        </div>
      </main>
    </div>
  )
}

export default App
