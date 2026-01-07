import { useState, useEffect } from 'preact/hooks';

interface Plan {
  id: string;
  name: string;
  price: number;
  features: string[];
}

export const Subscription = ({ path }: { path?: string }) => {
  const [loading, setLoading] = useState(false);
  const [message, setMessage] = useState('');
  const [error, setError] = useState('');

  const plans = [
    {
      id: 'basic',
      name: 'Basique',
      price: 5000,
      features: ['Jusqu\'à 2 boutiques', 'Inventaire illimité', 'Support email']
    },
    {
      id: 'premium',
      name: 'Premium',
      price: 15000,
      features: ['Boutiques illimitées', 'Gestion d\'équipe', 'Rapports avancés', 'Support prioritaire']
    }
  ];

  const handleSelectPlan = async (planId: string) => {
    setLoading(true);
    setMessage('');
    setError('');

    try {
      const response = await fetch('/api/subscription/select', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${localStorage.getItem('token')}`
        },
        body: JSON.stringify({ plan: planId })
      });

      const data = await response.json();

      if (!response.ok) {
        throw new Error(data.error || 'Erreur lors de la sélection du plan');
      }

      setMessage('Plan sélectionné avec succès ! Redirection vers le paiement...');
      // Normally we would redirect to a payment page here
    } catch (err: any) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="page subscription">
      <h1>Gestion de l'Abonnement</h1>
      
      {message && <div className="alert alert-success">{message}</div>}
      {error && <div className="alert alert-error">{error}</div>}

      <div className="layout-grid">
        {plans.map(plan => (
          <div className="card" key={plan.id}>
            <h2 style={{textAlign: 'center', marginBottom: '0.5rem'}}>{plan.name}</h2>
            <p style={{textAlign: 'center', fontSize: '2rem', fontWeight: '800', color: 'var(--primary)', marginBottom: '1.5rem'}}>
              {plan.price.toLocaleString()} CFA<span style={{fontSize: '1rem', color: 'var(--text-light)'}}>/mois</span>
            </p>
            
            <ul style={{listStyle: 'none', padding: 0, marginBottom: '2rem'}}>
              {plan.features.map(feature => (
                <li key={feature} style={{padding: '0.5rem 0', borderBottom: '1px solid var(--border)', fontSize: '0.9375rem'}}>
                  ✓ {feature}
                </li>
              ))}
            </ul>

            <button 
              onClick={() => handleSelectPlan(plan.id)}
              className="btn btn-primary btn-block"
              disabled={loading}
            >
              Choisir ce plan
            </button>
          </div>
        ))}
      </div>
    </div>
  );
};
