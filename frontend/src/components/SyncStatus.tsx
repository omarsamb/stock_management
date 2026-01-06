import { h } from 'preact';
import { useState, useEffect } from 'preact/hooks';
import { db } from '../db/offlineStorage';
import { syncOfflineData } from '../services/api';

export const SyncStatus = () => {
  const [isOnline, setIsOnline] = useState(navigator.onLine);
  const [pendingCount, setPendingCount] = useState(0);

  useEffect(() => {
    const handleOnline = () => {
      setIsOnline(true);
      syncOfflineData();
    };
    const handleOffline = () => setIsOnline(false);

    window.addEventListener('online', handleOnline);
    window.addEventListener('offline', handleOffline);

    // Initial count
    const updateCount = async () => {
      const count = await db.movements.count();
      setPendingCount(count);
    };
    
    updateCount();
    const interval = setInterval(updateCount, 5000);

    return () => {
      window.removeEventListener('online', handleOnline);
      window.removeEventListener('offline', handleOffline);
      clearInterval(interval);
    };
  }, []);

  return (
    <div className={`sync-status ${isOnline ? 'online' : 'offline'}`}>
      <span className="dot"></span>
      {isOnline ? 'Connecté' : 'Hors-ligne'}
      {pendingCount > 0 && (
        <span className="pending-badge">
          ({pendingCount} en attente)
        </span>
      )}
    </div>
  );
};
