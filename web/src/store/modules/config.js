import request from "@/helpers/request";
import { CONFIGS, CONFIGS_ID } from "@/constants/url";

const mutationConfigProcessing = "config.processing";

const state = {
  processing: false
};

export default {
  state,
  mutations: {
    [mutationConfigProcessing](state, processing) {
      state.processing = processing;
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
      try {
        const { data } = await request.get(CONFIGS, {
          params
        });
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
    }
  }
};
