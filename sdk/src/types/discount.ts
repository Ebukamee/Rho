import type { ID } from "./common";

export interface Discount {
  id: ID;
  code: string;
  type: "percentage" | "fixed";
  value: number;
  active: boolean;
  createdAt?: string;
  updatedAt?: string;
}

export interface CreateDiscountRequest {
  code: string;
  type: "percentage" | "fixed";
  value: number;
  active?: boolean;
}

export interface ApplyDiscountRequest {
  code: string;
  orderTotal?: number;
}
