import { db, OfflineMovement } from '../db/offlineStorage';

const API_BASE_URL = '/api';

export const recordMovement = async (movement: Omit<OfflineMovement, 'id' | 'timestamp'>) => {
  if (!navigator.onLine) {
    await db.movements.add({
      ...movement,
      timestamp: Date.now()
    });
    return { offline: true, message: 'Enregistré localement (hors-ligne)' };
  }

  try {
    const response = await fetch(`${API_BASE_URL}/stocks/movement`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      },
      body: JSON.stringify(movement)
    });

    if (!response.ok) throw new Error('Erreur réseau');
    return await response.json();
  } catch (error) {
    // Failover to offline storage if request fails (connection lost during request)
    await db.movements.add({
      ...movement,
      timestamp: Date.now()
    });
    return { offline: true, message: 'Echec connexion - sauvegardé pour plus tard' };
  }
};

export const syncOfflineData = async () => {
  if (!navigator.onLine) return;

  const pending = await db.movements.toArray();
  if (pending.length === 0) return;

  console.log(`Synchronisation de ${pending.length} mouvements...`);

  for (const move of pending) {
    try {
      const response = await fetch(`${API_BASE_URL}/stocks/movement`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${localStorage.getItem('token')}`
        },
        body: JSON.stringify({
          shop_id: move.shop_id,
          article_id: move.article_id,
          type: move.type,
          qty: move.qty,
          reason: `[Offline Sync] ${move.reason}`,
          device_id: 'PWA-Offline'
        })
      });

      if (response.ok) {
        await db.movements.delete(move.id!);
      }
    } catch (error) {
      console.error('Failed to sync item:', error);
      break; // Stop and retry later if network fails again
    }
  }
};

// Auto-sync when coming online
window.addEventListener('online', syncOfflineData);
