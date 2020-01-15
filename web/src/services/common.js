import { RANDOM_KEYS, CAPTCHA } from "../urls";
import request from "../request";

export async function getRandomKeys(params) {
  const { data } = await request.get(RANDOM_KEYS, {
    params
  });
  return data;
}

export async function getCaptcha() {
  const { data } = await request.get(CAPTCHA);
  return data;
}
