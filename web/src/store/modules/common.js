import request from "@/helpers/request";
import { COMMONS_CAPTCHA, COMMONS_ROUTERS } from "@/constants/url";

const mutationCommonProcessing = "common.processing";
const mutationCommonRouterList = "common.router.list";

const state = {
  processing: false,
  routers: null
};

export default {
  state,
  mutations: {
    [mutationCommonProcessing](state, processing) {
      state.processing = processing;
    },
    [mutationCommonRouterList](state, routers) {
      if (routers) {
        routers.forEach(item => {
          item.key = `${item.method} ${item.path}`;
        });
      }
      state.routers = routers;
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
        return;
      }
      commit(mutationCommonProcessing, true);
      try {
        const { data } = await request.get(COMMONS_ROUTERS);
        commit(mutationCommonRouterList, data.routers);
      } finally {
        commit(mutationCommonProcessing, false);
      }
    }
  }
};
