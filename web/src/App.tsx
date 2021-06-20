import { defineComponent, onMounted } from "vue";
import { NLayout, NLayoutSider, useLoadingBar } from "naive-ui";
import { setLoadingEvent } from "./routes/router";

import AppHeader from "./AppHeader";
import AppNavigation from "./AppNavigation";
import "./main.css";
import {
  mainHeaderHeight,
  mainNavigationWidth,
  padding,
} from "./constants/style";

const layoutStyle = {
  top: `${mainHeaderHeight}px`,
};

const contentLayoutStyle = {
  padding: `${2 * padding}px`,
};

export default defineComponent({
  name: "App",
  setup() {
    const loadingBar = useLoadingBar();
    if (loadingBar != undefined) {
      setLoadingEvent(loadingBar.start, loadingBar.finish);
      onMounted(() => {
        loadingBar.finish();
      });
    }
  },
  render() {
    return (
      <div>
        <AppHeader />
        <NLayout hasSider position="absolute" style={layoutStyle}>
          <NLayoutSider
            bordered
            collapseMode="width"
            collapsedWidth={64}
            width={mainNavigationWidth}
            showTrigger
          >
            <AppNavigation />
          </NLayoutSider>
          <NLayout style={contentLayoutStyle}>
            <router-view />
          </NLayout>
        </NLayout>
      </div>
    );
  },
});
