import type { RhoClient } from "../client";
import type { AdjustInventoryRequest, CreateInventoryRequest, InventoryItem } from "../types";

export function createInventoryModule(client: RhoClient) {
  return {
    create: (payload: CreateInventoryRequest) => client.post<InventoryItem>("/api/v1/inventory", payload),
    get: (id: string) => client.get<InventoryItem>(`/api/v1/inventory/${id}`),
    getByProduct: (productId: string) => client.get<InventoryItem>(`/api/v1/inventory/product/${productId}`),
    update: (id: string, payload: Partial<CreateInventoryRequest>) =>
      client.put<InventoryItem>(`/api/v1/inventory/${id}`, payload),
    adjust: (productId: string, payload: AdjustInventoryRequest) =>
      client.post<InventoryItem>(`/api/v1/inventory/product/${productId}/adjust`, payload),
    delete: (id: string) => client.delete<{ success: boolean }>(`/api/v1/inventory/${id}`),
  };
}
