import { Link, route } from 'preact-router';
import { SyncStatus } from './SyncStatus';

const TypedLink = Link as any;

export const Navbar = () => {
  const handleLogout = () => {
    localStorage.removeItem('token');
    localStorage.removeItem('user_phone');
    window.dispatchEvent(new Event('auth-change'));
    route('/login');
  };

  return (
    <nav className="navbar">
      <div className="nav-logo">StockManager</div>
      <div className="nav-links">
        <TypedLink href="/">Dashboard</TypedLink>
        <TypedLink href="/catalogue">Catalogue</TypedLink>
        <TypedLink href="/stocks">Stocks</TypedLink>
        <TypedLink href="/shops">Boutiques</TypedLink>
      </div>
      <div className="nav-right">
        <SyncStatus />
        <button onClick={handleLogout} className="btn btn-link" style={{color: 'var(--error)', fontSize: '0.875rem'}}>
          Déconnexion
        </button>
      </div>
    </nav>
  );
};
