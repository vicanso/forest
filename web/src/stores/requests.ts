import { defineStore } from "pinia";

import { FLUXES_TAG_VALUES, FLUXES_REQUESTS } from "../constants/url";
import request from "../helpers/request";
import { formatDate } from "../helpers/util";

export const measurementHttpRequest = "httpRequest";
// 后端HTTP请求记录
interface Request {
  _time: string;
  key: string;
  createdAt: string;
  hostname: string;
  addr: string;
  service: string;
  method: string;
  route: string;
  uri: string;
  status: number;
  reused: boolean;
  dnsUse: number;
  tcpUse: number;
  tlsUse: number;
  processingUse: number;
  use: number;
  result: string;
  errCategory: string;
  error: string;
  exception: boolean;
}

export const useRequestsStore = defineStore("requests", {
  state: () => {
    return {
      requests: [] as Request[],
      count: -1,
      processing: false,
      services: [] as string[],
      fetchingServices: false,
      routes: [] as string[],
      fetchingRoutes: false,
    };
  },
  actions: {
    // 获取请求记录列表
    async list(params: {
      route?: string;
      service?: string;
      errCategory?: string;
      begin: string;
      end: string;
      exception?: string;
      limit: number;
    }): Promise<void> {
      if (this.processing) {
        return;
      }
      try {
        this.processing = true;
        const { data } = await request.get<{
          requests: Request[];
          count: number;
        }>(FLUXES_REQUESTS, {
          params,
        });
        this.count = data.count || 0;
        const items = data.requests || [];
        this.requests = items.map((item) => {
          item.key = item._time;
          item.createdAt = formatDate(item._time);
          return item;
        });
      } finally {
        this.processing = false;
      }
    },
    // 获取请求服务列表
    async listService(): Promise<void> {
      if (this.fetchingServices) {
        return;
      }
      try {
        this.fetchingServices = true;
        const url = FLUXES_TAG_VALUES.replace(
          ":measurement",
          measurementHttpRequest
        ).replace(":tag", "service");
        const { data } = await request.get<{
          values: string[];
        }>(url);
        this.services = (data.values || []).sort();
      } finally {
        this.fetchingServices = false;
      }
    },
    // 获取请求服务路由
    async listRoute(): Promise<void> {
      if (this.fetchingRoutes) {
        return;
      }
      try {
        this.fetchingRoutes = true;
        const url = FLUXES_TAG_VALUES.replace(
          ":measurement",
          measurementHttpRequest
        ).replace(":tag", "route");
        const { data } = await request.get<{
          values: string[];
        }>(url);
        this.routes = (data.values || []).sort();
      } finally {
        this.fetchingRoutes = false;
      }
    },
  },
});
