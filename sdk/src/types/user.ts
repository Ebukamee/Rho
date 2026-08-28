import type { ID } from "./common";

export interface UserProfileUpdateRequest {
  firstName?: string;
  lastName?: string;
  email?: string;
}

export interface UserRoleUpdateRequest {
  role: "customer" | "admin" | "super_admin";
}

export interface UserListItem extends Record<string, unknown> {
  id: ID;
  email: string;
  firstName: string;
  lastName: string;
  role: string;
}
