import request from "@/helpers/request";
import {
  USERS_ME,
  USERS_LOGIN,
  USERS,
  USERS_ID,
  COMMONS_LIST_USER_ROLE,
  COMMONS_LIST_USER_STATUS
} from "@/constants/url";
import { generatePassword, formatDate } from "@/helpers/util";
import { sha256 } from "@/helpers/crypto";

const mutationUserProcessing = "user.processing";
const mutationUserInfo = "user.info";

const mutationUserListProcessing = "user.list.processing";
const mutationUserList = "user.list";

const mutationUserListRoleProcessing = "user.list.role.processing";
const mutationUserListRole = "user.list.role";

const mutationUserListStatusProcessing = "user.list.status.processing";
const mutationuserListStatus = "user.list.status";

const mutationUserUpdateProcessing = "user.updae.processing";
const mutationUserUpdate = "user.update";

const state = {
  roleListProcessing: false,
  roles: null,
  statusListProcessing: false,
  statuses: null,
  // 默认为处理中（程序一开始则拉取用户信息）
  processing: true,
  info: {
    account: "",
    trackID: "",
    roles: null
  },
  userListProcessing: false,
  list: {
    data: null,
    count: -1
  },
  updateProcessing: false
};

function commitUserInfo(commit, data) {
  commit(mutationUserInfo, {
    account: data.account || "",
    trackID: data.trackID || "",
    roles: data.roles || []
  });
}

function updateUserRoleDesc(user) {
  const { roles } = state;
  const roleDescList = [];
  user.roles.map(role => {
    roles.forEach(roleDesc => {
      if (roleDesc.value === role) {
        roleDescList.push(roleDesc.name);
      }
    });
  });
  user.roleDescList = roleDescList;
}

// listUserStatus 获取用户状态列表
async function listUserStatus({ commit }) {
  if (state.statuses) {
    return;
  }
  commit(mutationUserListStatusProcessing, true);
  try {
    const { data } = await request.get(COMMONS_LIST_USER_STATUS);
    commit(mutationuserListStatus, data);
  } finally {
    commit(mutationUserListStatusProcessing, false);
  }
}

// listUserRole 获取用户角色列表
async function listUserRole({ commit }) {
  if (state.roles) {
    return;
  }
  commit(mutationUserListRoleProcessing, true);
  try {
    const { data } = await request.get(COMMONS_LIST_USER_ROLE);
    commit(mutationUserListRole, data);
  } finally {
    commit(mutationUserListRoleProcessing, false);
  }
}

export default {
  state,
  mutations: {
    [mutationUserProcessing](state, processing) {
      state.processing = processing;
    },
    [mutationUserInfo](state, value) {
      Object.assign(state.info, value);
    },
    [mutationUserListProcessing](state, processing) {
      state.userListProcessing = processing;
    },
    [mutationUserList](state, { users, count }) {
      if (count >= 0) {
        state.list.count = count;
      }
      state.list.data = users;
    },
    [mutationUserListRoleProcessing](state, processing) {
      state.roleListProcessing = processing;
    },
    [mutationUserListRole](state, { roles }) {
      state.roles = roles;
    },
    [mutationUserListStatusProcessing](state, processing) {
      state.statusListProcessing = processing;
    },
    [mutationuserListStatus](state, { statuses }) {
      state.statuses = statuses;
    },
    [mutationUserUpdateProcessing](state, processing) {
      state.updateProcessing = processing;
    },
    [mutationUserUpdate](state, { id, data }) {
      if (!state.list.data) {
        return;
      }
      const users = state.list.data.slice(0);
      users.forEach(item => {
        if (item.id === id) {
          Object.assign(item, data);
          updateUserRoleDesc(item);
        }
      });
      state.list.data = users;
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
    // updateUser 更新用户信息，若无更新字段，则可刷新session有效期
    async updateUser() {
      await request.patch(USERS_ME, {});
    },
    // listUser 获取用户
    async listUser({ commit }, params) {
      commit(mutationUserListProcessing, true);
      try {
        await listUserRole({ commit });
        await listUserStatus({ commit });
        const { data } = await request.get(USERS, {
          params
        });
        const { statuses } = state;
        data.users.forEach(item => {
          if (!item.roles) {
            item.roles = [];
          }
          updateUserRoleDesc(item);
          item.updatedAtDesc = formatDate(item.updatedAt);
          statuses.forEach(status => {
            if (status.value === item.status) {
              item.statusDesc = status.name;
            }
          });
          if (!item.statusDesc) {
            item.statusDesc = "未知";
          }
        });
        commit(mutationUserList, data);
      } finally {
        commit(mutationUserListProcessing, false);
      }
    },
    listUserRole,
    listUserStatus,
    async updateUserByID({ commit }, { id, data }) {
      commit(mutationUserUpdateProcessing, true);
      try {
        await request.patch(USERS_ID.replace(":id", id), data);
        commit(mutationUserUpdate, {
          id,
          data
        });
      } finally {
        commit(mutationUserUpdateProcessing, false);
      }
    }
  }
};
