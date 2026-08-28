import type { RhoClient } from "../client";
import type { CreateOrderRequest, Order, UpdateOrderStatusRequest } from "../types";

export function createOrdersModule(client: RhoClient) {
  return {
    create: (payload: CreateOrderRequest) => client.post<Order>("/api/v1/orders", payload),
    get: (id: string) => client.get<Order>(`/api/v1/orders/${id}`),
    updateStatus: (id: string, payload: UpdateOrderStatusRequest) =>
      client.put<Order>(`/api/v1/orders/${id}/status`, payload),
  };
}
