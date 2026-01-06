import { useState, useEffect } from 'preact/hooks';

export const Dashboard = ({ path }: { path?: string }) => {
  const [stats, setStats] = useState({
    total_stock_value: 0,
    low_stock_alerts: 0,
    total_articles: 0,
    active_shops: 0
  });
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetch('/api/dashboard/stats', {
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    })
    .then(res => res.json())
    .then(data => {
      setStats(data);
      setLoading(false);
    })
    .catch(err => {
      console.error(err);
      setLoading(false);
    });
  }, []);

  if (loading) return <div className="loading">Chargement des statistiques...</div>;

  return (
    <div className="page dashboard">
      <h1>Tableau de Bord Global</h1>
      <div className="stats-grid">
        <div className="card">
          <h3>Valeur Totale Stock</h3>
          <p className="value">{stats.total_stock_value.toLocaleString()} CFA</p>
        </div>
        <div className="card highlight">
          <h3>Articles en Alerte</h3>
          <p className="value">{stats.low_stock_alerts}</p>
        </div>
        <div className="card">
          <h3>Articles au Catalogue</h3>
          <p className="value">{stats.total_articles}</p>
        </div>
        <div className="card">
          <h3>Boutiques Actives</h3>
          <p className="value">{stats.active_shops}</p>
        </div>
      </div>
    </div>
  );
};
