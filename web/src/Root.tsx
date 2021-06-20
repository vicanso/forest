import { defineComponent } from "vue";
import {
  darkTheme,
  NConfigProvider,
  NDialogProvider,
  NGlobalStyle,
  NLoadingBarProvider,
  NMessageProvider,
  NNotificationProvider,
} from "naive-ui";
import useCommonState from "./states/common";

import App from "./App";

export default defineComponent({
  name: "Root",
  setup() {
    const { settings } = useCommonState();
    return {
      settings,
    };
  },
  render() {
    const isDark = this.settings.theme === "dark";
    return (
      <NConfigProvider theme={isDark ? darkTheme : null}>
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
