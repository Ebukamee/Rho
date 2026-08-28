import type { RhoClient } from "../client";
import type { CheckoutPreviewRequest, CheckoutPreviewResponse, CheckoutRequest } from "../types";

export function createCheckoutModule(client: RhoClient) {
  return {
    preview: (payload: CheckoutPreviewRequest) =>
      client.post<CheckoutPreviewResponse>("/api/v1/checkout/preview", payload),
    create: (payload: CheckoutRequest) => client.post<{ success: boolean; orderId?: string }>("/api/v1/checkout", payload),
  };
}
