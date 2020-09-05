<template>
  <div class="mainNav">
    <h1>
      <router-link :to="{ name: home }">
        <i class="el-icon-eleme" />
        Origin
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
  CONFIG_ORDER_COMMISSION,
  CONFIG_MARKETING_GROUP,
  USERS,
  LOGINS,
  BRANDS,
  PRODUCTS,
  PRODUCT_CATEGORIES,
  SUPPLIERS,
  REGIONS,
  ADVERTISEMENTS,
  ORDERS
} from "@/constants/route";
import { USER_ADMIN, USER_SU, GROUP_MARKETING } from "@/constants/user";
import { mapState } from "vuex";
import { isAllowedUser } from "@/helpers/util";

const navs = [
  {
    name: "业务",
    icon: "el-icon-files",
    groups: [GROUP_MARKETING],
    children: [
      {
        name: "订单",
        route: ORDERS
      },
      {
        name: "品牌",
        route: BRANDS
      },
      {
        name: "产品",
        route: PRODUCTS
      },
      {
        name: "产品分类",
        route: PRODUCT_CATEGORIES
      },
      {
        name: "供应商",
        route: SUPPLIERS
      },
      {
        name: "广告",
        route: ADVERTISEMENTS
      },
      {
        name: "地区",
        route: REGIONS
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
  },
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
      },

      {
        name: "销售分组配置",
        route: CONFIG_MARKETING_GROUP
      },
      {
        name: "订单佣金配置",
        route: CONFIG_ORDER_COMMISSION
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
      const { roles, groups } = userInfo;
      const filterNavs = [];
      navs.forEach(item => {
        // 如果该栏目有配置权限，而且用户无该权限
        if (item.roles && !isAllowedUser(item.roles, roles)) {
          return;
        }
        // 如果该栏目配置了允许分级，而该用户不属于该组
        if (item.groups && !isAllowedUser(item.groups, groups)) {
          return;
        }
        const clone = Object.assign({}, item);
        const children = item.children.map(subItem =>
          Object.assign({}, subItem)
        );
        clone.children = children.filter(subItem => {
          // 如果未配置色色与分组限制
          if (!subItem.roles && !subItem.groups) {
            return true;
          }
          if (subItem.roles && !isAllowedUser(subItem.roles, roles)) {
            return false;
          }
          if (subItem.groups && !isAllowedUser(subItem.groups, groups)) {
            return false;
          }
          return true;
        });
        filterNavs.push(clone);
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
