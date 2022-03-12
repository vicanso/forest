import { defineStore } from "pinia";

import request from "../helpers/request";
import { FLUXES_HTTP_ERRORS, FLUXES_TAG_VALUES } from "../constants/url";
import { formatDate } from "../helpers/util";
export const measurementHttpError = "httpError";

// HTTPError 客户端HTTP请求出错记录
interface HTTPError {
  _time: string;
  createdAt: string;
  key: string;
  account: string;
  category: string;
  error: string;
  exception: boolean;
  hostname: string;
  ip: string;
  method: string;
  route: string;
  sid: string;
  status: number;
  tid: string;
  uri: string;
}

export const useHTTPErrorsStore = defineStore("httpErrors", {
  state: () => {
    return {
      httpErrors: [] as HTTPError[],
      count: -1,
      processing: false,
      categories: [] as string[],
      fetchingCategories: false,
    };
  },
  actions: {
    // 获取出错记录列表
    async list(params: {
      account?: string;
      category?: string;
      begin: string;
      end: string;
      exception?: string;
      limit: number;
    }): Promise<void> {
      if (this.processing) {
        return;
      }
      try {
        this.processing = false;
        const { data } = await request.get<{
          httpErrors: HTTPError[];
          count: number;
        }>(FLUXES_HTTP_ERRORS, {
          params,
        });
        const items = data.httpErrors || [];
        this.count = data.count || 0;
        this.httpErrors = items.map((item) => {
          item.key = item._time;
          item.createdAt = formatDate(item._time);
          return item;
        });
      } finally {
        this.processing = false;
      }
    },
    // 获取出错类型
    async listCategory(): Promise<void> {
      if (this.fetchingCategories) {
        return;
      }
      try {
        this.fetchingCategories = true;
        const url = FLUXES_TAG_VALUES.replace(
          ":measurement",
          measurementHttpError
        ).replace(":tag", "category");
        const { data } = await request.get<{
          values: string[];
        }>(url);
        this.categories = (data.values || []).sort();
      } finally {
        this.fetchingCategories = false;
      }
    },
  },
});
