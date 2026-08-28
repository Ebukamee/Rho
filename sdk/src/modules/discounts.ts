import type { RhoClient } from "../client";
import type { ApplyDiscountRequest, CreateDiscountRequest, Discount } from "../types";

export function createDiscountsModule(client: RhoClient) {
  return {
    apply: (payload: ApplyDiscountRequest) =>
      client.post<{ valid: boolean; discount?: Discount; total?: number }>("/api/v1/discounts/apply", payload),
    create: (payload: CreateDiscountRequest) => client.post<Discount>("/api/v1/discounts", payload),
    get: (id: string) => client.get<Discount>(`/api/v1/discounts/${id}`),
    update: (id: string, payload: Partial<CreateDiscountRequest>) =>
      client.put<Discount>(`/api/v1/discounts/${id}`, payload),
    delete: (id: string) => client.delete<{ success: boolean }>(`/api/v1/discounts/${id}`),
  };
}
