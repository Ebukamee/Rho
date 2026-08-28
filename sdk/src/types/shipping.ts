import type { ID } from "./common";

export type ShippingStatus = "pending" | "in_transit" | "delivered" | "failed";

export interface Shipping {
  id: ID;
  orderId: ID;
  provider?: string;
  trackingNumber?: string;
  status: ShippingStatus;
  address?: Record<string, unknown>;
  createdAt?: string;
  updatedAt?: string;
}

export interface CreateShippingRequest {
  orderId: ID;
  provider?: string;
  trackingNumber?: string;
  address?: Record<string, unknown>;
}
