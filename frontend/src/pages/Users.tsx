import { useState, useEffect } from 'preact/hooks';

interface User {
  id: string;
  phone: string;
  role: string;
  first_name?: string;
  last_name?: string;
  shop_id?: string;
  is_phone_verified?: boolean;
}

export const Users = ({ path }: { path?: string }) => {
  const [firstName, setFirstName] = useState('');
  const [lastName, setLastName] = useState('');
  const [phone, setPhone] = useState('');
  const [role, setRole] = useState('vendor');
  const [message, setMessage] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const [shops, setShops] = useState<any[]>([]);
  const [team, setTeam] = useState<User[]>([]);
  const [selectedShopId, setSelectedShopId] = useState('');
  const [tempPassword, setTempPassword] = useState('');
  const [editingUser, setEditingUser] = useState<User | null>(null);

  const fetchTeam = async () => {
    try {
      const response = await fetch('/api/users', {
        headers: { 'Authorization': `Bearer ${localStorage.getItem('token')}` }
      });
      if (response.ok) {
        const data = await response.json();
        setTeam(data);
      }
    } catch (err) {
      console.error(err);
    }
  };

  useEffect(() => {
    fetch('/api/shops', {
      headers: { 'Authorization': `Bearer ${localStorage.getItem('token')}` }
    })
    .then(res => res.json())
    .then(setShops)
    .catch(err => console.error(err));

    fetchTeam();
  }, []);

  const roles = [
    { value: 'admin', label: 'Administrateur' },
    { value: 'manager', label: 'Gestionnaire' },
    { value: 'vendor', label: 'Vendeur' }
  ];

  const handleEdit = (user: User) => {
    setEditingUser(user);
    setFirstName(user.first_name || '');
    setLastName(user.last_name || '');
    setPhone(user.phone);
    setRole(user.role);
    setSelectedShopId(user.shop_id || '');
    window.scrollTo({ top: 0, behavior: 'smooth' });
  };

  const cancelEdit = () => {
    setEditingUser(null);
    setFirstName('');
    setLastName('');
    setPhone('');
    setRole('vendor');
    setSelectedShopId('');
  };

  const handleDelete = async (userId: string) => {
    if (!confirm('Êtes-vous sûr de vouloir supprimer ce collaborateur ?')) return;

    try {
      const response = await fetch(`/api/users/${userId}`, {
        method: 'DELETE',
        headers: { 'Authorization': `Bearer ${localStorage.getItem('token')}` }
      });

      if (!response.ok) {
        const data = await response.json();
        throw new Error(data.error || 'Erreur lors de la suppression');
      }

      fetchTeam();
    } catch (err: any) {
      alert(err.message);
    }
  };

  const handleSubmit = async (e: Event) => {
    e.preventDefault();
    setLoading(true);
    setMessage('');
    setError('');

    const url = editingUser ? `/api/users/${editingUser.id}` : '/api/users/invite';
    const method = editingUser ? 'PUT' : 'POST';

    try {
      const response = await fetch(url, {
        method,
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${localStorage.getItem('token')}`
        },
        body: JSON.stringify({ 
          first_name: firstName,
          last_name: lastName,
          phone, 
          role,
          shop_id: selectedShopId || null
        })
      });

      const data = await response.json();

      if (!response.ok) {
        throw new Error(data.error || 'Une erreur est survenue');
      }

      setMessage(editingUser ? 'Utilisateur mis à jour !' : 'Utilisateur invité avec succès !');
      if (!editingUser) setTempPassword(data.temp_password);
      
      cancelEdit();
      fetchTeam();
    } catch (err: any) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="page users">
      <h1>Gestion de l'Équipe</h1>
      
      <div className="layout-grid">
        <div className="card">
          <h3>{editingUser ? 'Modifier le collaborateur' : 'Inviter un collaborateur'}</h3>
          <p style={{color: 'var(--text-light)', marginBottom: '1.5rem', fontSize: '0.875rem'}}>
            {editingUser ? 'Mettez à jour les informations de votre collaborateur.' : 'Envoyez une invitation pour permettre à un membre de votre équipe d\'accéder à la gestion des stocks.'}
          </p>
          
          {message && <div className="alert alert-success">{message}</div>}
          {error && <div className="alert alert-error">{error}</div>}

          <form onSubmit={handleSubmit}>
            <div className="layout-grid" style={{gridTemplateColumns: '1fr 1fr', gap: '1rem', marginTop: 0}}>
              <div className="form-group">
                <label>Prénom</label>
                <input 
                  type="text" 
                  value={firstName}
                  onInput={(e) => setFirstName(e.currentTarget.value)}
                  required
                />
              </div>
              <div className="form-group">
                <label>Nom</label>
                <input 
                  type="text" 
                  value={lastName}
                  onInput={(e) => setLastName(e.currentTarget.value)}
                  required
                />
              </div>
            </div>

            <div className="form-group">
              <label>Numéro de téléphone</label>
              <input 
                type="tel" 
                placeholder="Ex: 771234567"
                value={phone}
                onInput={(e) => setPhone(e.currentTarget.value)}
                required
              />
            </div>
            
            <div className="form-group">
              <label>Rôle</label>
              <select value={role} onChange={(e) => setRole(e.currentTarget.value)}>
                {roles.map(r => (
                  <option key={r.value} value={r.value}>{r.label}</option>
                ))}
              </select>
            </div>

            {role === 'vendor' && (
              <div className="form-group">
                <label>Boutique assignée</label>
                <select value={selectedShopId} onChange={(e) => setSelectedShopId(e.currentTarget.value)} required>
                  <option value="">Sélectionner une boutique</option>
                  {shops.map(s => (
                    <option key={s.id} value={s.id}>{s.name}</option>
                  ))}
                </select>
              </div>
            )}

            <div style={{display: 'flex', gap: '1rem'}}>
              <button type="submit" className="btn btn-primary" style={{flex: 1}} disabled={loading}>
                {loading ? 'Traitement...' : (editingUser ? 'Enregistrer' : 'Envoyer l\'invitation')}
              </button>
              {editingUser && (
                <button type="button" className="btn" style={{background: '#f1f5f9'}} onClick={cancelEdit}>
                  Annuler
                </button>
              )}
            </div>
          </form>

          {tempPassword && (
            <div className="card highlight" style={{marginTop: '2rem', borderLeft: '6px solid var(--success)', background: '#f0fdf4'}}>
              <h4 style={{margin: 0, color: 'var(--success)'}}>Invitation réussie !</h4>
              <p style={{margin: '0.5rem 0', fontSize: '0.875rem'}}>Partagez ce mot de passe temporaire avec l'utilisateur :</p>
              <div style={{background: 'white', padding: '1rem', borderRadius: '0.5rem', fontWeight: '800', fontSize: '1.25rem', textAlign: 'center', border: '2px dashed var(--success)'}}>
                {tempPassword}
              </div>
              <p style={{marginTop: '0.5rem', fontSize: '0.75rem', color: 'var(--text-light)'}}>Il devra obligatoirement le changer lors de sa première connexion.</p>
            </div>
          )}
        </div>

        <div className="card">
          <h3>Rôles et Permissions</h3>
          <div className="role-info" style={{marginTop: '1rem'}}>
            <div style={{marginBottom: '1rem'}}>
              <strong>Administrateur</strong>
              <p style={{fontSize: '0.8125rem', color: 'var(--text-light)'}}>Accès complet à toutes les fonctionnalités et paramètres.</p>
            </div>
            <div style={{marginBottom: '1rem'}}>
              <strong>Gestionnaire</strong>
              <p style={{fontSize: '0.8125rem', color: 'var(--text-light)'}}>Peut gérer les stocks, les articles et les boutiques.</p>
            </div>
            <div>
              <strong>Vendeur</strong>
              <p style={{fontSize: '0.8125rem', color: 'var(--text-light)'}}>Peut enregistrer des mouvements de stock uniquement.</p>
            </div>
          </div>
        </div>
      </div>

      <div className="card" style={{marginTop: '2rem'}}>
        <h3>Membres de l'Équipe</h3>
        <table className="data-table">
          <thead>
            <tr>
              <th>Nom complet</th>
              <th>Téléphone</th>
              <th>Rôle</th>
              <th>Boutique</th>
              <th>Statut</th>
              <th style={{textAlign: 'right'}}>Actions</th>
            </tr>
          </thead>
          <tbody>
            {team.map(user => (
              <tr key={user.id}>
                <td>{user.first_name} {user.last_name}</td>
                <td>{user.phone}</td>
                <td>
                  <span className="badge badge-info">{user.role}</span>
                </td>
                <td>{shops.find(s => s.id === user.shop_id)?.name || '-'}</td>
                <td>
                  {user.is_phone_verified ? 
                    <span className="badge badge-success">Vérifié</span> : 
                    <span className="badge" style={{background: '#f1f5f9', color: '#64748b'}}>En attente</span>
                  }
                </td>
                <td style={{textAlign: 'right'}}>
                  <button className="btn btn-sm" style={{background: 'var(--primary-light)', color: 'var(--primary)', marginRight: '0.5rem'}} onClick={() => handleEdit(user)}>
                    Modifier
                  </button>
                  {user.role !== 'owner' && (
                    <button className="btn btn-sm" style={{background: '#fef2f2', color: 'var(--error)'}} onClick={() => handleDelete(user.id)}>
                      Supprimer
                    </button>
                  )}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
};
