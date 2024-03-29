import { defineStore } from "pinia";

import request from "../helpers/request";
import { ADMINS_CACHES, ADMINS_CACHES_KEY } from "../constants/url";

export const useAdminStore = defineStore("admin", {
  state: () => {
    return {
      fetchingCacheKeys: false,
      cacheKeys: [] as string[],
    };
  },
  actions: {
    // 根据关键字查询缓存列表
    async listCacheKeys(params: {
      keyword: string;
      limit: number;
      offset: number;
    }): Promise<void> {
      if (this.fetchingCacheKeys) {
        return;
      }
      try {
        this.fetchingCacheKeys = true;
        const { data } = await request.get<{
          keys: string[];
        }>(ADMINS_CACHES, {
          params,
        });
        this.cacheKeys = data.keys || [];
      } finally {
        this.fetchingCacheKeys = false;
      }
    },
    // 删除缓存
    async removeCache(key: string): Promise<void> {
      if (this.fetchingCacheKeys) {
        return;
      }
      try {
        this.fetchingCacheKeys = true;
        const url = ADMINS_CACHES_KEY.replace(":key", key);
        await request.delete(url);
        this.cacheKeys = this.cacheKeys.filter((item) => item !== key);
      } finally {
        this.fetchingCacheKeys = false;
      }
    },
    // 获取缓存
    async getCache(key: string): Promise<string> {
      const url = ADMINS_CACHES_KEY.replace(":key", key);
      const { data } = await request.get<{
        data: string;
      }>(url);
      return data.data;
    },
  },
});
