import { useState, useEffect } from 'preact/hooks';

interface Article {
  id: string;
  sku: string;
  name: string;
  description: string;
  min_threshold: number;
  price: number;
  total_stock?: number;
}

export const Catalogue = ({ path }: { path?: string }) => {
  const [articles, setArticles] = useState<Article[]>([]);
  const [showForm, setShowForm] = useState(false);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [editingArticle, setEditingArticle] = useState<Article | null>(null);

  // Form state
  const [name, setName] = useState('');
  const [sku, setSku] = useState(''); // Optional, backend generates if empty
  const [description, setDescription] = useState('');
  const [price, setPrice] = useState(0);
  const [minThreshold, setMinThreshold] = useState(5);

  const resetForm = () => {
    setName('');
    setSku('');
    setDescription('');
    setPrice(0);
    setMinThreshold(5);
    setEditingArticle(null);
    setShowForm(false);
  };

  const handleEdit = (article: Article) => {
    setEditingArticle(article);
    setName(article.name);
    setSku(article.sku);
    setDescription(article.description);
    setPrice(article.price);
    setMinThreshold(article.min_threshold);
    setShowForm(true);
    window.scrollTo({ top: 0, behavior: 'smooth' });
  };

  const fetchArticles = async () => {
    setLoading(true);
    try {
      const response = await fetch('/api/articles', {
        headers: { 'Authorization': `Bearer ${localStorage.getItem('token')}` }
      });
      if (!response.ok) throw new Error('Erreur réseau');
      const data = await response.json();
      setArticles(data);
    } catch (err: any) {
      setError('Impossible de charger les articles');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchArticles();
  }, []);

  const handleSubmit = async (e: Event) => {
    e.preventDefault();
    setError('');

    const url = editingArticle ? `/api/articles/${editingArticle.id}` : '/api/articles';
    const method = editingArticle ? 'PUT' : 'POST';

    try {
      const response = await fetch(url, {
        method,
        headers: { 
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${localStorage.getItem('token')}` 
        },
        body: JSON.stringify({ 
          name, 
          sku, 
          description, 
          price: Number(price), 
          min_threshold: Number(minThreshold) 
        })
      });

      if (!response.ok) {
        const data = await response.json();
        throw new Error(data.error || 'Une erreur est survenue');
      }

      resetForm();
      fetchArticles();
    } catch (err: any) {
      setError(err.message);
    }
  };

  const handleImport = async (e: any) => {
    const file = e.target.files[0];
    if (!file) return;

    const formData = new FormData();
    formData.append('file', file);

    try {
      const response = await fetch('/api/articles/import', {
        method: 'POST',
        headers: { 
          'Authorization': `Bearer ${localStorage.getItem('token')}` 
        },
        body: formData
      });

      if (!response.ok) throw new Error('Erreur lors de l\'importation');
      
      const data = await response.json();
      alert(`${data.imported} articles importés avec succès !`);
      fetchArticles();
    } catch (err: any) {
      setError(err.message);
    }
  };

  return (
    <div className="page catalogue">
      <div className="section-header" style={{display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '2rem'}}>
        <h1>Catalogue d'Articles</h1>
        <div className="actions" style={{display: 'flex', gap: '1rem'}}>
          <button className="btn btn-primary" onClick={() => { if(showForm) resetForm(); else setShowForm(true); }}>
            {showForm ? 'Annuler' : 'Ajouter un Article'}
          </button>
          <label className="btn" style={{cursor: 'pointer', background: '#e2e8f0'}}>
            Importer CSV
            <input type="file" accept=".csv" onChange={handleImport} style={{display: 'none'}} />
          </label>
        </div>
      </div>

      {error && <div className="alert alert-error">{error}</div>}

      {showForm && (
        <div className="card" style={{marginBottom: '2rem'}}>
          <h3>{editingArticle ? 'Modifier l\'Article' : 'Nouvel Article'}</h3>
          <form onSubmit={handleSubmit} className="layout-grid">
            <div className="form-group">
              <label>Nom de l'article</label>
              <input type="text" value={name} onInput={(e) => setName(e.currentTarget.value)} required />
            </div>
            <div className="form-group">
              <label>SKU {editingArticle && '(non modifiable)'}</label>
              <input type="text" value={sku} onInput={(e) => setSku(e.currentTarget.value)} placeholder="Ex: ART-001" disabled={!!editingArticle} />
            </div>
            <div className="form-group">
              <label>Description</label>
              <input type="text" value={description} onInput={(e) => setDescription(e.currentTarget.value)} />
            </div>
            <div className="form-group">
              <label>Prix de vente (CFA)</label>
              <input type="number" value={price} onInput={(e) => setPrice(Number(e.currentTarget.value))} required min="0" />
            </div>
            <div className="form-group">
              <label>Seuil d'alerte stock</label>
              <input type="number" value={minThreshold} onInput={(e) => setMinThreshold(Number(e.currentTarget.value))} required min="1" />
            </div>
            <div style={{gridColumn: '1 / -1'}}>
              <button type="submit" className="btn btn-primary">
                {editingArticle ? 'Mettre à jour l\'article' : 'Enregistrer l\'article'}
              </button>
            </div>
          </form>
        </div>
      )}

      <div className="card">
        {loading ? (
          <div className="loading">Chargement des articles...</div>
        ) : (
          <table className="data-table">
            <thead>
              <tr>
                <th>SKU</th>
                <th>Nom</th>
                <th>Prix</th>
                <th>Seuil</th>
                <th>Stock Global</th>
                <th style={{textAlign: 'right'}}>Actions</th>
              </tr>
            </thead>
            <tbody>
              {articles.map(article => (
                <tr key={article.id}>
                  <td style={{fontWeight: '600', color: '#64748b'}}>{article.sku}</td>
                  <td>{article.name}</td>
                  <td>{article.price.toLocaleString()} CFA</td>
                  <td>{article.min_threshold}</td>
                  <td style={{fontWeight: '700', color: (article.total_stock || 0) < article.min_threshold ? 'var(--error)' : 'var(--success)'}}>
                    {article.total_stock || 0}
                  </td>
                  <td style={{textAlign: 'right'}}>
                    <button className="btn btn-sm" style={{background: 'var(--primary-light)', color: 'var(--primary)'}} onClick={() => handleEdit(article)}>
                      Modifier
                    </button>
                  </td>
                </tr>
              ))}
              {articles.length === 0 && (
                <tr>
                  <td colSpan={6} style={{textAlign: 'center', padding: '2rem'}}>Aucun article dans le catalogue</td>
                </tr>
              )}
            </tbody>
          </table>
        )}
      </div>
    </div>
  );
};
