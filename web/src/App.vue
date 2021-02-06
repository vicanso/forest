<template lang="pug">
#app(
  :class="{ shrinking: setting.mainNavShrinking }"
  v-loading="loadingSetting"
)
  //- 主头部
  main-header.header(
    v-if="!loadingSetting"
  )
  //- 主导航
  main-nav.nav(
    :shrinking="setting.mainNavShrinking"
    @toggle="toggleNav"
    v-if="!loadingSetting"
  )
  //- 内容区域
  .mainContent
    router-view(
      v-if="inited"
    )
    p.tac(
      v-else
    ) ...

</template>

<script lang="ts">
import { defineComponent } from "vue";
import MainHeader from "./components/MainHeader.vue";
import MainNav from "./components/MainNav.vue";

import useUserState, { userFetchInfo, userUpdate } from "./states/user";
import { ROUTE_LOGIN } from "./router";
import useSettingState, { settingLoad, settingSave } from "./states/setting";

export default defineComponent({
  name: "App",
  components: {
    MainHeader,
    MainNav,
  },
  setup() {
    const userState = useUserState();
    const settingState = useSettingState();
    return {
      setting: settingState,
      userInfo: userState.info,
    };
  },
  data() {
    return {
      loadingSetting: false,
      // 是否初始化完成
      inited: false,
    };
  },
  async beforeMount() {
    this.loadingSetting = true;
    try {
      await settingLoad();
    } catch (err) {
      this.$error(err);
    } finally {
      this.loadingSetting = false;
    }
  },
  mounted() {
    this.fetch();
  },
  methods: {
    toggleNav() {
      settingSave({
        mainNavShrinking: !this.setting.mainNavShrinking,
      });
    },
    async fetch() {
      const { userInfo, $router } = this;
      try {
        await userFetchInfo();
        // 如果未登录则跳转至登录
        if (!userInfo.account) {
          $router.push({
            name: ROUTE_LOGIN,
          });
        } else {
          // 如果已登录，刷新cookie有效期（不关注刷新是否成功，因此不用await）
          userUpdate({});
        }
      } catch (err) {
        this.$error(err);
      } finally {
        this.inited = true;
      }
    },
  },
});
</script>

<style lang="stylus" scoped>
@import "./common";
.shrinking
  .header
    left $mainNavShrinkingWidth
  .nav
    width $mainNavShrinkingWidth
  .mainContent
    padding-left $mainNavShrinkingWidth
.header
  position fixed
  left $mainNavWidth
  top 0
  right 0
  z-index 9
.nav
  position fixed
  width $mainNavWidth
  top 0
  bottom 0
  left 0
  overflow hidden
  overflow-y auto
.mainContent
  padding-left $mainNavWidth
  padding-top $mainHeaderHeight
</style>
