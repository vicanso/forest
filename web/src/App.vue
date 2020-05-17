<template>
  <div id="app">
    <MainHeader class="header" />
    <MainNav class="nav" />
    <router-view />
  </div>
</template>
<script>
import { mapActions } from "vuex";
import MainHeader from "@/components/MainHeader.vue";
import MainNav from "@/components/MainNav.vue";

export default {
  name: "App",
  components: {
    MainHeader,
    MainNav
  },
  methods: {
    ...mapActions(["fetchUserInfo"])
  },
  mounted() {
    this.fetchUserInfo();
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
