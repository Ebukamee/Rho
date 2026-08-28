import type { ID } from "./common";

export interface Product {
  id: ID;
  categoryId?: ID | null;
  name: string;
  slug: string;
  description: string;
  sku: string;
  price: number;
  currency: string;
  imageUrl?: string;
  active: boolean;
  createdAt?: string;
  updatedAt?: string;
}

export interface CreateProductRequest {
  categoryId?: ID | null;
  name: string;
  slug?: string;
  description?: string;
  sku?: string;
  price: number;
  currency?: string;
  imageUrl?: string;
  active?: boolean;
}

export interface UpdateProductRequest extends Partial<CreateProductRequest> {}
