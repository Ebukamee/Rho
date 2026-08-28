import type { ID } from "./common";

export type OrderStatus = "pending" | "paid" | "processing" | "shipped" | "completed" | "cancelled";

export interface OrderItem {
  id: ID;
  orderId: ID;
  productId: ID;
  name: string;
  sku: string;
  quantity: number;
  unitPrice: number;
  totalPrice: number;
}

export interface Order {
  id: ID;
  userId: ID;
  status: OrderStatus;
  subtotal: number;
  discount: number;
  total: number;
  currency: string;
  createdAt?: string;
  updatedAt?: string;
  items?: OrderItem[];
}

export interface CreateOrderRequest {
  items: Array<{
    productId: ID;
    quantity: number;
  }>;
  shippingAddress?: Record<string, unknown>;
  discountCode?: string;
}

export interface UpdateOrderStatusRequest {
  status: OrderStatus;
}
