import type { RhoClient } from "../client";
import type { Category, CreateCategoryRequest, UpdateCategoryRequest } from "../types";
import type { PageResponse } from "../types/common";

export function createCategoriesModule(client: RhoClient) {
  return {
    list: (params?: { page?: number; limit?: number }) =>
      client.get<PageResponse<Category>>("/api/v1/categories", params),
    get: (id: string) => client.get<Category>(`/api/v1/categories/${id}`),
    adminList: (params?: { page?: number; limit?: number }) =>
      client.get<PageResponse<Category>>("/api/v1/categories/admin", params),
    create: (payload: CreateCategoryRequest) => client.post<Category>("/api/v1/categories", payload),
    update: (id: string, payload: UpdateCategoryRequest) =>
      client.put<Category>(`/api/v1/categories/${id}`, payload),
    delete: (id: string) => client.delete<{ success: boolean }>(`/api/v1/categories/${id}`),
  };
}
