import { defineStore } from "pinia";

import request from "../helpers/request";
import { FLUXES_TRACKERS, FLUXES_TAG_VALUES } from "../constants/url";
import { formatDate } from "../helpers/util";

// 用户行为轨迹
interface UserTracker {
  [key: string]: unknown;
  _time: string;
  key: string;
  createdAt: string;
  account: string;
  action: string;
  hostname: string;
  ip: string;
  rslt: string;
  sid: string;
  tid: string;
  form: string;
  query: string;
  params: string;
  error: string;
}

export const measurementUserTracker = "userTracker";

export const useUserTrackersStore = defineStore("userTrackers", {
  state: () => {
    return {
      count: -1,
      trackers: [] as UserTracker[],
      processing: false,
      actions: [] as string[],
      fetchingActions: false,
    };
  },
  actions: {
    // 获取用户行为列表
    async list(params: {
      account?: string;
      action?: string;
      begin: string;
      end: string;
      limit: number;
      result?: string;
    }): Promise<void> {
      if (this.processing) {
        return;
      }
      try {
        this.processing = true;

        const { data } = await request.get<{
          trackers: UserTracker[];
          count: number;
        }>(FLUXES_TRACKERS, {
          params,
        });
        const items = data.trackers || [];
        this.count = data.count || 0;
        // 倒序
        this.trackers = items.reverse().map((item) => {
          if (item.error) {
            const reg = /, message=([\s\S]*)/;
            const result = reg.exec(item.error);
            if (result && result.length === 2) {
              item.error = `${result[1]}, ${item.error.replace(result[0], "")}`;
            }
          }
          item.key = item._time;
          item.createdAt = formatDate(item._time);
          return item;
        });
      } finally {
        this.processing = false;
      }
    },
    // 获取用户行为类型
    async listActions(): Promise<void> {
      if (this.fetchingActions) {
        return;
      }
      try {
        this.fetchingActions = true;
        const url = FLUXES_TAG_VALUES.replace(
          ":measurement",
          measurementUserTracker
        ).replace(":tag", "action");
        const { data } = await request.get<{
          values: string[];
        }>(url);
        this.actions = (data.values || []).sort();
      } finally {
        this.fetchingActions = false;
      }
    },
  },
});
