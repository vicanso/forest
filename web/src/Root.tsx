import { defineComponent } from "vue";
import {
  darkTheme,
  NConfigProvider,
  NDialogProvider,
  NGlobalStyle,
  NLoadingBarProvider,
  NMessageProvider,
  NNotificationProvider,
  zhCN,
} from "naive-ui";

import App from "./App";
import { storeToRefs } from "pinia";
import { useCommonStore } from "./stores/common";

export default defineComponent({
  name: "RootPage",
  setup() {
    const { setting } = storeToRefs(useCommonStore());
    return {
      setting,
    };
  },
  render() {
    const isDark = this.setting.theme === "dark";
    return (
      <NConfigProvider theme={isDark ? darkTheme : null} locale={zhCN}>
        <NLoadingBarProvider>
          <NMessageProvider>
            <NNotificationProvider>
              <NDialogProvider>
                <App />
              </NDialogProvider>
            </NNotificationProvider>
          </NMessageProvider>
        </NLoadingBarProvider>
        <NGlobalStyle />
      </NConfigProvider>
    );
  },
});
