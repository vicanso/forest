import request from "../helpers/request";
import {
  FLUXES_TRACKERS,
  FLUXES_HTTP_ERRORS,
  FLUXES_TAG_VALUES,
  FLUXES_ACTIONS,
  FLUXES_REQUESTS,
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
  count: number;
  flux: string;
}
const userTrackers: UserTrackers = reactive({
  processing: false,
  items: [],
  count: -1,
  flux: "",
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
  count: number;
  flux: string;
}
const clientActions: ClientActions = reactive({
  processing: false,
  items: [],
  count: -1,
  flux: "",
});

// HTTPError 客户端HTTP请求出错记录
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
  count: number;
  flux: string;
}
const httpErrors: HTTPErrors = reactive({
  processing: false,
  items: [],
  count: -1,
  flux: "",
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

// 后端HTTP请求记录
interface Request {
  _time: string;
  hostname: string;
  addr: string;
  service: string;
  method: string;
  route: string;
  uri: string;
  status: number;
  reused: boolean;
  dnsUse: number;
  tcpUse: number;
  tlsUse: number;
  processingUse: number;
  use: number;
  result: string;
  errCategory: string;
  error: string;
  exception: boolean;
}
interface Requests {
  processing: boolean;
  items: Request[];
  count: number;
  flux: string;
}
const requests: Requests = reactive({
  processing: false,
  items: [],
  count: -1,
  flux: "",
});

// request 服务名称
interface RequestServices {
  processing: boolean;
  items: string[];
}
const requestServices: RequestServices = reactive({
  processing: false,
  items: [],
});

// RequestRoutes 请求路由
interface RequestRoutes {
  processing: boolean;
  items: string[];
}
const requestRoutes: RequestRoutes = reactive({
  processing: false,
  items: [],
});

interface ReadonlyFluxState {
  userTrackers: DeepReadonly<UserTrackers>;
  userTrackerActions: DeepReadonly<UserTrackerActions>;
  httpErrorCategories: DeepReadonly<HTTPErrorCategories>;
  httpErrors: DeepReadonly<HTTPErrors>;
  requests: DeepReadonly<Requests>;
  requestServices: DeepReadonly<RequestServices>;
  requestRoutes: DeepReadonly<RequestRoutes>;
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
    userTrackers.count = data.count || 0;
    userTrackers.flux = data.flux || "";
  } finally {
    userTrackers.processing = false;
  }
}

// fluxListUserTrackerClear 清除tracker记录
export function fluxListUserTrackerClear(): void {
  userTrackers.items.length = 0;
  userTrackers.count = -1;
  userTrackers.flux = "";
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
    httpErrors.count = data.count || 0;
    httpErrors.flux = data.flux || "";
  } finally {
    httpErrors.processing = false;
  }
}

// fluxListHTTPErrorClear 清除http出错列表
export function fluxListHTTPErrorClear(): void {
  httpErrors.items.length = 0;
  httpErrors.count = -1;
  httpErrors.flux = "";
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
    clientActionCategories.items = data.values || [];
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
    clientActions.count = data.count || 0;
    clientActions.flux = data.flux || "";
  } finally {
    clientActions.processing = false;
  }
}

// fluxListClientActionClear 清空客户端行为
export function fluxListClientActionClear(): void {
  clientActions.items.length = 0;
  clientActions.count = -1;
  clientActions.flux = "";
}

// fluxListRequest 查询后端请求记录
export async function fluxListRequest(params: {
  route?: string;
  service?: string;
  errCategory?: string;
  begin: string;
  end: string;
  exception?: string;
  limit: number;
  offset: number;
}): Promise<void> {
  if (requests.processing) {
    return;
  }
  try {
    requests.processing = true;
    const { data } = await request.get(FLUXES_REQUESTS, {
      params,
    });
    requests.items = data.requests || [];
    requests.count = data.count || 0;
    requests.flux = data.flux || "";
  } finally {
    requests.processing = false;
  }
}

// fluxListRequestClear 清除request列表
export function fluxListRequestClear(): void {
  requests.items.length = 0;
  requests.count = -1;
  requests.flux = "";
}

// fluxListRequestService 获取request中的service列表
export async function fluxListRequestService(): Promise<void> {
  if (requestServices.processing || requestServices.items.length !== 0) {
    return;
  }
  try {
    requestServices.processing = true;
    const url = FLUXES_TAG_VALUES.replace(
      ":measurement",
      "httpRequest"
    ).replace(":tag", "service");
    const { data } = await request.get(url);
    requestServices.items = data.values || [];
  } finally {
    requestServices.processing = false;
  }
}

// fluxListRequestRoute 获取request中的route列表
export async function fluxListRequestRoute(): Promise<void> {
  if (requestRoutes.processing || requestRoutes.items.length !== 0) {
    return;
  }
  try {
    requestRoutes.processing = true;
    const url = FLUXES_TAG_VALUES.replace(
      ":measurement",
      "httpRequest"
    ).replace(":tag", "route");
    const { data } = await request.get(url);
    requestRoutes.items = data.values || [];
  } finally {
    requestRoutes.processing = false;
  }
}

const state = {
  userTrackers: readonly(userTrackers),
  userTrackerActions: readonly(userTrackerActions),
  httpErrorCategories: readonly(httpErrorCategories),
  httpErrors: readonly(httpErrors),
  requests: readonly(requests),
  requestServices: readonly(requestServices),
  requestRoutes: readonly(requestRoutes),
  clientActions: readonly(clientActions),
  clientActionCategories: readonly(clientActionCategories),
};
export default function useFluxState(): ReadonlyFluxState {
  return state;
}
