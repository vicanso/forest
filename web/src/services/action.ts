import store from "../helpers/store";
import request from "../helpers/request";
import { USERS_ACTIONS } from "../constants/url";

const userActionKey = "userActions";
const userActions: any[] = [];

// 定时flush的间隔
const flushInterval = 60 * 1000;
let timer: number;

interface UserActionData {
  category: string;
  route: string;
  path: string;
  time: number;
  extra: any;
}

async function loadFromStore() {
  const data = await store.getItem(userActionKey);
  if (!data) {
    return;
  }
  const arr = JSON.parse(data);
  userActions.push(...arr);
}
loadFromStore();

async function flush() {
  store.removeItem(userActionKey);
  const actions = userActions.slice(0);
  userActions.length = 0;
  request.post(USERS_ACTIONS, {
    actions,
  });
}

export function addUserAction(data: UserActionData): void {
  // 每次添加新的action时，清空定时器
  clearTimeout(timer);
  userActions.push(data);
  if (userActions.length < 10) {
    store.setItem(userActionKey, JSON.stringify(userActions));
    // 重新启动定时器
    timer = setTimeout(flush, flushInterval);
    return;
  }
  flush();
}
