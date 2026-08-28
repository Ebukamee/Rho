import type { ID } from "./common";

export interface Payment {
  id: ID;
  orderId: ID;
  provider: string;
  status: "pending" | "success" | "failed";
  amount: number;
  currency: string;
  reference?: string;
  createdAt?: string;
}

export interface InitializePaymentRequest {
  orderId: ID;
  email: string;
  amount: number;
  currency?: string;
}

export interface PaymentInitializationResponse {
  authorizationUrl?: string;
  accessCode?: string;
  reference?: string;
  paymentId?: ID;
}
