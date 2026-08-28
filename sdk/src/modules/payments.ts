import type { RhoClient } from "../client";
import type { InitializePaymentRequest, Payment, PaymentInitializationResponse } from "../types";

export function createPaymentsModule(client: RhoClient) {
  return {
    initialize: (payload: InitializePaymentRequest) =>
      client.post<PaymentInitializationResponse>("/api/v1/payments/initialize", payload),
    get: (id: string) => client.get<Payment>(`/api/v1/payments/${id}`),
    verify: (id: string, payload?: { reference?: string }) =>
      client.post<Payment>(`/api/v1/payments/${id}/verify`, payload ?? {}),
  };
}
