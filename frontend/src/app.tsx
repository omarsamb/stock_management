import { useState, useEffect } from 'preact/hooks';
import { Router, route } from 'preact-router';
import { Dashboard } from './pages/Dashboard';
import { Catalogue } from './pages/Catalogue';
import { Stocks } from './pages/Stocks';
import { Shops } from './pages/Shops';
import { Users } from './pages/Users';
import { Transfers } from './pages/Transfers';
import { Subscription } from './pages/Subscription';
import { Stats } from './pages/Stats';
import { Settings } from './pages/Settings';
import { ChangePassword } from './pages/ChangePassword';
import { Login } from './pages/Login';
import { Register } from './pages/Register';
import { Navbar } from './components/Navbar';

export function App() {
  const [isAuthenticated, setIsAuthenticated] = useState(!!localStorage.getItem('token'));

  useEffect(() => {
    const applyTheme = () => {
      const savedAccount = localStorage.getItem('user_account');
      if (savedAccount) {
        const account = JSON.parse(savedAccount);
        if (account.primary_color) {
          document.documentElement.style.setProperty('--primary', account.primary_color);
        }
        if (account.background_image) {
          if (account.background_image.startsWith('http')) {
            document.documentElement.style.setProperty('--app-bg-image', `url(${account.background_image})`);
            document.documentElement.style.setProperty('--app-bg-color', 'transparent');
          } else {
            document.documentElement.style.setProperty('--app-bg-color', account.background_image);
            document.documentElement.style.setProperty('--app-bg-image', 'none');
          }
        }
      }
    };

    const checkAuth = () => {
      const token = localStorage.getItem('token');
      setIsAuthenticated(!!token);
      applyTheme();
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
          <Transfers path="/transfers" />
          <Users path="/team" />
          <Subscription path="/subscription" />
          <Stats path="/stats" />
          <Settings path="/settings" />
          <ChangePassword path="/change-password" />
          <Login path="/login" />
          <Register path="/register" />
        </Router>
      </main>
    </div>
  );
}
