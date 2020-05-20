import request from "@/helpers/request";
import { CONFIGS, CONFIGS_ID } from "@/constants/url";
import { formatDate } from "@/helpers/util";
import { CONFIG_ENABLED, CONFIG_DISABLED } from "@/constants/config";

const mutationConfigProcessing = "config.processing";
const mutationConfigList = "config.list";
const mutationConfigListReset = "config.list.reset";
const mutationConfigListDelete = "config.list.delete";

const statusList = [
  {
    label: "启用",
    value: CONFIG_ENABLED
  },
  {
    label: "禁用",
    value: CONFIG_DISABLED
  }
];
const state = {
  status: statusList,
  processing: false,
  items: null
};

function formatJSON(data) {
  const item = JSON.parse(data);
  if (item.response && item.response[0] === "{") {
    item.response = JSON.parse(item.response);
  }
  return JSON.stringify(item, null, 2);
}

export default {
  state,
  mutations: {
    // 设置状态为处理中
    [mutationConfigProcessing](state, processing) {
      state.processing = processing;
    },
    // 重置列表数据
    [mutationConfigListReset](state) {
      state.count = -1;
      state.items = null;
    },
    // 设置列表数据
    [mutationConfigList](state, { configs }) {
      state.items = configs;
    },
    // 删除该id对应数据
    [mutationConfigListDelete](state, id) {
      if (!state.items) {
        return;
      }
      const arr = [];
      state.items.slice(0).forEach(item => {
        if (item.id !== id) {
          arr.push(item);
        }
      });
      state.items = arr;
    }
  },
  actions: {
    // addConfig 添加config
    async addConfig({ commit }, config) {
      commit(mutationConfigProcessing, true);
      try {
        const { data } = await request.post(CONFIGS, config);
        return data;
      } finally {
        commit(mutationConfigProcessing, false);
      }
    },
    // listConfig 获取config列表
    async listConfig({ commit }, params) {
      commit(mutationConfigProcessing, true);
      if (!params.offset) {
        commit(mutationConfigListReset);
      }
      try {
        const { data } = await request.get(CONFIGS, {
          params
        });
        data.configs.forEach(item => {
          item.isJSON = false;
          if (item.data[0] === "{") {
            item.isJSON = true;
            item.data = formatJSON(item.data);
          }
          if (item.beginDate) {
            item.beginDateDesc = formatDate(item.beginDate);
          }
          if (item.endDate) {
            item.endDateDesc = formatDate(item.endDate);
          }
          item.updatedAtDesc = formatDate(item.updatedAt);
          statusList.forEach(status => {
            if (item.status === status.value) {
              item.statusDesc = status.label;
            }
          });
        });
        commit(mutationConfigList, data);
        return data;
      } finally {
        commit(mutationConfigProcessing, false);
      }
    },
    // getConfigByID 通过id获取config
    async getConfigByID({ commit }, id) {
      commit(mutationConfigProcessing, true);
      try {
        const url = CONFIGS_ID.replace(":id", id);
        const { data } = await request.get(url);
        return data;
      } finally {
        commit(mutationConfigProcessing, false);
      }
    },
    // updateConfigByID 通过id指定config
    async updateConfigByID({ commit }, { id, data }) {
      commit(mutationConfigProcessing, true);
      try {
        const url = CONFIGS_ID.replace(":id", id);
        const res = await request.patch(url, data);
        return res.data;
      } finally {
        commit(mutationConfigProcessing, false);
      }
    },
    // removeConfigByID 通过id删除config
    async removeConfigByID({ commit }, id) {
      commit(mutationConfigProcessing, true);
      try {
        const url = CONFIGS_ID.replace(":id", id);
        await request.delete(url);
        commit(mutationConfigListDelete, id);
      } finally {
        commit(mutationConfigProcessing, false);
      }
    }
  }
};
