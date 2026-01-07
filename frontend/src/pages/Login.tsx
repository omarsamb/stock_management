import { useState } from 'preact/hooks';
import { route } from 'preact-router';

export const Login = ({ path }: { path?: string }) => {
  const [phone, setPhone] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');

  const handleSubmit = async (e: Event) => {
    e.preventDefault();
    setError('');

    try {
      const response = await fetch('/api/auth/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ phone, password })
      });

      const data = await response.json();
      if (!response.ok) {
        if (data.requires_verification) {
          setError('Veuillez vérifier votre numéro de téléphone.');
          // Logic to redirect to verification could go here
          return;
        }
        throw new Error(data.error || 'Identifiants incorrects');
      }

      localStorage.setItem('token', data.token);
      localStorage.setItem('user_phone', data.user.phone);
      localStorage.setItem('user_role', data.user.role);
      if (data.user.shop_id) localStorage.setItem('user_shop_id', data.user.shop_id);
      if (data.account) localStorage.setItem('user_account', JSON.stringify(data.account));
      
      window.dispatchEvent(new Event('auth-change'));

      if (data.must_change_password) {
        localStorage.setItem('must_change_password', 'true');
        route('/change-password');
      } else {
        route('/');
      }
    } catch (err: any) {
      setError(err.message);
    }
  };

  return (
    <div className="auth-page">
      <div className="auth-card">
        <h1>Connexion</h1>
        <p>Gérez votre stock en toute simplicité</p>
        
        {error && <div className="alert alert-error">{error}</div>}
        
        <form onSubmit={handleSubmit}>
          <div className="form-group">
            <label>Téléphone (WhatsApp)</label>
            <input 
              type="tel" 
              value={phone} 
              onInput={(e) => setPhone(e.currentTarget.value)} 
              required 
              placeholder="+221..."
            />
          </div>
          <div className="form-group">
            <label>Mot de passe</label>
            <input 
              type="password" 
              value={password} 
              onInput={(e) => setPassword(e.currentTarget.value)} 
              required 
            />
          </div>
          <button type="submit" className="btn btn-primary btn-block">Se connecter</button>
        </form>
        
        <p className="auth-footer">
          Pas de compte ? <a href="/register">Inscrivez-vous</a>
        </p>
      </div>
    </div>
  );
};
