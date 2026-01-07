import { Link, route } from 'preact-router';
import { SyncStatus } from './SyncStatus';

const TypedLink = Link as any;

export const Navbar = () => {
  const role = localStorage.getItem('user_role');
  
  const handleLogout = () => {
    localStorage.removeItem('token');
    localStorage.removeItem('user_phone');
    localStorage.removeItem('user_role');
    localStorage.removeItem('user_shop_id');
    localStorage.removeItem('user_account');
    localStorage.removeItem('must_change_password');
    window.dispatchEvent(new Event('auth-change'));
    route('/login');
  };

  return (
    <nav className="navbar">
      <div className="nav-logo">StockManager</div>
      <div className="nav-links">
        <TypedLink href="/" activeClassName="active">Dashboard</TypedLink>
        <TypedLink href="/catalogue" activeClassName="active">Catalogue</TypedLink>
        <TypedLink href="/stocks" activeClassName="active">Stocks</TypedLink>
        {role !== 'vendor' && (
          <>
            <TypedLink href="/transfers" activeClassName="active">Transferts</TypedLink>
            <TypedLink href="/shops" activeClassName="active">Boutiques</TypedLink>
            <TypedLink href="/team" activeClassName="active">Équipe</TypedLink>
            <TypedLink href="/stats" activeClassName="active">Statistiques</TypedLink>
            <TypedLink href="/settings" activeClassName="active">Paramètres</TypedLink>
            <TypedLink href="/subscription" activeClassName="active">Abonnement</TypedLink>
          </>
        )}
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
