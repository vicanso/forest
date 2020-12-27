import request from "@/helpers/request";
import {
  COMMONS_CAPTCHA,
  COMMONS_ROUTERS,
  COMMONS_STATUSES
} from "@/constants/url";
import { formatDate } from "@/helpers/util";

const prefix = "common";
const mutationCommonProcessing = `${prefix}.processing`;
const mutationCommonRouterList = `${prefix}.router.list`;
const mutationCommonStatusList = `${prefix}.status.list`;

const state = {
  processing: false,
  routers: null,
  statuses: null
};

// listStatus 获取状态列表
export async function listStatus({ commit }) {
  if (state.statuses) {
    return {
      statuses: state.statuses
    };
  }
  commit(mutationCommonProcessing, true);
  try {
    const { data } = await request.get(COMMONS_STATUSES);
    commit(mutationCommonStatusList, data);
    return data;
  } finally {
    commit(mutationCommonProcessing, false);
  }
}

export function attachStatusDesc(item) {
  if (!state.statuses) {
    return;
  }
  state.statuses.forEach(status => {
    if (item.status === status.value) {
      item.statusDesc = status.name;
    }
  });
}
export function attachUpdatedAtDesc(item) {
  item.updatedAtDesc = formatDate(item.updatedAt);
}

export default {
  state,
  mutations: {
    [mutationCommonProcessing](state, processing) {
      state.processing = processing;
    },
    [mutationCommonRouterList](state, { routers }) {
      if (routers) {
        routers.forEach(item => {
          item.key = `${item.method} ${item.route}`;
        });
      }
      state.routers = routers;
    },
    [mutationCommonStatusList](state, { statuses }) {
      state.statuses = statuses;
    }
  },
  actions: {
    // getCaptcha 获取图形验证码
    async getCaptcha() {
      const { data } = await request.get(COMMONS_CAPTCHA);
      return data;
    },
    // listRouter 获取路由列表
    async listRouter({ commit }) {
      if (state.routers) {
        return {
          routers: state.routers
        };
      }
      commit(mutationCommonProcessing, true);
      try {
        const { data } = await request.get(COMMONS_ROUTERS);
        commit(mutationCommonRouterList, data);
        return data;
      } finally {
        commit(mutationCommonProcessing, false);
      }
    },
    // listStatus 获取状态列表
    listStatus
  }
};
