import { NLayout, NLayoutSider, useLoadingBar } from "naive-ui";
import { css } from "@linaria/core";
import { defineComponent, onMounted } from "vue";
import AppHeader from "./AppHeader";
import AppNavigation from "./AppNavigation";
import {
  mainHeaderHeight,
  mainNavigationWidth,
  padding,
} from "./constants/style";
import "./main.css";
import { setLoadingEvent } from "./routes/router";
import { storeToRefs } from "pinia";
import { useCommonStore } from "./stores/common";

const layoutClass = css`
  top: ${mainHeaderHeight}px !important;
`;

const contentLayoutClass = css`
  padding: ${2 * padding}px;
`;

export default defineComponent({
  name: "App",
  setup() {
    const commonStore = useCommonStore();
    const { setting } = storeToRefs(commonStore);
    const loadingBar = useLoadingBar();
    if (loadingBar != undefined) {
      setLoadingEvent(loadingBar.start, loadingBar.finish);
      onMounted(() => {
        loadingBar.finish();
      });
    }
    return {
      setting,
      setCollapsed: commonStore.setCollapsed,
    };
  },
  render() {
    const { setting, setCollapsed } = this;
    return (
      <div>
        <AppHeader />
        <NLayout hasSider position="absolute" class={layoutClass}>
          <NLayoutSider
            bordered
            collapseMode="width"
            collapsed={setting.collapsed}
            collapsedWidth={64}
            width={mainNavigationWidth}
            showTrigger
            onCollapse={() => {
              setCollapsed(true);
            }}
            onExpand={() => {
              setCollapsed(false);
            }}
          >
            <AppNavigation />
          </NLayoutSider>
          <NLayout class={contentLayoutClass}>
            <router-view />
          </NLayout>
        </NLayout>
      </div>
    );
  },
});
