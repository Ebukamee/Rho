import type { ID } from "./common";

export interface CartItem {
  id: ID;
  productId: ID;
  quantity: number;
  unitPrice: number;
  totalPrice: number;
}

export interface Cart {
  id: ID;
  userId: ID;
  items: CartItem[];
  subtotal: number;
  total: number;
  currency: string;
}

export interface AddCartItemRequest {
  productId: ID;
  quantity: number;
}

export interface UpdateCartItemRequest {
  quantity: number;
}
