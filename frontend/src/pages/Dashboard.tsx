import { useState, useEffect } from 'preact/hooks';

export const Dashboard = ({ path }: { path?: string }) => {
  const role = localStorage.getItem('user_role');
  const shopId = localStorage.getItem('user_shop_id');
  
  const [stats, setStats] = useState({
    total_stock_value: 0,
    low_stock_alerts: 0,
    total_articles: 0,
    active_shops: 0
  });
  const [loading, setLoading] = useState(true);
  const [articles, setArticles] = useState<any[]>([]);
  const [shops, setShops] = useState<any[]>([]);
  const [searchTerm, setSearchTerm] = useState('');
  const [selectedShopId, setSelectedShopId] = useState(shopId || '');
  const [saleLoading, setSaleLoading] = useState(false);
  const [cart, setCart] = useState<any[]>([]);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [selectedArticle, setSelectedArticle] = useState<any>(null);
  const [quantity, setQuantity] = useState(1);

  const fetchStats = () => {
    let url = `/api/dashboard/stats`;
    if (selectedShopId) url += `?shop_id=${selectedShopId}`;

    fetch(url, {
      headers: { 'Authorization': `Bearer ${localStorage.getItem('token')}` }
    })
    .then(res => res.json())
    .then(setStats)
    .catch(console.error);
  };

  const fetchArticles = () => {
    fetch('/api/articles', {
      headers: { 'Authorization': `Bearer ${localStorage.getItem('token')}` }
    })
    .then(res => res.json())
    .then(setArticles)
    .catch(console.error);
  };

  const fetchShops = () => {
    fetch('/api/shops', {
      headers: { 'Authorization': `Bearer ${localStorage.getItem('token')}` }
    })
    .then(res => res.json())
    .then(setShops)
    .catch(console.error);
  };

  useEffect(() => {
    Promise.all([fetchStats(), fetchArticles(), fetchShops()])
      .finally(() => setLoading(false));
  }, [selectedShopId]);

  const addToCart = () => {
    if (!selectedArticle) return;
    if (quantity <= 0) return;
    if ((selectedArticle.total_stock || 0) < quantity) {
      alert("Stock insuffisant !");
      return;
    }

    const existing = cart.find(item => item.id === selectedArticle.id);
    if (existing) {
      setCart(cart.map(item => item.id === selectedArticle.id ? { ...item, qty: item.qty + quantity } : item));
    } else {
      setCart([...cart, { ...selectedArticle, qty: quantity }]);
    }
    
    setIsModalOpen(false);
    setSelectedArticle(null);
    setQuantity(1);
  };

  const removeFromCart = (id: string) => {
    setCart(cart.filter(item => item.id !== id));
  };

  const processSale = async () => {
    if (!selectedShopId) {
      alert("Veuillez sélectionner une boutique.");
      return;
    }
    if (cart.length === 0) return;

    setSaleLoading(true);
    try {
      // Process articles sequentially for now (could be optimized with a bulk API)
      for (const item of cart) {
        const response = await fetch('/api/stocks/movement', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${localStorage.getItem('token')}`
          },
          body: JSON.stringify({
            shop_id: selectedShopId,
            article_id: item.id,
            type: 'out',
            qty: item.qty,
            reason: 'Vente Panier (Dashboard)'
          })
        });

        if (!response.ok) {
          const data = await response.json();
          throw new Error(`Erreur pour ${item.name}: ${data.error}`);
        }
      }

      alert("Vente enregistrée avec succès !");
      setCart([]);
      fetchStats();
      fetchArticles();
    } catch (err: any) {
      alert(err.message);
    } finally {
      setSaleLoading(false);
    }
  };

  const cartTotal = cart.reduce((total, item) => total + (item.price * item.qty), 0);

  const filteredArticles = articles.filter(a => 
    a.name.toLowerCase().includes(searchTerm.toLowerCase()) || 
    a.sku.toLowerCase().includes(searchTerm.toLowerCase())
  );

  const currentShopName = shops.find(s => s.id === selectedShopId)?.name;

  if (loading) return <div className="loading">Chargement des statistiques...</div>;

  return (
    <div className="page dashboard">
      <div className="section-header" style={{display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '2rem'}}>
        <h1>{role === 'vendor' ? (currentShopName || 'Ma Boutique') : 'Vente Rapide'}</h1>
        
        {role !== 'vendor' && (
          <div className="form-group" style={{margin: 0, width: '250px'}}>
            <select value={selectedShopId} onChange={(e) => setSelectedShopId(e.currentTarget.value)}>
              <option value="">Toutes les boutiques</option>
              {shops.map(s => <option key={s.id} value={s.id}>{s.name}</option>)}
            </select>
          </div>
        )}
      </div>

      <div className="dashboard-layout">
        <aside className="dashboard-sidebar">
          <div className="stats-grid vertical">
            {role !== 'vendor' && (
              <div className="card">
                <h3>Valeur Stock</h3>
                <p className="value">{stats.total_stock_value.toLocaleString()} CFA</p>
              </div>
            )}
            <div className="card highlight">
              <h3>Alertes Stock</h3>
              <p className="value">{stats.low_stock_alerts}</p>
            </div>
            <div className="card">
              <h3>Catalogue</h3>
              <p className="value">{stats.total_articles}</p>
            </div>
          </div>

          <div className="cart-container" style={{position: 'static', margin: 0, width: '100%'}}>
            <div className="cart-title">
              <span>Panier</span>
              <span className="badge badge-info">{cart.length}</span>
            </div>

            <div className="cart-items" style={{maxHeight: '300px', overflowY: 'auto'}}>
              {cart.map(item => (
                <div key={item.id} className="cart-item">
                  <div className="cart-item-info">
                    <h5>{item.name}</h5>
                    <span>{item.qty} x {item.price.toLocaleString()} CFA</span>
                  </div>
                  <button className="btn btn-sm" style={{color: 'var(--error)', background: 'transparent'}} onClick={(e) => { e.stopPropagation(); removeFromCart(item.id); }}>✕</button>
                </div>
              ))}
              {cart.length === 0 && (
                <div style={{textAlign: 'center', padding: '1.5rem', color: 'var(--text-light)', fontSize: '0.85rem'}}>
                  Le panier est vide
                </div>
              )}
            </div>

            <div className="cart-total">
              <span>Total</span>
              <span>{cartTotal.toLocaleString()} CFA</span>
            </div>

            <button 
              className="btn btn-primary btn-block" 
              style={{marginTop: '1.25rem'}} 
              disabled={cart.length === 0 || saleLoading}
              onClick={processSale}
            >
              {saleLoading ? 'Traitement...' : 'Valider la vente'}
            </button>
          </div>
        </aside>

        <main className="dashboard-main">
          <div className="search-bar" style={{marginBottom: '2rem'}}>
            <input 
              type="text" 
              placeholder="Rechercher un article (Nom ou SKU)..." 
              value={searchTerm}
              onInput={(e) => setSearchTerm(e.currentTarget.value)}
              style={{padding: '1rem 1.5rem', fontSize: '1.1rem', borderRadius: '1.25rem', border: '2px solid var(--primary-light)', width: '100%', outline: 'none'}}
            />
          </div>

          <div className="articles-grid" style={{display: 'grid', gridTemplateColumns: 'repeat(auto-fill, minmax(220px, 1fr))', gap: '1.25rem'}}>
            {filteredArticles.map(article => (
              <div key={article.id} className="card article-card" style={{padding: '1.15rem', cursor: 'pointer', display: 'flex', flexDirection: 'column', justifyContent: 'space-between', height: '100%'}} onClick={() => { setSelectedArticle(article); setIsModalOpen(true); }}>
                <div className="article-info">
                  <span className="sku" style={{fontSize: '0.65rem', fontWeight: 800, color: 'var(--text-light)', display: 'block', marginBottom: '0.25rem'}}>{article.sku}</span>
                  <h4 style={{margin: '0 0 0.5rem 0', fontSize: '0.95rem', lineHeight: '1.3'}}>{article.name}</h4>
                  <p style={{fontSize: '1.05rem', fontWeight: 900, color: 'var(--primary)', margin: 0}}>{article.price.toLocaleString()} CFA</p>
                </div>
                <div className="article-stock" style={{marginTop: '0.75rem', display: 'flex', justifyContent: 'space-between', alignItems: 'center'}}>
                  <span className={`badge ${article.total_stock < article.min_threshold ? 'badge-error' : 'badge-success'}`} style={{fontSize: '0.7rem'}}>
                    Stock: {article.total_stock || 0}
                  </span>
                  <button className="btn btn-sm btn-primary" style={{padding: '0.3rem 0.6rem', fontSize: '0.75rem'}}>Vendre</button>
                </div>
              </div>
            ))}
            {filteredArticles.length === 0 && (
              <div style={{gridColumn: '1 / -1', textAlign: 'center', padding: '3rem', color: 'var(--text-light)'}}>
                Aucun article correspondant trouvé.
              </div>
            )}
          </div>
        </main>
      </div>

      {isModalOpen && (
        <div className="modal-overlay" onClick={() => { setIsModalOpen(false); setSelectedArticle(null); }}>
          <div className="modal-content" onClick={(e) => e.stopPropagation()}>
            <h3 style={{marginBottom: '0.5rem', fontSize: '1.25rem'}}>Quantité pour {selectedArticle?.name}</h3>
            <p style={{fontSize: '0.85rem', color: 'var(--text-light)', marginBottom: '1.5rem'}}>
              Stock disponible: <strong style={{color: 'var(--primary)'}}>{selectedArticle?.total_stock || 0}</strong>
            </p>
            <div className="form-group">
              <label style={{fontSize: '0.9rem', fontWeight: 600}}>Quantité à vendre</label>
              <input 
                type="number" 
                value={quantity} 
                onInput={(e) => setQuantity(Number(e.currentTarget.value))} 
                min="1" 
                max={selectedArticle?.total_stock}
                autoFocus
                style={{padding: '0.75rem', fontSize: '1.1rem'}}
              />
            </div>
            <div style={{display: 'flex', gap: '1rem', marginTop: '2rem'}}>
              <button className="btn btn-block" onClick={() => { setIsModalOpen(false); setSelectedArticle(null); }}>Annuler</button>
              <button className="btn btn-primary btn-block" onClick={addToCart}>Ajouter au panier</button>
            </div>
          </div>
        </div>
      )}

    </div>
  );
};
