import { useState, useEffect } from 'preact/hooks';

export const Stats = ({ path }: { path?: string }) => {
  const [shops, setShops] = useState<any[]>([]);
  const [selectedShopId, setSelectedShopId] = useState('');
  const [loading, setLoading] = useState(true);
  const [stats, setStats] = useState<any>(null);

  const fetchShops = () => {
    fetch('/api/shops', {
      headers: { 'Authorization': `Bearer ${localStorage.getItem('token')}` }
    })
    .then(res => res.json())
    .then(setShops)
    .catch(console.error);
  };

  const fetchStats = () => {
    setLoading(true);
    let url = `/api/dashboard/stats`;
    if (selectedShopId) url += `?shop_id=${selectedShopId}`;

    fetch(url, {
      headers: { 'Authorization': `Bearer ${localStorage.getItem('token')}` }
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
  };

  useEffect(() => {
    fetchShops();
  }, []);

  useEffect(() => {
    fetchStats();
  }, [selectedShopId]);

  return (
    <div className="page stats">
      <div className="section-header" style={{display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '2rem'}}>
        <h1>Statistiques & Analyses</h1>
        
        <div className="form-group" style={{margin: 0, width: '250px'}}>
          <select value={selectedShopId} onChange={(e) => setSelectedShopId(e.currentTarget.value)}>
            <option value="">Toutes les boutiques</option>
            {shops.map(s => <option key={s.id} value={s.id}>{s.name}</option>)}
          </select>
        </div>
      </div>

      {loading ? (
        <div className="loading">Chargement des données...</div>
      ) : stats && (
        <>
          <div className="stats-grid">
            <div className="card">
              <h3>Valeur du Stock</h3>
              <p className="value" style={{color: 'var(--primary)'}}>{stats.total_stock_value.toLocaleString()} CFA</p>
              <p style={{fontSize: '0.8rem', color: 'var(--text-light)', marginTop: '0.5rem'}}>Valeur marchande actuelle</p>
            </div>
            <div className="card">
              <h3>Alertes Stock Bas</h3>
              <p className="value" style={{color: 'var(--error)'}}>{stats.low_stock_alerts}</p>
              <p style={{fontSize: '0.8rem', color: 'var(--text-light)', marginTop: '0.5rem'}}>Articles nécessitant un réapprovisionnement</p>
            </div>
            <div className="card">
              <h3>Diversité Produits</h3>
              <p className="value">{stats.total_articles}</p>
              <p style={{fontSize: '0.8rem', color: 'var(--text-light)', marginTop: '0.5rem'}}>Nombre de références au catalogue</p>
            </div>
            <div className="card">
              <h3>Points de Vente</h3>
              <p className="value">{selectedShopId ? '1' : stats.active_shops}</p>
              <p style={{fontSize: '0.8rem', color: 'var(--text-light)', marginTop: '0.5rem'}}>Boutiques actives suivies</p>
            </div>
          </div>

          <div className="card" style={{marginTop: '2rem'}}>
            <h3>Analyse des Performances</h3>
            <div style={{padding: '2rem', textAlign: 'center', background: '#f8fafc', borderRadius: '1rem', marginTop: '1rem'}}>
              <p style={{color: 'var(--text-light)'}}>Graphiques et tendances détaillées bientôt disponibles dans la version Pro.</p>
            </div>
          </div>
        </>
      )}
    </div>
  );
};
