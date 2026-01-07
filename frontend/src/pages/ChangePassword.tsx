import { useState } from 'preact/hooks';
import { route } from 'preact-router';

export const ChangePassword = ({ path }: { path?: string }) => {
  const [newPassword, setNewPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e: Event) => {
    e.preventDefault();
    if (newPassword !== confirmPassword) {
      setError('Les mots de passe ne correspondent pas');
      return;
    }

    setLoading(true);
    setError('');

    try {
      const response = await fetch('/api/auth/change-password', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${localStorage.getItem('token')}`
        },
        body: JSON.stringify({ new_password: newPassword })
      });

      if (!response.ok) {
        const data = await response.json();
        throw new Error(data.error || 'Erreur lors du changement de mot de passe');
      }

      // Success
      localStorage.removeItem('must_change_password');
      route('/');
    } catch (err: any) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="auth-page">
      <div className="auth-card">
        <h1>Sécurisez votre compte</h1>
        <p>C'est votre première connexion. Veuillez définir un nouveau mot de passe pour continuer.</p>

        {error && <div className="alert alert-error">{error}</div>}

        <form onSubmit={handleSubmit}>
          <div className="form-group">
            <label>Nouveau mot de passe</label>
            <input 
              type="password" 
              value={newPassword}
              onInput={(e) => setNewPassword(e.currentTarget.value)}
              required
              minLength={6}
            />
          </div>
          <div className="form-group">
            <label>Confirmer le mot de passe</label>
            <input 
              type="password" 
              value={confirmPassword}
              onInput={(e) => setConfirmPassword(e.currentTarget.value)}
              required
            />
          </div>
          <button type="submit" className="btn btn-primary btn-block" disabled={loading}>
            {loading ? 'Mise à jour...' : 'Changer le mot de passe'}
          </button>
        </form>
      </div>
    </div>
  );
};
