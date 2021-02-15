import { reactive, readonly, DeepReadonly } from "vue";

import request from "../helpers/request";
import { sha256 } from "../helpers/crypto";
import {
  USERS_ME,
  USERS_LOGIN,
  USERS_INNER_LOGIN,
  USERS_LOGINS,
  USERS_ROLES,
  USERS,
  USERS_ID,
  USERS_ME_DETAIL,
} from "../constants/url";
import { generatePassword } from "../helpers/util";

// 用户信息
interface UserInfo {
  processing: boolean;
  date: string;
  account: string;
  groups: string[];
  roles: string[];
}
const info: UserInfo = reactive({
  processing: false,
  date: "",
  account: "",
  groups: [],
  roles: [],
});

interface UserRole {
  name: string;
  value: string;
}
interface UserRoles {
  processing: boolean;
  items: UserRole[];
}

const roles: UserRoles = reactive({
  processing: false,
  items: [],
});

// 用户账户信息
interface UserAccount {
  id: number;
  account: string;
  groups: string[];
  roles: string[];
  email: string;
  status: number;
}
// 用户账户列表
interface UserAccounts {
  processing: boolean;
  count: number;
  items: UserAccount[];
}
const users: UserAccounts = reactive({
  processing: false,
  count: -1,
  items: [],
});

// 用户登录信息
interface UserLoginRecord {
  account: string;
  userAgent?: string;
  ip?: string;
  trackID?: string;
  sessionID?: string;
  xForwardedFor?: string;
  country?: string;
  province?: string;
  city?: string;
  isp?: string;
  updatedAt?: string;
  createdAt?: string;
}
// 用户登录列表
interface UserLoginRecords {
  processing: boolean;
  count: number;
  items: UserLoginRecord[];
}

const logins: UserLoginRecords = reactive({
  processing: false,
  count: -1,
  items: [],
});

// 仅读用户state
interface ReadonlyUserState {
  info: DeepReadonly<UserInfo>;
  logins: DeepReadonly<UserLoginRecords>;
  users: DeepReadonly<UserAccounts>;
  roles: DeepReadonly<UserRoles>;
}

function fillUserInfo(data: UserInfo) {
  info.account = data.account;
  info.date = data.date;
  info.roles = data.roles || [];
  info.groups = data.groups || [];
}

// userFetchInfo 拉取用户信息
export async function userFetchInfo(): Promise<void> {
  // TODO 是否需要针对并发调用出错
  if (info.processing) {
    return;
  }
  try {
    info.processing = true;
    const { data } = await request.get(USERS_ME);
    fillUserInfo(<UserInfo>data);
  } finally {
    info.processing = false;
  }
}

// userFetchDetail 拉取个人详细信息， 仅在个人信息页使用，因此不记录state
export async function userFetchDetail(): Promise<UserAccount> {
  const { data } = await request.get(USERS_ME_DETAIL);
  return <UserAccount>data;
}

// userLogin 用户登录
export async function userLogin(params: {
  account: string;
  password: string;
  captcha: string;
}): Promise<void> {
  if (info.processing) {
    return;
  }
  try {
    info.processing = true;
    const resp = await request.get(USERS_LOGIN);
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
    fillUserInfo(<UserInfo>data);
  } finally {
    info.processing = false;
  }
}

// userRegister 用户注册
export async function userRegister(params: {
  account: string;
  password: string;
  captcha: string;
}): Promise<void> {
  if (info.processing) {
    return;
  }
  try {
    info.processing = true;
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
    info.processing = false;
  }
}

// userLogout 退出登录
export async function userLogout(): Promise<void> {
  if (info.processing) {
    return;
  }
  try {
    info.processing = true;
    await request.delete(USERS_ME);
    fillUserInfo({
      account: "",
      roles: [],
      groups: [],
      date: "",
      processing: false,
    });
  } finally {
    info.processing = false;
  }
}

// userUpdate 更新用户信息
export async function userUpdate(params: {
  password?: string;
  newPassword?: string;
  roles?: string[];
}): Promise<void> {
  if (info.processing) {
    return;
  }
  try {
    info.processing = true;
    const data = Object.assign({}, params);
    if (data.password) {
      data.password = generatePassword(data.password);
    }
    if (data.newPassword) {
      data.newPassword = generatePassword(data.newPassword);
    }
    await request.patch(USERS_ME, data);
  } finally {
    info.processing = false;
  }
}

// userListLogin 查询用户登录记录
export async function userListLogin(params: {
  account?: string;
  begin: string;
  end: string;
  limit: number;
  offset: number;
  order?: string;
}): Promise<void> {
  if (logins.processing) {
    return;
  }
  try {
    logins.processing = true;
    const { data } = await request.get(USERS_LOGINS, {
      params,
    });
    const count = data.count || 0;
    if (count >= 0) {
      logins.count = count;
    }
    logins.items = data.userLogins || [];
  } finally {
    logins.processing = false;
  }
}

// userLoginClear 清空登录记录
export function userLoginClear(): void {
  logins.count = -1;
  logins.items.length = 0;
}

// userList 查询用户
export async function userList(params: {
  keyword?: string;
  limit: number;
  offset: number;
  role?: string;
  status?: number;
  order?: string;
}): Promise<void> {
  if (users.processing) {
    return;
  }
  try {
    users.processing = true;
    const { data } = await request.get(USERS, {
      params,
    });
    const count = data.count || 0;
    if (count >= 0) {
      users.count = count;
    }
    users.items = data.users || [];
  } finally {
    users.processing = false;
  }
}

// userListClear 清空用户记录
export function userListClear(): void {
  users.count = -1;
  users.items.length = 0;
}

// userFindByID 通过ID查询用户
export async function userFindByID(id: number): Promise<UserAccount> {
  const { data } = await request.get(USERS_ID.replace(":id", `${id}`));
  return <UserAccount>data;
}

// userUpdateByID 通过ID更新用户
export async function userUpdateByID(params: {
  id: number;
  data: Record<string, unknown>;
}): Promise<void> {
  if (users.processing) {
    return;
  }
  const data = Object.assign({}, params.data);
  try {
    users.processing = true;
    // 如果groups未设置，则清空
    ["groups"].forEach((key) => {
      const value = data[key];
      if (value && Array.isArray(value) && value.length === 0) {
        delete data[key];
      }
    });
    await request.patch(USERS_ID.replace(":id", `${params.id}`), data);
    const items = users.items.slice(0);
    items.forEach((item) => {
      if (item.id === params.id) {
        Object.assign(item, data);
      }
    });
    users.items = items;
  } finally {
    users.processing = false;
  }
}

// userListRole 查询用户角色
export async function userListRole(): Promise<void> {
  if (roles.processing || roles.items.length !== 0) {
    return;
  }
  try {
    roles.processing = true;
    const { data } = await request.get(USERS_ROLES);
    roles.items = data.userRoles || [];
  } finally {
    roles.processing = false;
  }
}

const state = {
  info: readonly(info),
  logins: readonly(logins),
  users: readonly(users),
  roles: readonly(roles),
};
// useUserState 使用用户state
export default function useUserState(): ReadonlyUserState {
  return state;
}
