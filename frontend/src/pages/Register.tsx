import { useState } from 'preact/hooks';
import { route } from 'preact-router';

export const Register = ({ path }: { path?: string }) => {
  const [phone, setPhone] = useState('');
  const [password, setPassword] = useState('');
  const [companyName, setCompanyName] = useState('');
  const [error, setError] = useState('');
  const [step, setStep] = useState(1); // 1: Info, 2: Verification
  const [code, setCode] = useState('');

  const handleSubmit = async (e: Event) => {
    e.preventDefault();
    setError('');

    try {
      const response = await fetch('/api/auth/register', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ phone, password, company_name: companyName })
      });

      const data = await response.json();
      if (!response.ok) throw new Error(data.error || 'Erreur lors de l\'inscription');

      setStep(2);
    } catch (err: any) {
      setError(err.message);
    }
  };

  const handleVerify = async (e: Event) => {
    e.preventDefault();
    setError('');

    try {
      const response = await fetch('/api/auth/verify', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ phone, code })
      });

      const data = await response.json();
      if (!response.ok) throw new Error(data.error || 'Code incorrect');

      localStorage.setItem('token', data.token);
      window.dispatchEvent(new Event('auth-change'));
      route('/');
    } catch (err: any) {
      setError(err.message);
    }
  };

  return (
    <div className="auth-page">
      <div className="auth-card">
        <h1>{step === 1 ? 'Créer un compte' : 'Vérification'}</h1>
        <p>{step === 1 ? 'Commencez votre essai gratuit de 14 jours' : 'Entrez le code envoyé sur votre WhatsApp'}</p>
        
        {error && <div className="alert alert-error">{error}</div>}
        
        {step === 1 ? (
          <form onSubmit={handleSubmit}>
            <div className="form-group">
              <label>Nom de l'entreprise</label>
              <input 
                type="text" 
                value={companyName} 
                onInput={(e) => setCompanyName(e.currentTarget.value)} 
                required 
              />
            </div>
            <div className="form-group">
              <label>Numéro de téléphone (WhatsApp)</label>
              <input 
                type="tel" 
                value={phone} 
                onInput={(e) => setPhone(e.currentTarget.value)} 
                required 
                placeholder="Ex: +22177..."
              />
            </div>
            <div className="form-group">
              <label>Mot de passe</label>
              <input 
                type="password" 
                value={password} 
                onInput={(e) => setPassword(e.currentTarget.value)} 
                required 
                minLength={6}
              />
            </div>
            <button type="submit" className="btn btn-primary btn-block">Suivant</button>
          </form>
        ) : (
          <form onSubmit={handleVerify}>
            <div className="form-group">
              <label>Code de vérification (6 chiffres)</label>
              <input 
                type="text" 
                value={code} 
                onInput={(e) => setCode(e.currentTarget.value)} 
                required 
                maxLength={6}
                style={{ textAlign: 'center', fontSize: '1.5rem', letterSpacing: '0.5rem' }}
              />
            </div>
            <button type="submit" className="btn btn-primary btn-block">Vérifier</button>
            <button type="button" onClick={() => setStep(1)} className="btn btn-link">Retour</button>
          </form>
        )}
        
        <p className="auth-footer">
          Déjà un compte ? <a href="/login">Connectez-vous</a>
        </p>
      </div>
    </div>
  );
};
