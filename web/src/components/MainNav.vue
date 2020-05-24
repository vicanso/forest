<template>
  <div class="mainNav">
    <h1>
      <router-link :to="{ name: home }">
        <i class="el-icon-eleme" />
        Forest
      </router-link>
    </h1>
    <nav>
      <el-menu
        class="menu"
        :default-active="active"
        background-color="#000c17"
        text-color="#fff"
        active-text-color="#fff"
      >
        <el-submenu
          class="submenu"
          v-for="(nav, i) in navs"
          :index="`${i}`"
          :key="`${i}`"
        >
          <template slot="title">
            <i :class="nav.icon" />
            <span>{{ nav.name }}</span>
          </template>
          <el-menu-item
            class="menuItem"
            v-for="(subItem, j) in nav.children"
            :index="`${i}-${j}`"
            :key="`${i}-${j}`"
            @click="goTo(subItem)"
          >
            <span>{{ subItem.name }}</span>
          </el-menu-item>
        </el-submenu>
      </el-menu>
    </nav>
  </div>
</template>
<script>
import {
  HOME,
  CONFIG_MOCK_TIME,
  CONFIG_BLOCK_IP,
  CONFIG_SIGNED_KEY,
  CONFIG_ROUTER,
  CONFIG_ROUTER_CONCURRENCY,
  USERS,
  LOGINS
} from "@/constants/route";
import { USER_ADMIN, USER_SU } from "@/constants/user";
import { mapState } from "vuex";
import { isAllowedUser } from "@/helpers/util";

const navs = [
  {
    name: "配置",
    icon: "el-icon-setting",
    roles: [USER_SU],
    children: [
      {
        name: "MockTime配置",
        route: CONFIG_MOCK_TIME
      },
      {
        name: "黑名单IP",
        route: CONFIG_BLOCK_IP
      },
      {
        name: "SignedKey",
        route: CONFIG_SIGNED_KEY
      },
      {
        name: "路由配置",
        route: CONFIG_ROUTER
      },
      {
        name: "路由并发配置",
        route: CONFIG_ROUTER_CONCURRENCY
      }
    ]
  },
  {
    name: "用户",
    icon: "el-icon-user",
    roles: [USER_ADMIN, USER_SU],
    children: [
      {
        name: "用户列表",
        route: USERS
      },
      {
        name: "登录记录",
        route: LOGINS
      }
    ]
  }
];

export default {
  name: "MainNav",
  data() {
    return {
      home: HOME,
      active: ""
    };
  },
  computed: {
    ...mapState({
      userInfo: state => state.user.info
    }),
    navs() {
      const { userInfo } = this;
      if (!userInfo || !userInfo.account) {
        return [];
      }
      const { roles } = userInfo;
      const filterNavs = [];
      navs.forEach(item => {
        // 如果该栏目有配置权限，而且用户无该权限
        if (item.roles && !isAllowedUser(item.roles, roles)) {
          return;
        }
        filterNavs.push(item);
      });
      return filterNavs;
    }
  },
  watch: {
    // 路由变化时设置对应的导航为活动状态
    $route(to) {
      const { navs } = this;
      let active = "";
      navs.forEach((nav, i) => {
        nav.children.forEach((item, j) => {
          if (item.route === to.name) {
            active = `${i}-${j}`;
          }
        });
      });
      this.active = active;
    }
  },
  methods: {
    goTo({ route }) {
      if (!route || this.$route.name === route) {
        return;
      }
      this.$router.push({
        name: route
      });
    }
  }
};
</script>
<style lang="sass" scoped>
@import "@/common.sass"
$mainNavColor: #000c17
.mainNav
  min-height: 100vh
  overflow-y: auto
  background-color: $mainNavColor
h1
  height: $mainHeaderHeight
  line-height: $mainHeaderHeight
  color: $white
  padding-left: 20px
  font-size: 18px
  i
    font-weight: bold
nav
  border-top: 1px solid rgba($white, 0.3)
.menu
  border-right: 1px solid $mainNavColor
.menuItem
  color: rgba($white, 0.65)
  &.is-active
    background-color: $darkBlue !important
</style>
