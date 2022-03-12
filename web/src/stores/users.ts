import { defineStore } from "pinia";

import { USERS, USERS_ID } from "../constants/url";
import request from "../helpers/request";

// 用户账户信息
export interface UserAccount {
  [key: string]: unknown;
  id: number;
  account: string;
  groups: string[];
  roles: string[];
  email: string;
  status: number;
}

export const useUsersStore = defineStore("users", {
  state: () => {
    return {
      processing: false,
      count: -1,
      users: [] as UserAccount[],
    };
  },
  actions: {
    // 获取用户列表
    async list(params: {
      keyword?: string;
      limit: number;
      offset: number;
      role?: string;
      status?: string;
      order?: string;
    }): Promise<void> {
      if (this.processing) {
        return;
      }
      try {
        this.processing = true;
        const { data } = await request.get<{
          count: number;
          users: UserAccount[];
        }>(USERS, {
          params,
        });
        const count = data.count || 0;
        if (count >= 0) {
          this.count = count;
        }
        this.users = data.users;
      } finally {
        this.processing = false;
      }
    },
    // 更新用户信息
    async update(params: {
      id: number;
      data: Record<string, unknown>;
    }): Promise<void> {
      await request.patch(USERS_ID.replace(":id", `${params.id}`), params.data);
    },
  },
});
