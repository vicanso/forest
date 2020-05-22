<template>
  <div id="app">
    <MainHeader class="header" />
    <MainNav class="nav" />
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
  computed: mapState({
    userAccount: state => state.user.info.account
  }),
  methods: {
    ...mapActions(["fetchUserInfo", "updateUser"]),
    refreshSessionTTL() {
      if (!this.userAccount) {
        return;
      }
      this.updateUser({});
    }
  },
  async mounted() {
    setInterval(() => {
      this.refreshSessionTTL();
    }, 5 * 60 * 1000);
    try {
      await this.fetchUserInfo();
    } catch (err) {
      this.$message.error(err.message);
    }
  }
};
</script>
<style lang="sass" scoped>
@import "@/common.sass"
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
.mainContent
  padding-left: $mainNavWidth
  padding-top: $mainHeaderHeight
</style>
