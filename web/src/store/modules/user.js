import request from "@/helpers/request";
import {
  USERS_ME,
  USERS_LOGIN,
  USERS,
  USERS_ID,
  USERS_ME_PROFILE,
  COMMONS_LIST_USER_ROLE,
  COMMONS_LIST_USER_STATUS,
  COMMONS_LIST_USER_GROUPS
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

const mutationUserListGroupProcessing = "user.list.group.processing";
const mutationuserListGroup = "user.list.group";

const mutationUserUpdateProcessing = "user.update.processing";
const mutationUserUpdate = "user.update";

const mutationUserProfileProcessing = "user.profile.processing";
const mutationUserProfile = "user.profile";
const mutationUserProfileUpdate = "user.profile.update";

const state = {
  roleListProcessing: false,
  roles: null,
  statusListProcessing: false,
  statuses: null,
  groupListProcessing: false,
  groups: null,

  // 默认为处理中（程序一开始则拉取用户信息）
  processing: true,
  info: {
    account: "",
    trackID: "",
    roles: null
  },

  // 用户详情信息
  profileProcessing: false,
  profile: null,

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

function updateDescList(user, key) {
  if (!user[key]) {
    user[key] = [];
  }
  const descList = [];
  user[key].map(item => {
    state[key].forEach(desc => {
      if (desc.value === item) {
        descList.push(desc.name);
      }
    });
  });
  user[`${key}Desc`] = descList;
}

function updateUserRoleDesc(user) {
  updateDescList(user, "roles");
}

function updateUserGroupDesc(user) {
  updateDescList(user, "groups");
}

function updateStatusDesc(user) {
  state.statuses.forEach(status => {
    if (status.value === user.status) {
      user.statusDesc = status.name;
    }
  });
  if (!user.statusDesc) {
    user.statusDesc = "未知";
  }
}

function updateUserDesc(user) {
  user.updatedAtDesc = formatDate(user.updatedAt);
  updateStatusDesc(user);
  updateUserRoleDesc(user);
  updateUserGroupDesc(user);
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

// listUserGroup 获取用户分组列表
async function listUserGroup({ commit }) {
  if (state.groups) {
    return;
  }
  commit(mutationUserListGroupProcessing, true);
  try {
    const { data } = await request.get(COMMONS_LIST_USER_GROUPS);
    commit(mutationuserListGroup, data);
  } finally {
    commit(mutationUserListGroupProcessing, false);
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
    [mutationUserListGroupProcessing](state, processing) {
      state.groupListProcessing = processing;
    },
    [mutationuserListGroup](state, { groups }) {
      state.groups = groups;
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
      state.profile = profile;
    },
    [mutationUserProfileUpdate](state, data) {
      Object.assign(state.profile, data);
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
        commit(mutationUserProfileUpdate, data);
      } finally {
        commit(mutationUserProfileProcessing, false);
      }
    },
    // listUser 获取用户
    async listUser({ commit }, params) {
      commit(mutationUserListProcessing, true);
      try {
        await listUserRole({ commit });
        await listUserStatus({ commit });
        await listUserGroup({ commit });
        const { data } = await request.get(USERS, {
          params
        });
        data.users.forEach(item => {
          updateUserDesc(item);
        });
        commit(mutationUserList, data);
      } finally {
        commit(mutationUserListProcessing, false);
      }
    },
    listUserRole,
    listUserStatus,
    listUserGroup,
    // updateUserByID 更新用户信息
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
    },
    // getUserProfile 获取用户详细信息
    async getUserProfile({ commit }) {
      if (state.profile) {
        return;
      }
      commit(mutationUserProfileProcessing, true);
      try {
        await listUserRole({ commit });
        await listUserStatus({ commit });
        await listUserGroup({ commit });
        const { data } = await request.get(USERS_ME_PROFILE);
        updateUserDesc(data);
        commit(mutationUserProfile, data);
      } finally {
        commit(mutationUserProfileProcessing, false);
      }
    }
  }
};
