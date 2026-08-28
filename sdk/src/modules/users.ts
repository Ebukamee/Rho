import type { RhoClient } from "../client";
import type { User, UserListItem, UserProfileUpdateRequest, UserRoleUpdateRequest } from "../types";
import type { PageResponse } from "../types/common";

export function createUsersModule(client: RhoClient) {
  return {
    list: (params?: { page?: number; limit?: number }) =>
      client.get<PageResponse<UserListItem>>("/api/v1/users", params),
    get: (id: string) => client.get<User>(`/api/v1/users/${id}`),
    update: (id: string, payload: UserProfileUpdateRequest) =>
      client.put<User>(`/api/v1/users/${id}`, payload),
    delete: (id: string) => client.delete<{ success: boolean }>(`/api/v1/users/${id}`),
    updateRole: (id: string, payload: UserRoleUpdateRequest) =>
      client.put<User>(`/api/v1/users/${id}/role`, payload),
  };
}
