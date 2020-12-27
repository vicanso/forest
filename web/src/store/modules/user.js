import request from "@/helpers/request";
import {
  USERS_ME,
  USERS_LOGIN,
  USERS,
  USERS_ID,
  USERS_ROLES,
  USERS_LOGINS
} from "@/constants/url";
import {
  generatePassword,
  formatDate,
  queryOmitEmpty,
  addNoCacheQueryParam
} from "@/helpers/util";
import { sha256 } from "@/helpers/crypto";
import {
  listStatus,
  attachStatusDesc,
  attachUpdatedAtDesc
} from "@/store/modules/common";

const prefix = "user";
const mutationUserProcessing = `${prefix}.processing`;
const mutationUserInfo = `${prefix}.info`;

const mutationUserList = `${prefix}.list`;
const mutationUserListProcessing = `${mutationUserList}.processing`;

const mutationUserListRole = `${mutationUserList}.role`;
const mutationUserListRoleProcessing = `${mutationUserListRole}.processing`;

const mutationUserListGroup = `${mutationUserList}.group`;
const mutationUserListGroupProcessing = `${mutationUserListGroup}.processing`;

const mutationUserListMarketingGroup = `${mutationUserList}.marketingGroup`;
const mutationUserListMarketingGroupProcessing = `${mutationUserListMarketingGroup}.processing`;

const mutationUserUpdate = `${prefix}.update`;
const mutationUserUpdateProcessing = `${mutationUserUpdate}.processing`;

const mutationUserProfile = `${prefix}.profile`;
const mutationUserProfileProcessing = `${mutationUserProfile}.processing`;
const mutationUserProfileUpdate = `${mutationUserProfile}.update`;

const mutationUserListLogin = `${mutationUserList}.loign`;
const mutationUserListLoginProcessing = `${mutationUserListLogin}.processing`;

const state = {
  // 用户角色
  roleListProcessing: false,
  roles: null,
  // 用户分组
  groupListProcessing: false,
  groups: null,

  // 市场分组
  marketingGroupListProcessing: false,
  marketingGroups: null,

  // 默认为处理中（程序一开始则拉取用户信息）
  processing: true,
  info: {
    account: "",
    trackID: "",
    roles: null,
    groups: null
  },

  // 用户详情信息
  profileProcessing: false,
  profile: null,
  // 用户列表
  listProcessing: false,
  list: {
    data: null,
    count: -1
  },
  // 更新信息
  updateProcessing: false,

  // 用户登录
  userLoginListProcessing: false,
  userLogins: {
    data: null,
    count: -1
  }
};

function commitUserInfo(commit, data) {
  commit(
    mutationUserInfo,
    Object.assign(data, {
      account: data.account || "",
      name: data.name || "",
      trackID: data.trackID || "",
      roles: data.roles || [],
      groups: data.groups || []
    })
  );
}
function updateUserDesc(user) {
  user.updatedAtDesc = formatDate(user.updatedAt);
  attachUpdatedAtDesc(user);
  attachStatusDesc(user);
}

// listUserRole 获取用户角色列表
async function listUserRole({ commit }) {
  if (state.roles) {
    return {
      roles: state.roles
    };
  }
  commit(mutationUserListRoleProcessing, true);
  try {
    const { data } = await request.get(USERS_ROLES, {
      params: addNoCacheQueryParam()
    });
    data.roles = data.userRoles;
    commit(mutationUserListRole, data);
    return data;
  } finally {
    commit(mutationUserListRoleProcessing, false);
  }
}

// listUserGroup 获取用户分组列表
// async function listUserGroup({ commit }) {
//   if (state.groups) {
//     return {
//       groups: state.groups
//     };
//   }
//   commit(mutationUserListGroupProcessing, true);
//   try {
//     const { data } = await request.get(USERS_GROUPS, {
//       params: addNoCacheQueryParam()
//     });
//     commit(mutationUserListGroup, data);
//     return data;
//   } finally {
//     commit(mutationUserListGroupProcessing, false);
//   }
// }

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
      state.listProcessing = processing;
    },
    [mutationUserList](state, { users = [], count }) {
      if (count >= 0) {
        state.list.count = count;
      }
      users.forEach(updateUserDesc);
      state.list.data = users;
    },
    [mutationUserListRoleProcessing](state, processing) {
      state.roleListProcessing = processing;
    },
    [mutationUserListRole](state, { roles = [] }) {
      state.roles = roles;
    },
    [mutationUserListGroupProcessing](state, processing) {
      state.groupListProcessing = processing;
    },
    [mutationUserListGroup](state, { groups = [] }) {
      state.groups = groups;
    },
    [mutationUserListMarketingGroupProcessing](state, processing) {
      state.marketingGroupListProcessing = processing;
    },
    [mutationUserListMarketingGroup](state, { marketingGroups = [] }) {
      state.marketingGroups = marketingGroups;
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
          updateUserDesc(item);
        }
      });
      state.list.data = users;
    },
    [mutationUserProfileProcessing](state, processing) {
      state.processing = processing;
    },
    [mutationUserProfile](state, profile) {
      updateUserDesc(profile);
      state.profile = profile;
    },
    [mutationUserProfileUpdate](state, data) {
      Object.assign(state.info, data);
    },
    [mutationUserListLoginProcessing](state, processing) {
      state.userLoginListProcessing = processing;
    },
    [mutationUserListLogin](state, { userLogins = [], count = 0 }) {
      if (count >= 0) {
        state.userLogins.count = count;
      }
      userLogins.forEach(item => {
        item.createdAtDesc = formatDate(item.createdAt);
        const locations = [];
        ["country", "province", "city"].forEach(key => {
          if (item[key]) {
            locations.push(item[key]);
          }
        });
        item.location = locations.join(" ") || "--";
      });
      state.userLogins.data = userLogins;
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
    // updateMe 更新用户信息，若无更新字段，则可刷新session有效期
    async updateMe({ commit }, data = {}) {
      commit(mutationUserProfileProcessing, true);
      try {
        if (data.newPassword) {
          data.password = generatePassword(data.password);
          data.newPassword = generatePassword(data.newPassword);
        }
        await request.patch(USERS_ME, data);
        if (Object.keys(data).length !== 0) {
          commit(mutationUserProfileUpdate, data);
        }
      } finally {
        commit(mutationUserProfileProcessing, false);
      }
    },
    // listUser 获取用户
    async listUser({ commit }, params) {
      commit(mutationUserListProcessing, true);
      try {
        const { data } = await request.get(USERS, {
          params
        });
        commit(mutationUserList, data);
      } finally {
        commit(mutationUserListProcessing, false);
      }
    },
    listUserRole,
    listUserStatus: listStatus,
    // listUserGroup,
    // updateUserByID 更新用户信息
    async updateUserByID({ commit }, { id, data }) {
      commit(mutationUserUpdateProcessing, true);
      try {
        ["groups"].forEach(key => {
          if (data[key] && data[key].length === 0) {
            delete data[key];
          }
        });
        await request.patch(USERS_ID.replace(":id", id), data);
        commit(mutationUserUpdate, {
          id,
          data
        });
      } finally {
        commit(mutationUserUpdateProcessing, false);
      }
    },
    // getUserByID get user by id
    async getUserByID(_, id) {
      const { data } = await request.get(USERS_ID.replace(":id", id));
      return data;
    },
    // listUserLogins 获取用户登录记录
    async listUserLogins({ commit }, params) {
      commit(mutationUserListLoginProcessing, true);
      try {
        const { data } = await request.get(USERS_LOGINS, {
          params: queryOmitEmpty(params)
        });
        commit(mutationUserListLogin, data);
      } finally {
        commit(mutationUserListLoginProcessing, false);
      }
    }
  }
};
