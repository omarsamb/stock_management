import Dexie, { Table } from 'dexie';

export interface OfflineMovement {
  id?: number;
  shop_id: string;
  article_id: string;
  type: 'in' | 'out' | 'adjust';
  qty: number;
  reason: string;
  timestamp: number;
}

export class MyDatabase extends Dexie {
  movements!: Table<OfflineMovement>;

  constructor() {
    super('StockManagerDB');
    this.version(1).stores({
      movements: '++id, shop_id, article_id, type, timestamp'
    });
  }
}

export const db = new MyDatabase();
