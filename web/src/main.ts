import { createApp } from "vue";
import ElementPlus from "element-plus";
import { ElMessage } from "element-plus";
import "element-plus/lib/theme-chalk/index.css";
import App from "./App.vue";
import ExButton from "./components/ExButton.vue";
import Router, { getCurrentLocation } from "./router";
import "./main.styl";
import { actionAdd, ERROR, FAIL } from "./states/action";
import { isDevelopment } from "./constants/env";
import HTTPError from "./helpers/http-error";

const app = createApp(App);
// 全局出错处理
app.config.errorHandler = (err: unknown, vm, info) => {
  // 处理错误
  let message = "未知错误";
  if (err instanceof Error) {
    message = err.message;
  }
  if (info) {
    message += ` [${info}]`;
  }
  const currentLocation = getCurrentLocation();
  actionAdd({
    category: ERROR,
    route: currentLocation.name,
    path: currentLocation.path,
    result: FAIL,
    message,
  });
  throw err;
};
// 自定义全局出错提示
app.config.globalProperties.$error = function (err: Error | HTTPError) {
  let message = err.toString();
  if (err instanceof HTTPError) {
    message = err.message;
    if (err.category) {
      message += ` [${err.category.toUpperCase()}]`;
    }
    if (err.code) {
      message += ` [${err.code}]`;
    }
    // 如果是异常（客户端异常，如请求超时，中断等），则上报user action
    if (err.exception) {
      const currentLocation = getCurrentLocation();
      actionAdd({
        category: ERROR,
        route: currentLocation.name,
        path: currentLocation.path,
        result: FAIL,
        message,
      });
    }
  } else if (err instanceof Error) {
    message =  err.message;
  }
  ElMessage.error(message);
  if (isDevelopment()) {
    console.error(err);
  }
  return;
};

// 全局注册组件
app.component("ExButton", ExButton);

app.use(Router).use(ElementPlus).mount("#app");
