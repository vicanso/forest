import { ROUTERS } from "../urls";
import request from "../request";

export async function list(params) {
  const { data } = await request.get(ROUTERS, {
    params
  });
  return data;
}
