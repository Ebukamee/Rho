import type { RhoClient } from "../client";
import type { CreateProductRequest, Product, UpdateProductRequest } from "../types";
import type { PageResponse } from "../types/common";

export function createProductsModule(client: RhoClient) {
  return {
    list: (params?: { page?: number; limit?: number; search?: string }) =>
      client.get<PageResponse<Product>>("/api/v1/products", params),
    get: (id: string) => client.get<Product>(`/api/v1/products/${id}`),
    adminList: (params?: { page?: number; limit?: number }) =>
      client.get<PageResponse<Product>>("/api/v1/products/admin", params),
    create: (payload: CreateProductRequest) => client.post<Product>("/api/v1/products", payload),
    update: (id: string, payload: UpdateProductRequest) =>
      client.put<Product>(`/api/v1/products/${id}`, payload),
    delete: (id: string) => client.delete<{ success: boolean }>(`/api/v1/products/${id}`),
  };
}
