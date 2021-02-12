<template lang="pug">
//- 登录后的用户相关功能
mixin UserFunctions
  router-link(
    :to="{ name: profileRoute }"
  )
    i.el-icon-user
    span {{user.account}}
  span.divided |
  a.logout(
    href="#"
    title="退出登录"
    @click.preventDefault="onLogout"
  )
    i.el-icon-switch-button

//- 登录与注册
mixin LoginAndRegister
  router-link.login(
    :to="{ name: loginRoute }"
  )
    i.el-icon-user
    | 登录
  span.divided |
  router-link.register(
    :to="{ name: registerRoute }"
  )
    i.el-icon-circle-plus
    | 注册

//- 页图标
mixin HomeLogo
  h1(
    v-if="$props.shrinking"
  ): router-link(
    :to='{name: homeRoute}'
  )
    i.el-icon-cpu
    | Forest

header.header
  //- 用户信息
  .userInfo
    span(
      v-if="user.processing"
    ) 正在加载...
    .functions(
      v-else-if="user.account"
    )
      +UserFunctions
    div(
      v-else
    )
      +LoginAndRegister

  +HomeLogo

</template>

<script lang="ts">
import { defineComponent } from "vue";

import useUserState, { userLogout } from "../states/user";
import {
  ROUTE_HOME,
  ROUTE_LOGIN,
  ROUTE_REGISTER,
  ROUTE_PROFILE,
} from "../router";

export default defineComponent({
  name: "MainHeader",
  props: {
    shrinking: {
      type: Boolean,
      default: false,
    },
  },
  setup() {
    const userState = useUserState();
    return {
      homeRoute: ROUTE_HOME,
      profileRoute: ROUTE_PROFILE,
      loginRoute: ROUTE_LOGIN,
      registerRoute: ROUTE_REGISTER,
      user: userState.info,
    };
  },
  methods: {
    async onLogout() {
      try {
        await userLogout();
        this.$router.push({
          name: ROUTE_HOME,
        });
      } catch (err) {
        this.$error(err);
      }
    },
  },
});
</script>

<style lang="stylus" scoped>
@import "../common";
.header
  height $mainHeaderHeight
  background-color $white
  padding 5px 0
  line-height $mainHeaderHeight - 10
  color $darkBlue
  box-shadow 0 1px 4px rgba(0, 21, 41, 0.08)
.userInfo
  float right
  font-size 13px
  margin-right $mainMargin
  i
    margin-right 3px
    font-weight bold
.divided
  margin 0 15px
h1
  font-size 18px
  margin-left 10px
  a
    color $dark
  i
    margin-right 5px
</style>
