import { RhoApiError } from "./types/common";
import { createAuthModule } from "./modules/auth";
import { createCartModule } from "./modules/cart";
import { createCategoriesModule } from "./modules/categories";
import { createCheckoutModule } from "./modules/checkout";
import { createDiscountsModule } from "./modules/discounts";
import { createInventoryModule } from "./modules/inventory";
import { createOrdersModule } from "./modules/orders";
import { createPaymentsModule } from "./modules/payments";
import { createProductsModule } from "./modules/products";
import { createShippingModule } from "./modules/shipping";
import { createUsersModule } from "./modules/users";

export type HttpMethod = "GET" | "POST" | "PUT" | "PATCH" | "DELETE";

export interface RhoClientOptions {
  baseURL: string;
  token?: string;
  fetch?: typeof fetch;
  headers?: Record<string, string>;
}

export class RhoClient {
  baseURL: string;
  token?: string;
  fetchFn: typeof fetch;
  headers: Record<string, string>;

  auth = createAuthModule(this);
  products = createProductsModule(this);
  categories = createCategoriesModule(this);
  cart = createCartModule(this);
  orders = createOrdersModule(this);
  payments = createPaymentsModule(this);
  shipping = createShippingModule(this);
  discounts = createDiscountsModule(this);
  checkout = createCheckoutModule(this);
  inventory = createInventoryModule(this);
  users = createUsersModule(this);

  constructor(options: RhoClientOptions) {
    this.baseURL = options.baseURL.replace(/\/+$/, "");
    this.token = options.token;
    this.fetchFn = options.fetch ?? fetch;
    this.headers = {
      "Content-Type": "application/json",
      ...(options.headers ?? {}),
    };
  }

  setToken(token: string) {
    this.token = token;
  }

  clearToken() {
    this.token = undefined;
  }

  buildUrl(path: string, params?: Record<string, string | number | boolean | undefined>) {
    const url = new URL(path.startsWith("http") ? path : `${this.baseURL}${path}`);

    if (params) {
      Object.entries(params).forEach(([key, value]) => {
        if (value !== undefined && value !== null) {
          url.searchParams.set(key, String(value));
        }
      });
    }

    return url.toString();
  }

  private async request<T>(
    method: HttpMethod,
    path: string,
    body?: unknown,
    params?: Record<string, string | number | boolean | undefined>,
    headers?: Record<string, string>
  ): Promise<T> {
    const url = this.buildUrl(path, params);

    const response = await this.fetchFn(url, {
      method,
      headers: {
        ...this.headers,
        ...headers,
        ...(this.token ? { Authorization: `Bearer ${this.token}` } : {}),
      },
      body: body !== undefined ? JSON.stringify(body) : undefined,
    });

    const contentType = response.headers.get("content-type") ?? "";
    const isJson = contentType.includes("application/json");
    const payload = isJson ? await response.json() : await response.text();

    if (!response.ok) {
      throw new RhoApiError(response.status, payload as any);
    }

    return payload as T;
  }

  get<T>(path: string, params?: Record<string, string | number | boolean | undefined>) {
    return this.request<T>("GET", path, undefined, params);
  }

  post<T>(path: string, body?: unknown, params?: Record<string, string | number | boolean | undefined>) {
    return this.request<T>("POST", path, body, params);
  }

  put<T>(path: string, body?: unknown, params?: Record<string, string | number | boolean | undefined>) {
    return this.request<T>("PUT", path, body, params);
  }

  patch<T>(path: string, body?: unknown, params?: Record<string, string | number | boolean | undefined>) {
    return this.request<T>("PATCH", path, body, params);
  }

  delete<T>(path: string, params?: Record<string, string | number | boolean | undefined>) {
    return this.request<T>("DELETE", path, undefined, params);
  }
}

export function createClient(options: RhoClientOptions) {
  return new RhoClient(options);
}
