import type { ID } from "./common";

export interface InventoryItem {
  id: ID;
  productId: ID;
  sku?: string;
  quantity: number;
  reserved: number;
  available: number;
  location?: string;
  createdAt?: string;
  updatedAt?: string;
}

export interface CreateInventoryRequest {
  productId: ID;
  quantity: number;
  reserved?: number;
  location?: string;
}

export interface AdjustInventoryRequest {
  delta: number;
  reason?: string;
}
