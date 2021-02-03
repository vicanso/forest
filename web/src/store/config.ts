import request from "../helpers/request";
import { CONFIG_ENABLED, CONFIG_DISABLED } from "../constants/common";
import { CONFIGS, CONFIGS_ID, CONFIGS_CURRENT_VALID } from "../constants/url";
import { DeepReadonly, reactive, readonly } from "vue";

// 配置状态
interface Status {
  label: string;
  value: number;
}
// 配置状态列表
interface Statuses {
  items: Status[];
}

const statues: Statuses = reactive({
  items: [
    {
      label: "启用",
      value: CONFIG_ENABLED,
    },
    {
      label: "禁用",
      value: CONFIG_DISABLED,
    },
  ],
});

// 配置信息
interface Config {
  id: number;
  createdAt: string;
  updatedAt: string;
  status: number;
  name: string;
  category: string;
  owner: string;
  data: string;
  startedAt: string;
  endedAt: string;
}
interface Configs {
  processing: boolean;
  current?: Config;
  items: Config[];
  count: number;
}
const configs: Configs = reactive({
  processing: false,
  items: [],
  count: -1,
});

interface CurrentValidConfig {
  processing: boolean;
  data: string;
}
const currentValidConfig: CurrentValidConfig = reactive({
  processing: false,
  data: "",
});

// 仅读配置state
interface ReadonlyConfigState {
  configs: DeepReadonly<Configs>;
  currentValidConfig: DeepReadonly<CurrentValidConfig>;
  statuses: DeepReadonly<Statuses>;
}

// configAdd
export async function configAdd(params: {
  name: string;
  status: number;
  category: string;
  startedAt: string;
  endedAt: string;
  data: string;
}): Promise<void> {
  if (configs.processing) {
    return;
  }
  configs.processing = true;
  try {
    const { data } = await request.post(CONFIGS, params);
    configs.current = <Config>data;
  } finally {
    configs.processing = false;
  }
}

// configList 查询配置列表
export async function configList(params: {
  name?: string;
  category?: string;
  limit?: number;
  offset?: number;
}): Promise<void> {
  if (configs.processing) {
    return;
  }
  if (!params.limit) {
    params.limit = 50;
  }
  configs.processing = true;
  try {
    const { data } = await request.get(CONFIGS, {
      params,
    });
    const count = data.count || 0;
    if (count >= 0) {
      configs.count = count;
    }
    configs.items = data.configurations || [];
  } finally {
    configs.processing = false;
  }
}

// configListValid 查询有效配置
export async function configListValid(): Promise<void> {
  if (currentValidConfig.processing) {
    return;
  }
  currentValidConfig.processing = true;
  try {
    const { data } = await request.get(CONFIGS_CURRENT_VALID);
    currentValidConfig.data = JSON.stringify(data, null, 2);
  } finally {
    currentValidConfig.processing = false;
  }
}

// configFindByID 通过ID查询config
export async function configFindByID(id: number): Promise<Config> {
  const url = CONFIGS_ID.replace(":id", `${id}`);
  const { data } = await request.get(url);
  return <Config>data;
}

// configUpdateByID 通过ID更新config
export async function configUpdateByID(params: {
  id: number;
  data: any;
}): Promise<void> {
  if (configs.processing) {
    return;
  }
  configs.processing = true;
  try {
    const url = CONFIGS_ID.replace(":id", `${params.id}`);
    const { data } = await request.patch(url, params.data);
    const items = configs.items.slice(0);
    items.forEach((item) => {
      if (item.id === params.id) {
        Object.assign(item, params.data);
      }
    });
    configs.current = <Config>data;
    configs.items = items;
  } finally {
    configs.processing = false;
  }
}

export default function useConfigState(): ReadonlyConfigState {
  return {
    configs: readonly(configs),
    currentValidConfig: readonly(currentValidConfig),
    statuses: readonly(statues),
  };
}
