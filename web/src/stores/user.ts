import { defineStore } from "pinia";

import {
  USERS_ME,
  USERS_LOGIN,
  USERS_INNER_LOGIN,
  USERS_ME_DETAIL,
} from "../constants/url";
// eslint-disable-next-line
// @ts-ignore
import { sha256 } from "../helpers/crypto";
import request from "../helpers/request";

const hash = "JT";

function generatePassword(pass: string): string {
  return sha256(hash + sha256(pass + hash));
}

export const useUserStore = defineStore("user", {
  state: () => {
    return {
      processing: false,
      date: "",
      account: "",
      groups: [] as string[],
      roles: [] as string[],
      status: 0,
      name: "",
      email: "",
    };
  },
  getters: {
    anonymous: (state) => state.account == "",
  },
  actions: {
    // 填充客户信息
    _fillUserInfo(data: {
      account: string;
      date: string;
      roles: string[];
      groups: string[];
    }) {
      this.account = data.account;
      this.date = data.date;
      this.roles = data.roles || [];
      this.groups = data.groups || [];
    },
    // 登录
    async login(params: {
      account: string;
      password: string;
      captcha: string;
    }): Promise<void> {
      if (this.processing) {
        return;
      }
      try {
        this.processing = true;
        const resp = await request.get<{
          token: string;
        }>(USERS_LOGIN);
        const { token } = resp.data;
        const { data } = await request.post(
          USERS_INNER_LOGIN,
          {
            account: params.account,
            password: sha256(generatePassword(params.password) + token),
          },
          {
            headers: {
              "X-Captcha": params.captcha,
            },
          }
        );
        this._fillUserInfo(data);
      } finally {
        this.processing = false;
      }
    },
    // 退出登录
    async logout(): Promise<void> {
      if (this.processing) {
        return;
      }
      try {
        this.processing = true;
        await request.delete(USERS_ME);
        this.$reset();
      } finally {
        this.processing = false;
      }
    },
    // 拉取用户信息
    async fetch(): Promise<void> {
      if (this.processing) {
        return;
      }
      try {
        this.processing = true;
        const { data } = await request.get(USERS_ME);
        this._fillUserInfo(data);
      } finally {
        this.processing = false;
      }
    },
    // 注册
    async register(params: {
      account: string;
      password: string;
      captcha: string;
    }): Promise<void> {
      if (this.processing) {
        return;
      }
      try {
        // 如果密码小于6位或者纯数字
        if (params.password.length < 6 || /^\d+$/.exec(params.password)) {
          throw new Error("密码过于简单，请使用数字加字母且长度大于6位");
        }
        this.processing = true;
        await request.post(
          USERS_ME,
          {
            account: params.account,
            password: generatePassword(params.password),
          },
          {
            headers: {
              "X-Captcha": params.captcha,
            },
          }
        );
      } finally {
        this.processing = false;
      }
    },
    // 详细信息
    async detail(): Promise<void> {
      const { data } = await request.get<{
        status: number;
        email: string;
        name: string;
      }>(USERS_ME_DETAIL);
      this.status = data.status;
      this.email = data.email;
      this.name = data.name;
    },
    // 更新个人信息
    async update(params: Record<string, unknown>): Promise<void> {
      await request.patch(USERS_ME, params);
    },
  },
});
