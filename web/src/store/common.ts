import request from "../helpers/request";
import {
  COMMONS_CAPTCHA,
  COMMONS_ROUTERS,
  COMMONS_STATUSES,
  COMMONS_RANDOM_KEYS,
} from "../constants/url";
import { DeepReadonly, reactive, readonly } from "vue";

interface Captcha {
  data: string;
  expiredAt: string;
  id: string;
  type: string;
}

interface Status {
  name: string;
  value: number;
}
interface Statuses {
  processing: boolean;
  items: Status[];
}
const statuses: Statuses = reactive({
  processing: false,
  items: [],
});

// 路由配置
interface Router {
  method: string;
  route: string;
}
interface Routers {
  processing: boolean;
  items: Router[];
}
const routers: Routers = reactive({
  processing: false,
  items: [],
});

// 随机数
interface RandomKeys {
  processing: boolean;
  items: string[];
}
const randomKeys: RandomKeys = reactive({
  processing: false,
  items: [],
});

// 仅读通用state
interface ReadonlyCommonState {
  statuses: DeepReadonly<Statuses>;
  routers: DeepReadonly<Routers>;
  randomKeys: DeepReadonly<RandomKeys>;
}

// commonGetCaptcha 获取图形验证码
export async function commonGetCaptcha(): Promise<Captcha> {
  const { data } = await request.get(COMMONS_CAPTCHA);
  return <Captcha>data;
}

// commonListStatus 获取状态列表
export async function commonListStatus(): Promise<void> {
  if (statuses.processing || statuses.items.length !== 0) {
    return;
  }
  statuses.processing = true;
  try {
    const { data } = await request.get(COMMONS_STATUSES);
    statuses.items = data.statuses || [];
  } finally {
    statuses.processing = false;
  }
}

// commonListRouter 获取路由列表
export async function commonListRouter(): Promise<void> {
  if (routers.processing || routers.items.length !== 0) {
    return;
  }
  routers.processing = true;
  try {
    const { data } = await request.get(COMMONS_ROUTERS);
    routers.items = data.routers || [];
  } finally {
    routers.processing = false;
  }
}

// commonListRandomKey 获取随机字符串
export async function commonListRandomKey(): Promise<void> {
  if (randomKeys.processing) {
    return;
  }
  randomKeys.processing = true;
  try {
    const { data } = await request.get(COMMONS_RANDOM_KEYS, {
      params: {
        size: 10,
        n: 5,
      },
    });
    randomKeys.items = data.keys || [];
  } finally {
    randomKeys.processing = false;
  }
}

// useCommonState 使用通用state
export default function useCommonState(): ReadonlyCommonState {
  return {
    statuses: readonly(statuses),
    routers: readonly(routers),
    randomKeys: readonly(randomKeys),
  };
}
