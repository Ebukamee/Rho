import type { RhoClient } from "../client";
import type { CreateShippingRequest, Shipping } from "../types";

export function createShippingModule(client: RhoClient) {
  return {
    create: (payload: CreateShippingRequest) => client.post<Shipping>("/api/v1/shipping", payload),
    get: (id: string) => client.get<Shipping>(`/api/v1/shipping/${id}`),
    getByOrder: (orderId: string) => client.get<Shipping>(`/api/v1/shipping/order/${orderId}`),
    update: (id: string, payload: Partial<CreateShippingRequest>) =>
      client.put<Shipping>(`/api/v1/shipping/${id}`, payload),
    delete: (id: string) => client.delete<{ success: boolean }>(`/api/v1/shipping/${id}`),
  };
}
