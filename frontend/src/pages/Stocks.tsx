import { useState, useEffect } from 'preact/hooks';
import { recordMovement } from '../services/api';

interface Article {
  id: string;
  name: string;
}

interface Shop {
  id: string;
  name: string;
}

export const Stocks = ({ path }: { path?: string }) => {
  const [articles, setArticles] = useState<Article[]>([]);
  const [shops, setShops] = useState<Shop[]>([]);
  const [selectedShop, setSelectedShop] = useState('');
  const [stockLevels, setStockLevels] = useState<any[]>([]);
  const [moveType, setMoveType] = useState<'in' | 'out' | 'transfer' | 'adjust'>('in');
  
  // Form fields
  const [articleId, setArticleId] = useState('');
  const [shopId, setShopId] = useState('');
  const [toShopId, setToShopId] = useState('');
  const [qty, setQty] = useState(0);
  const [reason, setReason] = useState('');
  
  const [message, setMessage] = useState({ type: '', text: '' });

  const role = localStorage.getItem('user_role');
  const userShopId = localStorage.getItem('user_shop_id');

  useEffect(() => {
    // If vendor, auto-select shop
    if (role === 'vendor' && userShopId) {
      setShopId(userShopId);
      setSelectedShop(userShopId);
    }

    // Fetch shops
    fetch('/api/shops', {
      headers: { 'Authorization': `Bearer ${localStorage.getItem('token')}` }
    })
    .then(res => res.json())
    .then(setShops)
    .catch(err => console.error('Fetch shops failed', err));
  }, []);

  useEffect(() => {
    // Fetch articles (filtered by shop if not 'in' movement)
    let url = '/api/articles';
    if (shopId && moveType !== 'in') {
      url += `?shop_id=${shopId}`;
    }

    fetch(url, {
      headers: { 'Authorization': `Bearer ${localStorage.getItem('token')}` }
    })
    .then(res => res.json())
    .then(setArticles)
    .catch(err => console.error('Fetch articles failed', err));
  }, [shopId, moveType]);

  const fetchStockLevels = async () => {
    if (!selectedShop) {
      setStockLevels([]);
      return;
    }
    
    try {
      const res = await fetch(`/api/stocks/levels?shop_id=${selectedShop}`, {
        headers: { 'Authorization': `Bearer ${localStorage.getItem('token')}` }
      });
      const data = await res.json();
      setStockLevels(Array.isArray(data) ? data : []);
    } catch (err) {
      console.error('Fetch stock levels failed', err);
    }
  };

  useEffect(() => {
    fetchStockLevels();
  }, [selectedShop]);

  const handleMovement = async (e: Event) => {
    e.preventDefault();
    setMessage({ type: '', text: '' });

    if (moveType === 'transfer') {
      try {
        const response = await fetch('/api/transfers', {
          method: 'POST',
          headers: { 
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${localStorage.getItem('token')}` 
          },
          body: JSON.stringify({
            from_shop_id: shopId,
            to_shop_id: toShopId,
            article_id: articleId,
            qty: Number(qty),
            reason
          })
        });
        const data = await response.json();
        if (!response.ok) throw new Error(data.error);
        setMessage({ type: 'success', text: 'Transfert initié !' });
        if (shopId === selectedShop) fetchStockLevels();
      } catch (err: any) {
        setMessage({ type: 'error', text: err.message });
      }
    } else {
      try {
        const res = await recordMovement({
          shop_id: shopId,
          article_id: articleId,
          type: moveType,
          qty: Number(qty),
          reason
        });
        setMessage({ type: res.offline ? 'info' : 'success', text: res.message || 'Mouvement enregistré !' });
        if (shopId === selectedShop) fetchStockLevels();
      } catch (err: any) {
        setMessage({ type: 'error', text: err.message });
      }
    }
    
    // Reset form
    setQty(0);
    setReason('');
  };

  return (
    <div className="page stocks">
      <h1>Gestion des Stocks</h1>

      <div className="layout-grid">
        <div className="card">
          <h3>Nouveau Mouvement</h3>
          
          <div className="tab-group">
            <button className={`tab ${moveType === 'in' ? 'active' : ''}`} onClick={() => setMoveType('in')}>Entrée</button>
            <button className={`tab ${moveType === 'out' ? 'active' : ''}`} onClick={() => setMoveType('out')}>Sortie</button>
            <button className={`tab ${moveType === 'transfer' ? 'active' : ''}`} onClick={() => setMoveType('transfer')}>Transfert</button>
            <button className={`tab ${moveType === 'adjust' ? 'active' : ''}`} onClick={() => setMoveType('adjust')}>Ajustement</button>
          </div>

          {message.text && (
            <div className={`alert alert-${message.type}`}>
              {message.text}
            </div>
          )}

          <form onSubmit={handleMovement}>
            {/* Shop Selector - Hidden for Vendor */}
            {role !== 'vendor' && (
              <div className="form-group">
                <label>{moveType === 'transfer' ? 'Boutique Source' : 'Boutique'}</label>
                <select value={shopId} onChange={(e) => {
                  setShopId(e.currentTarget.value);
                  setArticleId('');
                }} required>
                  <option value="">Sélectionner une boutique</option>
                  {shops.map(s => <option key={s.id} value={s.id}>{s.name}</option>)}
                </select>
              </div>
            )}

            <div className="form-group">
              <label>Article</label>
              <select value={articleId} onChange={(e) => setArticleId(e.currentTarget.value)} required>
                <option value="">Sélectionner un article</option>
                {articles.map(a => <option key={a.id} value={a.id}>{a.name}</option>)}
              </select>
            </div>

            {moveType === 'transfer' && (
              <div className="form-group">
                <label>Boutique Destination</label>
                <select value={toShopId} onChange={(e) => setToShopId(e.currentTarget.value)} required>
                  <option value="">Sélectionner la destination</option>
                  {shops.filter(s => s.id !== shopId).map(s => (
                    <option key={s.id} value={s.id}>{s.name}</option>
                  ))}
                </select>
              </div>
            )}

            <div className="form-group">
              <label>Quantité</label>
              <input type="number" value={qty} onInput={(e) => setQty(Number(e.currentTarget.value))} required min="1" />
            </div>

            <div className="form-group">
              <label>Raison / Commentaire</label>
              <input type="text" value={reason} onInput={(e) => setReason(e.currentTarget.value)} placeholder="Ex: Vente comptoir, Arrivage fournisseur..." />
            </div>

            <button type="submit" className="btn btn-primary btn-block">Enregistrer</button>
          </form>
        </div>

        <div className="card">
          <h3>Niveaux de Stock</h3>
          {/* Shop Filter - Hidden for Vendor */}
          {role !== 'vendor' && (
            <div className="form-group">
              <label>Filtrer par boutique</label>
              <select value={selectedShop} onChange={(e) => setSelectedShop(e.currentTarget.value)}>
                <option value="">Sélectionner une boutique</option>
                {shops.map(s => <option key={s.id} value={s.id}>{s.name}</option>)}
              </select>
            </div>
          )}
          
          {selectedShop ? (
            <div className="table-responsive">
              <table className="data-table">
                <thead>
                  <tr>
                    <th>Article</th>
                    <th>Quantité</th>
                    <th>Mise à jour</th>
                  </tr>
                </thead>
                <tbody>
                  {stockLevels.map(level => (
                    <tr key={level.article_id}>
                      <td style={{fontWeight: '600'}}>{level.Article?.name}</td>
                      <td style={{fontWeight: '700', color: 'var(--primary)'}}>{level.quantity}</td>
                      <td style={{fontSize: '0.8rem', color: '#666'}}>
                        {new Date(level.updated_at).toLocaleString()}
                      </td>
                    </tr>
                  ))}
                  {stockLevels.length === 0 && (
                    <tr>
                      <td colSpan={3} style={{textAlign: 'center', padding: '2rem', color: 'var(--text-light)'}}>Aucun stock enregistré pour cette boutique</td>
                    </tr>
                  )}
                </tbody>
              </table>
            </div>
          ) : (
            <p style={{color: '#666', fontSize: '0.9rem', textAlign: 'center', padding: '2rem'}}>Sélectionnez une boutique pour voir les niveaux de stock.</p>
          )}
        </div>
      </div>
    </div>
  );
};
