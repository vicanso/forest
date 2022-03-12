import { defineStore } from "pinia";

import { USERS_LOGINS } from "../constants/url";
import request from "../helpers/request";

// 用户登录信息
interface UserLoginRecord {
  account: string;
  userAgent?: string;
  ip?: string;
  trackID?: string;
  sessionID?: string;
  xForwardedFor?: string;
  country?: string;
  province?: string;
  city?: string;
  isp?: string;
  updatedAt?: string;
  createdAt?: string;
  location?: string;
}

export const useUserLoginsStore = defineStore("userLogins", {
  state: () => {
    return {
      processing: false,
      count: -1,
      logins: [] as UserLoginRecord[],
    };
  },
  actions: {
    // 获取登录记录
    async list(params: {
      account?: string;
      begin: string;
      end: string;
      limit: number;
      offset: number;
      order?: string;
    }): Promise<void> {
      if (this.processing) {
        return;
      }
      try {
        this.processing = true;
        const { data } = await request.get<{
          count: number;
          userLogins: UserLoginRecord[];
        }>(USERS_LOGINS, {
          params,
        });
        const count = data.count || 0;
        if (count >= 0) {
          this.count = count;
        }
        data.userLogins.forEach((item: UserLoginRecord) => {
          const arr: string[] = [];
          if (item.province) {
            arr.push(item.province);
          }
          if (item.city) {
            arr.push(item.city);
          }
          item.location = arr.join("");
        });

        this.logins = data.userLogins;
      } finally {
        this.processing = false;
      }
    },
  },
});
