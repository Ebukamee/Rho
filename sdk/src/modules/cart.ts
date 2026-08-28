import type { RhoClient } from "../client";
import type { AddCartItemRequest, Cart, UpdateCartItemRequest } from "../types";

export function createCartModule(client: RhoClient) {
  return {
    get: () => client.get<Cart>("/api/v1/cart"),
    addItem: (payload: AddCartItemRequest) => client.post<Cart>("/api/v1/cart/items", payload),
    updateItem: (itemId: string, payload: UpdateCartItemRequest) =>
      client.put<Cart>(`/api/v1/cart/items/${itemId}`, payload),
    removeItem: (itemId: string) => client.delete<{ success: boolean }>(`/api/v1/cart/items/${itemId}`),
    clear: () => client.delete<{ success: boolean }>("/api/v1/cart"),
  };
}
