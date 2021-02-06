import request from "../helpers/request";
import {
  FLUXES_TRACKERS,
  FLUXES_HTTP_ERRORS,
  FLUXES_TAG_VALUES,
  FLUXES_ACTIONS,
} from "../constants/url";
import { DeepReadonly, reactive, readonly } from "vue";

// 用户行为轨迹
interface UserTracker {
  _time: string;
  account: string;
  action: string;
  hostname: string;
  ip: string;
  result: string;
  sid: string;
  tid: string;
  form: string;
  query: string;
  params: string;
  error: string;
}
interface UserTrackers {
  processing: boolean;
  items: UserTracker[];
}
const userTrackers: UserTrackers = reactive({
  processing: false,
  items: [],
});

// 用户行为轨迹类型
interface UserTrackerActions {
  processing: boolean;
  items: string[];
}
const userTrackerActions: UserTrackerActions = reactive({
  processing: false,
  items: [],
});

// 客户端行为记录类型
interface ClientActionCategories {
  processing: boolean;
  items: string[];
}
const clientActionCategories: ClientActionCategories = reactive({
  processing: false,
  items: [],
});

// 客户端行为记录
interface ClientAction {
  _time: string;
  account: string;
  category: string;
  hostname: string;
  path: string;
  result: string;
  route: string;
  tid: string;
  message: string;
}
interface ClientActions {
  processing: boolean;
  items: ClientAction[];
}
const clientActions: ClientActions = reactive({
  processing: false,
  items: [],
});

interface HTTPError {
  _time: string;
  account: string;
  category: string;
  error: string;
  exception: boolean;
  hostname: string;
  ip: string;
  method: string;
  route: string;
  sid: string;
  status: number;
  tid: string;
  uri: string;
}
interface HTTPErrors {
  processing: boolean;
  items: HTTPError[];
}
const httpErrors: HTTPErrors = reactive({
  processing: false,
  items: [],
});

// HTTP出错类型
interface HTTPErrorCategories {
  processing: boolean;
  items: string[];
}
const httpErrorCategories: HTTPErrorCategories = reactive({
  processing: false,
  items: [],
});

interface ReadonlyFluxState {
  userTrackers: DeepReadonly<UserTrackers>;
  userTrackerActions: DeepReadonly<UserTrackerActions>;
  httpErrorCategories: DeepReadonly<HTTPErrorCategories>;
  httpErrors: DeepReadonly<HTTPErrors>;
  clientActions: DeepReadonly<ClientActions>;
  clientActionCategories: DeepReadonly<ClientActionCategories>;
}

// fluxListUserTracker 查询用户跟踪轨迹记录
export async function fluxListUserTracker(params: {
  account?: string;
  action?: string;
  begin: string;
  end: string;
  limit: number;
  offset: number;
  result?: string;
}): Promise<void> {
  if (userTrackers.processing) {
    return;
  }
  try {
    userTrackers.processing = true;
    const { data } = await request.get(FLUXES_TRACKERS, {
      params,
    });
    userTrackers.items = data.trackers || [];
  } finally {
    userTrackers.processing = false;
  }
}

// fluxListUserTrackerClear 清除tracker记录
export function fluxListUserTrackerClear(): void {
  userTrackers.items.length = 0;
}

// fluxListUserTrackAction 查询用户轨迹action列表
export async function fluxListUserTrackAction(): Promise<void> {
  if (userTrackerActions.processing || userTrackerActions.items.length !== 0) {
    return;
  }
  try {
    userTrackerActions.processing = true;
    const url = FLUXES_TAG_VALUES.replace(
      ":measurement",
      "userTracker"
    ).replace(":tag", "action");
    const { data } = await request.get(url);
    userTrackerActions.items = data.values || [];
  } finally {
    userTrackerActions.processing = false;
  }
}

// fluxListHTTPCategory 查询HTTP出错类型列表
export async function fluxListHTTPCategory(): Promise<void> {
  if (
    httpErrorCategories.processing ||
    httpErrorCategories.items.length !== 0
  ) {
    return;
  }
  try {
    httpErrorCategories.processing = true;
    const url = FLUXES_TAG_VALUES.replace(":measurement", "httpError").replace(
      ":tag",
      "category"
    );
    const { data } = await request.get(url);
    httpErrorCategories.items = data.values || [];
  } finally {
    httpErrorCategories.processing = false;
  }
}

// fluxListHTTPError 查询HTTP出错记录
export async function fluxListHTTPError(params: {
  account?: string;
  category?: string;
  begin: string;
  end: string;
  exception?: string;
  limit: number;
  offset: number;
}): Promise<void> {
  if (httpErrors.processing) {
    return;
  }
  try {
    httpErrors.processing = true;
    const { data } = await request.get(FLUXES_HTTP_ERRORS, {
      params,
    });
    httpErrors.items = data.httpErrors || [];
  } finally {
    httpErrors.processing = false;
  }
}

// fluxListHTTPErrorClear 清除http出错列表
export function fluxListHTTPErrorClear(): void {
  httpErrors.items.length = 0;
}

// fluxListClientActionCategory 查询客户端行为分类
export async function fluxListClientActionCategory(): Promise<void> {
  if (
    clientActionCategories.processing ||
    clientActionCategories.items.length !== 0
  ) {
    return;
  }
  try {
    clientActionCategories.processing = true;
    const url = FLUXES_TAG_VALUES.replace(":measurement", "userAction").replace(
      ":tag",
      "category"
    );
    const { data } = await request.get(url);
    clientActionCategories.items = data.value || [];
  } finally {
    clientActionCategories.processing = false;
  }
}

// fluxListClientAction 查询客户端行为记录
export async function fluxListClientAction(params: {
  account?: string;
  category?: string;
  begin: string;
  end: string;
  limit: number;
  offset: number;
  result?: string;
}): Promise<void> {
  if (clientActions.processing) {
    return;
  }
  try {
    clientActions.processing = true;
    const { data } = await request.get(FLUXES_ACTIONS, {
      params,
    });
    clientActions.items = data.actions || [];
  } finally {
    clientActions.processing = false;
  }
}

// fluxListClientActionClear 清空客户端行为
export function fluxListClientActionClear(): void {
  clientActions.items.length = 0;
}

const state = {
  userTrackers: readonly(userTrackers),
  userTrackerActions: readonly(userTrackerActions),
  httpErrorCategories: readonly(httpErrorCategories),
  httpErrors: readonly(httpErrors),
  clientActions: readonly(clientActions),
  clientActionCategories: readonly(clientActionCategories),
};
export default function useFluxState(): ReadonlyFluxState {
  return state;
}
