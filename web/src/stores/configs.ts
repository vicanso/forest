import { defineStore } from "pinia";

import request from "../helpers/request";

import { CONFIGS, CONFIGS_ID, CONFIGS_CURRENT_VALID } from "../constants/url";

export enum ConfigCategory {
  MockTime = "mockTime",
  BlockIP = "blockIP",
  SignedKey = "signedKey",
  RouterConcurrency = "routerConcurrency",
  SessionInterceptor = "sessionInterceptor",
  RequestConcurrency = "requestConcurrency",
  Router = "router",
  Email = "email",
  HTTPServerInterceptor = "httpServerInterceptor",
}

export enum ConfigStatus {
  Enabled = 1,
  Disabled,
}

// 配置信息
export interface Config {
  [key: string]: unknown;
  key: string;
  id: number;
  createdAt: string;
  updatedAt: string;
  status: number;
  name: string;
  category: string;
  owner: string;
  data: string;
  startedAt: string;
  endedAt: string;
  description?: string;
}

export const useConfigsStore = defineStore("configs", {
  state: () => {
    return {
      count: -1,
      configs: [] as Config[],
      processing: false,
    };
  },
  actions: {
    // 获取mock time的配置
    async getMockTime(): Promise<Config> {
      const { data } = await request.get<{
        configurations: Config[];
      }>(CONFIGS, {
        params: {
          category: ConfigCategory.MockTime,
          name: ConfigCategory.MockTime,
          limit: 1,
        },
      });
      const items = data.configurations || [];
      if (items.length === 0) {
        return <Config>{};
      }
      return items[0];
    },
    // 新增配置
    async add(params: {
      name: string;
      status: number;
      category: string;
      startedAt: string;
      endedAt: string;
      data: string;
    }): Promise<Config> {
      const { data } = await request.post<Config>(CONFIGS, params);
      return data;
    },
    // 通过ID查询配置
    async findByID(id: number): Promise<Config> {
      const url = CONFIGS_ID.replace(":id", `${id}`);
      const { data } = await request.get<Config>(url);
      return data;
    },
    // 通过ID更新配置
    async updateByID(params: {
      id: number;
      data: Record<string, unknown>;
    }): Promise<void> {
      const url = CONFIGS_ID.replace(":id", `${params.id}`);
      await request.patch(url, params.data);
    },
    // 获取配置列表
    async list(params: {
      name?: string;
      category?: string;
      limit?: number;
      offset?: number;
    }): Promise<void> {
      if (this.processing) {
        return;
      }
      if (!params.limit) {
        params.limit = 50;
      }
      try {
        this.processing = true;
        const { data } = await request.get<{
          count: number;
          configurations: Config[];
        }>(CONFIGS, {
          params,
        });
        const count = data.count || 0;
        if (count >= 0) {
          this.count = count;
        }
        this.configs = (data.configurations || []).map((item) => {
          item.key = `${item.id}`;
          return item;
        });
      } finally {
        this.processing = false;
      }
    },
    // 获取当前有效配置
    async getCurrentValid(): Promise<Record<string, unknown>> {
      const { data } = await request.get<Record<string, unknown>>(
        CONFIGS_CURRENT_VALID
      );
      return data;
    },
  },
});
