<template>
  <div id="app" :class="{ shrinking: shrinking }">
    <MainHeader class="header" />
    <MainNav class="nav" :onToggle="toggleNav" :shrinking="shrinking" />
    <div class="mainContent">
      <router-view />
    </div>
  </div>
</template>
<script>
import { mapActions, mapState } from "vuex";
import MainHeader from "@/components/MainHeader.vue";
import MainNav from "@/components/MainNav.vue";

export default {
  name: "App",
  components: {
    MainHeader,
    MainNav
  },
  data() {
    return {
      shrinking: false
    };
  },
  computed: mapState({
    userAccount: state => state.user.info.account
  }),
  methods: {
    ...mapActions(["fetchUserInfo", "updateMe", "listStatus"]),
    refreshSessionTTL() {
      if (!this.userAccount) {
        return;
      }
      this.updateMe({});
    },
    toggleNav() {
      this.shrinking = !this.shrinking;
    }
  },
  async mounted() {
    setInterval(() => {
      this.refreshSessionTTL();
    }, 5 * 60 * 1000);
    try {
      await this.listStatus();
      await this.fetchUserInfo();
    } catch (err) {
      this.$message.error(err.message);
    }
  }
};
</script>
<style lang="sass" scoped>
@import "@/common.sass"
.shrinking
  .header
    left: $mainNavShrinkingWidth
  .nav
    width: $mainNavShrinkingWidth
  .mainContent
    padding-left: $mainNavShrinkingWidth
.header
  position: fixed
  left: $mainNavWidth
  top: 0
  right: 0
  z-index: 9
.nav
  position: fixed
  width: $mainNavWidth
  top: 0
  bottom: 0
  left: 0
  overflow: hidden
  overflow-y: auto
.mainContent
  padding-left: $mainNavWidth
  padding-top: $mainHeaderHeight
</style>
