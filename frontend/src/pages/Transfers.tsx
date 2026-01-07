import { useState, useEffect } from 'preact/hooks';

interface Transfer {
  id: string;
  from_shop_name: string;
  to_shop_name: string;
  article_name: string;
  qty: number;
  status: string;
  created_at: string;
}

export const Transfers = ({ path }: { path?: string }) => {
  const [transfers, setTransfers] = useState<Transfer[]>([]);
  const [loading, setLoading] = useState(true);

  const fetchTransfers = async () => {
    try {
      const response = await fetch('/api/transfers', {
        headers: { 'Authorization': `Bearer ${localStorage.getItem('token')}` }
      });
      const data = await response.json();
      setTransfers(data || []);
    } catch (err) {
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchTransfers();
  }, []);

  const handleReceive = async (id: string) => {
    try {
      const response = await fetch(`/api/transfers/${id}/receive`, {
        method: 'POST',
        headers: { 'Authorization': `Bearer ${localStorage.getItem('token')}` }
      });
      if (response.ok) {
        fetchTransfers();
      }
    } catch (err) {
      console.error(err);
    }
  };

  if (loading) return <div className="loading">Chargement des transferts...</div>;

  return (
    <div className="page transfers">
      <h1>Transferts de Stock</h1>
      
      <div className="card">
        <h3>Historique des transferts</h3>
        <table className="data-table">
          <thead>
            <tr>
              <th>Date</th>
              <th>Article</th>
              <th>De</th>
              <th>Vers</th>
              <th>Qté</th>
              <th>Statut</th>
              <th>Action</th>
            </tr>
          </thead>
          <tbody>
            {transfers.map(t => (
              <tr key={t.id}>
                <td>{new Date(t.created_at).toLocaleDateString()}</td>
                <td>{t.article_name}</td>
                <td>{t.from_shop_name}</td>
                <td>{t.to_shop_name}</td>
                <td>{t.qty}</td>
                <td>
                  <span className={`badge badge-${t.status === 'pending' ? 'info' : 'success'}`}>
                    {t.status === 'pending' ? 'En attente' : 'Terminé'}
                  </span>
                </td>
                <td>
                  {t.status === 'pending' && (
                    <button onClick={() => handleReceive(t.id)} className="btn btn-primary btn-sm">
                      Confirmer Réception
                    </button>
                  )}
                </td>
              </tr>
            ))}
            {transfers.length === 0 && (
              <tr>
                <td colSpan={7} style={{textAlign: 'center'}}>Aucun transfert récent</td>
              </tr>
            )}
          </tbody>
        </table>
      </div>
    </div>
  );
};
