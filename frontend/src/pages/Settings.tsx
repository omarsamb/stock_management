import { useState, useEffect } from 'preact/hooks';

export const Settings = ({ path }: { path?: string }) => {
  const [primaryColor, setPrimaryColor] = useState('#4f46e5');
  const [bgType, setBgType] = useState<'color' | 'image'>('color');
  const [bgColor, setBgColor] = useState('#f8fafc');
  const [bgImage, setBgImage] = useState('');
  const [loading, setLoading] = useState(false);
  const [message, setMessage] = useState({ type: '', text: '' });

  useEffect(() => {
    const savedAccount = localStorage.getItem('user_account');
    if (savedAccount) {
      const account = JSON.parse(savedAccount);
      setPrimaryColor(account.primary_color || '#4f46e5');
      if (account.background_image && account.background_image.startsWith('http')) {
        setBgType('image');
        setBgImage(account.background_image);
      } else if (account.background_image) {
        setBgType('color');
        setBgColor(account.background_image);
      }
    }
  }, []);

  const handleSave = async (e: Event) => {
    e.preventDefault();
    setLoading(true);
    setMessage({ type: '', text: '' });

    const backgroundImage = bgType === 'image' ? bgImage : bgColor;

    try {
      const response = await fetch('/api/auth/theme', {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${localStorage.getItem('token')}`
        },
        body: JSON.stringify({
          primary_color: primaryColor,
          background_image: backgroundImage
        })
      });

      if (!response.ok) {
        const data = await response.json();
        throw new Error(data.error || 'Erreur lors de la mise à jour');
      }

      // Update local storage
      const savedAccount = JSON.parse(localStorage.getItem('user_account') || '{}');
      savedAccount.primary_color = primaryColor;
      savedAccount.background_image = backgroundImage;
      localStorage.setItem('user_account', JSON.stringify(savedAccount));

      // Apply changes globally
      document.documentElement.style.setProperty('--primary', primaryColor);
      document.documentElement.style.setProperty('--app-bg-color', bgType === 'color' ? bgColor : 'transparent');
      document.documentElement.style.setProperty('--app-bg-image', bgType === 'image' ? `url(${bgImage})` : 'none');

      setMessage({ type: 'success', text: 'Paramètres enregistrés avec succès !' });
    } catch (err: any) {
      setMessage({ type: 'error', text: err.message });
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="page settings">
      <div className="section-header">
        <h1>Paramètres de l'Espace</h1>
        <p className="subtitle">Personnalisez l'apparence de votre logiciel pour qu'il vous ressemble.</p>
      </div>

      <div className="layout-grid" style={{gridTemplateColumns: 'minmax(0, 2fr) 1fr'}}>
        <div className="card">
          {message.text && (
            <div className={`alert alert-${message.type}`} style={{marginBottom: '2rem'}}>
              {message.text}
            </div>
          )}

          <form onSubmit={handleSave}>
            <section style={{marginBottom: '3rem'}}>
              <h3 style={{marginBottom: '1.5rem', display: 'flex', alignItems: 'center', gap: '0.5rem'}}>
                <span style={{width: '8px', height: '24px', background: 'var(--primary)', borderRadius: '4px'}}></span>
                Couleur Identitaire
              </h3>
              <div className="form-group">
                <label>Couleur principale (Boutons, liens, accents)</label>
                <div style={{display: 'flex', gap: '1rem', alignItems: 'center'}}>
                  <input 
                    type="color" 
                    value={primaryColor} 
                    onInput={(e) => setPrimaryColor(e.currentTarget.value)}
                    style={{width: '60px', height: '60px', padding: '4px', borderRadius: '12px', cursor: 'pointer', border: '2px solid var(--border)'}}
                  />
                  <input 
                    type="text" 
                    value={primaryColor} 
                    onInput={(e) => setPrimaryColor(e.currentTarget.value)}
                    style={{fontFamily: 'monospace', fontSize: '1.1rem'}}
                  />
                </div>
              </div>
            </section>

            <section style={{marginBottom: '3rem'}}>
              <h3 style={{marginBottom: '1.5rem', display: 'flex', alignItems: 'center', gap: '0.5rem'}}>
                <span style={{width: '8px', height: '24px', background: 'var(--primary)', borderRadius: '4px'}}></span>
                Arrière-plan
              </h3>
              
              <div className="tab-group" style={{marginBottom: '2rem'}}>
                <button 
                  type="button"
                  className={`tab ${bgType === 'color' ? 'active' : ''}`} 
                  onClick={() => setBgType('color')}
                >
                  Couleur Unie
                </button>
                <button 
                  type="button"
                  className={`tab ${bgType === 'image' ? 'active' : ''}`} 
                  onClick={() => setBgType('image')}
                >
                  Image (URL)
                </button>
              </div>

              {bgType === 'color' ? (
                <div className="form-group">
                  <label>Choisir une couleur de fond</label>
                  <div style={{display: 'flex', gap: '1rem', alignItems: 'center'}}>
                    <input 
                      type="color" 
                      value={bgColor} 
                      onInput={(e) => setBgColor(e.currentTarget.value)}
                      style={{width: '60px', height: '60px', padding: '4px', borderRadius: '12px', cursor: 'pointer', border: '2px solid var(--border)'}}
                    />
                    <input 
                      type="text" 
                      value={bgColor} 
                      onInput={(e) => setBgColor(e.currentTarget.value)}
                      style={{fontFamily: 'monospace', fontSize: '1.1rem'}}
                    />
                  </div>
                </div>
              ) : (
                <div className="form-group">
                  <label>Lien de l'image de fond (URL)</label>
                  <input 
                    type="url" 
                    value={bgImage} 
                    onInput={(e) => setBgImage(e.currentTarget.value)}
                    placeholder="https://images.unsplash.com/..." 
                    style={{padding: '1rem'}}
                  />
                  <p style={{fontSize: '0.8rem', color: 'var(--text-light)', marginTop: '0.5rem'}}>Utilisez une image de haute qualité pour un meilleur rendu.</p>
                </div>
              )}
            </section>

            <button type="submit" className="btn btn-primary btn-block" disabled={loading} style={{padding: '1.25rem'}}>
              {loading ? 'Enregistrement...' : 'Enregistrer les préférences'}
            </button>
          </form>
        </div>

        <div className="settings-preview">
          <div className="card" style={{position: 'sticky', top: '2rem'}}>
            <h3>Aperçu</h3>
            <div style={{
              marginTop: '1.5rem', 
              borderRadius: '1rem', 
              border: '1px solid var(--border)', 
              overflow: 'hidden',
              background: bgType === 'color' ? bgColor : `url(${bgImage})`,
              backgroundSize: 'cover',
              backgroundPosition: 'center',
              height: '300px',
              display: 'flex',
              flexDirection: 'column',
              padding: '1.5rem'
            }}>
              <div style={{background: 'white', padding: '1rem', borderRadius: '0.75rem', boxShadow: 'var(--shadow)'}}>
                <div style={{height: '10px', width: '40%', background: 'var(--border)', borderRadius: '5px', marginBottom: '1rem'}}></div>
                <div style={{height: '40px', width: '100%', background: primaryColor, borderRadius: '8px', display: 'flex', alignItems: 'center', justifyContent: 'center', color: 'white', fontWeight: 600, fontSize: '0.8rem'}}>
                  Bouton de Test
                </div>
              </div>
              <div style={{background: 'white', padding: '1rem', borderRadius: '0.75rem', boxShadow: 'var(--shadow)', marginTop: '1rem', flex: 1}}>
                <div style={{height: '10px', width: '60%', background: 'var(--border)', borderRadius: '5px', marginBottom: '1rem'}}></div>
                <div style={{height: '10px', width: '80%', background: 'var(--border)', borderRadius: '5px', marginBottom: '1rem'}}></div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};
