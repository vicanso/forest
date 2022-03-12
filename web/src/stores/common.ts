import { defineStore } from "pinia";
import {
  COMMONS_CAPTCHA,
  COMMONS_ROUTERS,
  COMMONS_HTTP_STATS,
} from "../constants/url";
import request from "../helpers/request";
import { settingStorage } from "../storages/local";

export interface Captcha {
  data: string;
  expiredAt: string;
  id: string;
  type: string;
}

// 路由
interface Router {
  method: string;
  route: string;
}

// http实例
interface HTTPInstance {
  name: string;
  maxConcurrency: number;
  concurrency: number;
}

export enum Mode {
  Add = "add",
  Update = "update",
  List = "list",
}

export const useCommonStore = defineStore("commonStore", {
  state: () => {
    return {
      setting: {
        theme: "",
        collapsed: false,
      },
      routers: [] as Router[],
      fetchingRouters: false,
      httpInstances: [] as HTTPInstance[],
      fetchingHTTPInstances: false,
    };
  },
  actions: {
    // 获取图形验证码
    async getCaptcha(): Promise<Captcha> {
      const { data } = await request.get<Captcha>(COMMONS_CAPTCHA);
      return data;
    },
    // 获取配置
    getSetting(): void {
      const data = settingStorage.getData();
      if (data.theme) {
        this.setting.theme = data.theme as string;
      } else if (
        window.matchMedia &&
        window.matchMedia("(prefers-color-scheme: light)").matches
      ) {
        this.setting.theme = "light";
      }
      this.setting.collapsed = data.collapsed as boolean;
    },
    // 设置主题
    async setTheme(theme: string): Promise<void> {
      await settingStorage.set("theme", theme);
      this.setting.theme = theme;
    },
    async setCollapsed(collapsed: boolean): Promise<void> {
      await settingStorage.set("collapsed", collapsed);
      this.setting.collapsed = collapsed;
    },
    // 获取路由列表
    async listRouter(): Promise<void> {
      if (this.fetchingRouters) {
        return;
      }
      try {
        this.fetchingRouters = true;
        const { data } = await request.get<{
          routers: Router[];
        }>(COMMONS_ROUTERS);
        this.routers = data.routers || [];
      } finally {
        this.fetchingRouters = false;
      }
    },
    // 获取http实例配置
    async listHTTPInstance(): Promise<void> {
      if (this.fetchingHTTPInstances) {
        return;
      }
      try {
        this.fetchingHTTPInstances = true;
        const { data } = await request.get<{
          statsList: HTTPInstance[];
        }>(COMMONS_HTTP_STATS);
        this.httpInstances = data.statsList || [];
      } finally {
        this.fetchingHTTPInstances = false;
      }
    },
  },
});
