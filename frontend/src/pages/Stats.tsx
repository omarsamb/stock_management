import { useState, useEffect } from 'preact/hooks';
import { StockPieChart } from '../components/charts/StockPieChart';
import { MovementBarChart } from '../components/charts/MovementBarChart';
import { LowStockTable } from '../components/charts/LowStockTable';
import { SalesChart } from '../components/charts/SalesChart';

export const Stats = ({ path }: { path?: string }) => {
  const [shops, setShops] = useState<any[]>([]);
  const [selectedShopId, setSelectedShopId] = useState('');
  const [loading, setLoading] = useState(true);
  const [stats, setStats] = useState<any>(null);
  
  // Sales Stats State
  const [salesPeriod, setSalesPeriod] = useState('month');
  const [salesData, setSalesData] = useState<any[]>([]);
  const [salesLoading, setSalesLoading] = useState(false);

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

  const fetchSalesStats = () => {
    setSalesLoading(true);
    let url = `/api/dashboard/sales?period=${salesPeriod}`;
    if (selectedShopId) url += `&shop_id=${selectedShopId}`;

    fetch(url, {
      headers: { 'Authorization': `Bearer ${localStorage.getItem('token')}` }
    })
    .then(res => res.json())
    .then(data => {
      setSalesData(data || []);
      setSalesLoading(false);
    })
    .catch(err => {
      console.error(err);
      setSalesLoading(false);
    });
  };

  useEffect(() => {
    fetchShops();
  }, []);

  useEffect(() => {
    fetchStats();
  }, [selectedShopId]);

  useEffect(() => {
    fetchSalesStats();
  }, [selectedShopId, salesPeriod]);

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
            <div style={{display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '1rem'}}>
              <h3 style={{margin: 0}}>Statistiques de Vente</h3>
              <div className="period-selector">
                {['day', 'week', 'month', 'year'].map(p => (
                  <button 
                    key={p} 
                    className={`btn btn-sm ${salesPeriod === p ? 'btn-primary' : ''}`} 
                    style={{marginLeft: '0.5rem'}}
                    onClick={() => setSalesPeriod(p)}
                  >
                    {p === 'day' ? 'Jour' : p === 'week' ? 'Semaine' : p === 'month' ? 'Mois' : 'Année'}
                  </button>
                ))}
              </div>
            </div>
            
            <div style={{ height: '350px' }}>
              {salesLoading ? (
                <div className="loading">Chargement des ventes...</div>
              ) : (
                <SalesChart data={salesData} period={salesPeriod} />
              )}
            </div>
          </div>

          <div className="stats-charts-grid" style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(350px, 1fr))', gap: '1.5rem', marginTop: '2rem' }}>
            <div className="card">
              <h3>Répartition du Stock</h3>
              <div style={{ height: '300px', display: 'flex', justifyContent: 'center' }}>
                <StockPieChart data={stats.stock_by_category || []} />
              </div>
            </div>
            <div className="card">
              <h3>Mouvements (30 Jours)</h3>
              <div style={{ height: '300px' }}>
                <MovementBarChart data={stats.daily_movements || []} />
              </div>
            </div>
          </div>

          <div className="card" style={{ marginTop: '2rem' }}>
            <h3>Alertes Stock Critique</h3>
            <LowStockTable items={stats.low_stock_items || []} />
          </div>
        </>
      )}
    </div>
  );
};
