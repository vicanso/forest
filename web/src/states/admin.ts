import request from "../helpers/request";

import { ADMINS_SESSION_ID } from "../constants/url";

interface Session {
  data: string;
}

// adminFindSessionByID 查询session
export async function adminFindSessionByID(id: string): Promise<Session> {
  const url = ADMINS_SESSION_ID.replace(":id", id);
  const { data } = await request.get(url);
  return <Session>data;
}

// adminCleanSessionByID 清除session
export async function adminCleanSessionByID(id: string): void {
  const url = ADMINS_SESSION_ID.replace(":id", id);
  await request.delete(url);
}
