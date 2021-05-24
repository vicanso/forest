import request from "../helpers/request";

import { ADMINS_CACHE_ID } from "../constants/url";

interface Session {
  data: string;
}

// adminFindCacheByKey 查询缓存
export async function adminFindCacheByKey(key: string): Promise<Session> {
  const url = ADMINS_CACHE_ID.replace(":key", key);
  const { data } = await request.get(url);
  return <Session>data;
}

// adminCleanCacheByKey 清除缓存
export async function adminCleanCacheByKey(key: string): Promise<void> {
  const url = ADMINS_CACHE_ID.replace(":key", key);
  await request.delete(url);
  return;
}
