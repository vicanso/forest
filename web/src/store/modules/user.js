import request from "@/helpers/request";

import { USERS_ME } from "@/constants/url";

const mutationUserProcessing = "user.processing";
const mutationUserInfo = "user.info";

const state = {
  // 默认为处理中（程序一开始则拉取用户信息）
  processing: true,
  info: {
    account: "",
    trackID: ""
  }
};

function commitUserInfo(commit, data) {
  commit(mutationUserInfo, {
    account: data.account || "",
    trackID: data.trackID || ""
  });
}

export default {
  state,
  mutations: {
    [mutationUserProcessing](state, processing) {
      state.processing = processing;
    },
    [mutationUserInfo](state, value) {
      Object.assign(state.info, value);
    }
  },
  actions: {
    async fetchUserInfo({ commit }) {
      commit(mutationUserProcessing, true);
      try {
        const { data } = await request.get(USERS_ME);
        commitUserInfo(commit, data);
      } finally {
        commit(mutationUserProcessing, false);
      }
    },
    async logout({ commit }) {
      commit(mutationUserProcessing, true);
      try {
        await request.delete(USERS_ME);
        commitUserInfo(commit, {});
      } finally {
        commit(mutationUserProcessing, false);
      }
    }
  }
};
