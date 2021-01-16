import axios from "axios";

import { isDevelopment } from "../constants/env";

const request = axios.create({
  // 默认超时为10秒
  timeout: 10 * 1000,
});

request.interceptors.request.use(
  (config) => {
    if (isDevelopment()) {
      config.url = `/api${config.url}`;
    }
    return config;
  },
  (err) => {
    return Promise.reject(err);
  }
);
request.interceptors.response.use(
  (res) => {
    return res;
  },
  (err) => {
    const { response } = err;
    if (err.code === "ECONNABORTED") {
      err.message = "请求超时，请稍候再试";
    } else if (response) {
      if (response.data && response.data.message) {
        err.message = response.data.message;
        err.code = response.data.code;
        err.category = response.data.category;
      } else {
        err.message = `unknown error[${response.statusCode || -1}]`;
      }
    }
    return Promise.reject(err);
  }
);

export default request;
