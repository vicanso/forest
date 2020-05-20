import request from "@/helpers/request";
import { USERS_ME, USERS_LOGIN } from "@/constants/url";
import { generatePassword } from "@/helpers/util";
import { sha256 } from "@/helpers/crypto";

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
    // fetchUserInfo 获取用户信息
    async fetchUserInfo({ commit }) {
      commit(mutationUserProcessing, true);
      try {
        const { data } = await request.get(USERS_ME);
        commitUserInfo(commit, data);
      } finally {
        commit(mutationUserProcessing, false);
      }
    },
    // logout 退出登录
    async logout({ commit }) {
      // 设置处理中
      commit(mutationUserProcessing, true);
      try {
        await request.delete(USERS_ME);
        commitUserInfo(commit, {});
      } finally {
        commit(mutationUserProcessing, false);
      }
    },
    // login 用户登录
    async login({ commit }, { account, password, captcha }) {
      commit(mutationUserProcessing, true);
      try {
        // 先获取登录用的token
        const res = await request.get(USERS_LOGIN);
        const { token } = res.data;
        // 根据token与密码生成登录密码
        const { data } = await request.post(
          USERS_LOGIN,
          {
            account,
            password: sha256(generatePassword(password) + token)
          },
          {
            headers: {
              // 图形验证码
              "X-Captcha": captcha
            }
          }
        );
        commitUserInfo(commit, data);
      } finally {
        commit(mutationUserProcessing, false);
      }
    },
    // register 用户注册
    async register({ commit }, { account, password, captcha }) {
      commit(mutationUserProcessing, true);
      try {
        await request.post(
          USERS_ME,
          {
            account,
            // 密码加密
            password: generatePassword(password)
          },
          {
            headers: {
              "X-Captcha": captcha
            }
          }
        );
      } finally {
        commit(mutationUserProcessing, false);
      }
    },
    async updateUser() {
      await request.patch(USERS_ME, {});
    }
  }
};
