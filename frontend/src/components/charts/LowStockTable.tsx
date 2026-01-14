import { h } from 'preact';

export const LowStockTable = ({ items }: { items: any[] }) => {
  if (!items || items.length === 0) return <p>Aucune alerte de stock.</p>;

  return (
    <div style={{ overflowX: 'auto' }}>
      <table style={{ width: '100%', borderCollapse: 'collapse', fontSize: '0.9rem' }}>
        <thead>
          <tr style={{ background: '#f8fafc', color: 'var(--text-light)' }}>
            <th style={{ padding: '0.5rem', textAlign: 'left', borderBottom: '2px solid #e2e8f0' }}>Article</th>
            <th style={{ padding: '0.5rem', textAlign: 'right', borderBottom: '2px solid #e2e8f0' }}>Qté</th>
            <th style={{ padding: '0.5rem', textAlign: 'right', borderBottom: '2px solid #e2e8f0' }}>Seuil</th>
            <th style={{ padding: '0.5rem', textAlign: 'left', borderBottom: '2px solid #e2e8f0' }}>Boutique</th>
          </tr>
        </thead>
        <tbody>
          {items.map((item, i) => (
            <tr key={i} style={{ borderBottom: '1px solid #e2e8f0' }}>
              <td style={{ padding: '0.75rem 0.5rem' }}>{item.article_name}</td>
              <td style={{ padding: '0.75rem 0.5rem', textAlign: 'right', color: 'var(--error)', fontWeight: 'bold' }}>{item.quantity}</td>
              <td style={{ padding: '0.75rem 0.5rem', textAlign: 'right' }}>{item.min_threshold}</td>
              <td style={{ padding: '0.75rem 0.5rem' }}>{item.shop_name}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
};
