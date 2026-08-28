import type { ID } from "./common";

export interface CheckoutPreviewRequest {
  cartId?: ID;
  items?: Array<{
    productId: ID;
    quantity: number;
  }>;
}

export interface CheckoutPreviewResponse {
  subtotal: number;
  discount: number;
  total: number;
  currency: string;
}

export interface CheckoutRequest {
  cartId?: ID;
  items?: Array<{
    productId: ID;
    quantity: number;
  }>;
  discountCode?: string;
}
