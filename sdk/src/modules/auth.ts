import type { RhoClient } from "../client";
import type {
  AuthTokens,
  LoginRequest,
  RefreshTokenRequest,
  SignupRequest,
  UpdateProfileRequest,
  User,
} from "../types";

export function createAuthModule(client: RhoClient) {
  return {
    signup: (payload: SignupRequest) => client.post<AuthTokens>("/api/v1/auth/signup", payload),
    login: (payload: LoginRequest) => client.post<AuthTokens>("/api/v1/auth/login", payload),
    refresh: (payload: RefreshTokenRequest) => client.post<AuthTokens>("/api/v1/auth/refresh", payload),
    logout: () => client.post<{ success: boolean }>("/api/v1/auth/logout"),
    getProfile: () => client.get<User>("/api/v1/auth/me"),
    updateProfile: (payload: UpdateProfileRequest) => client.put<User>("/api/v1/auth/profile", payload),
    changePassword: (payload: { currentPassword: string; newPassword: string }) =>
      client.put<{ success: boolean }>("/api/v1/auth/password", payload),
    googleLogin: () => client.get<string>("/api/v1/auth/google/login"),
  };
}
