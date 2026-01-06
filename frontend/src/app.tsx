import { useState, useEffect } from 'preact/hooks';
import { Router, route } from 'preact-router';
import { Dashboard } from './pages/Dashboard';
import { Catalogue } from './pages/Catalogue';
import { Stocks } from './pages/Stocks';
import { Shops } from './pages/Shops';
import { Login } from './pages/Login';
import { Register } from './pages/Register';
import { Navbar } from './components/Navbar';

export function App() {
  const [isAuthenticated, setIsAuthenticated] = useState(!!localStorage.getItem('token'));

  useEffect(() => {
    const checkAuth = () => {
      const token = localStorage.getItem('token');
      setIsAuthenticated(!!token);
      if (!token && window.location.pathname !== '/login' && window.location.pathname !== '/register') {
        route('/login');
      }
    };

    window.addEventListener('auth-change', checkAuth);
    checkAuth();

    return () => window.removeEventListener('auth-change', checkAuth);
  }, []);

  return (
    <div className="app-container">
      {isAuthenticated && <Navbar />}
      <main className={isAuthenticated ? "content" : "auth-content"}>
        <Router>
          <Dashboard path="/" />
          <Catalogue path="/catalogue" />
          <Stocks path="/stocks" />
          <Shops path="/shops" />
          <Login path="/login" />
          <Register path="/register" />
        </Router>
      </main>
    </div>
  );
}
