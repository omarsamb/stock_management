import { useState, useEffect } from 'preact/hooks';

interface Shop {
  id: string;
  name: string;
  location: string;
}

export const Shops = ({ path }: { path?: string }) => {
  const [shops, setShops] = useState<Shop[]>([]);
  const [name, setName] = useState('');
  const [location, setLocation] = useState('');
  const [error, setError] = useState('');

  const fetchShops = async () => {
    try {
      const response = await fetch('/api/shops', {
        headers: { 'Authorization': `Bearer ${localStorage.getItem('token')}` }
      });
      const data = await response.json();
      setShops(data);
    } catch (err) {
      console.error(err);
    }
  };

  useEffect(() => {
    fetchShops();
  }, []);

  const handleSubmit = async (e: Event) => {
    e.preventDefault();
    setError('');

    try {
      const response = await fetch('/api/shops', {
        method: 'POST',
        headers: { 
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${localStorage.getItem('token')}` 
        },
        body: JSON.stringify({ name, location })
      });

      if (!response.ok) throw new Error('Erreur lors de la création de la boutique');

      setName('');
      setLocation('');
      fetchShops();
    } catch (err: any) {
      setError(err.message);
    }
  };

  return (
    <div className="page shops">
      <h1>Gestion des Boutiques</h1>

      <div className="layout-grid">
        <div className="card">
          <h3>Nouvelle Boutique</h3>
          {error && <div className="alert alert-error">{error}</div>}
          <form onSubmit={handleSubmit}>
            <div className="form-group">
              <label>Nom de la boutique</label>
              <input 
                type="text" 
                value={name} 
                onInput={(e) => setName(e.currentTarget.value)} 
                required 
              />
            </div>
            <div className="form-group">
              <label>Emplacement / Ville</label>
              <input 
                type="text" 
                value={location} 
                onInput={(e) => setLocation(e.currentTarget.value)} 
              />
            </div>
            <button type="submit" className="btn btn-primary">Créer</button>
          </form>
        </div>

        <div className="card">
          <h3>Boutiques Existantes</h3>
          <table className="data-table">
            <thead>
              <tr>
                <th>Nom</th>
                <th>Emplacement</th>
              </tr>
            </thead>
            <tbody>
              {shops.map(shop => (
                <tr key={shop.id}>
                  <td>{shop.name}</td>
                  <td>{shop.location}</td>
                </tr>
              ))}
              {shops.length === 0 && (
                <tr>
                  <td colSpan={2} style={{textAlign: 'center'}}>Aucune boutique</td>
                </tr>
              )}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
};
