<template lang="pug">
//- 用户角色
mixin Roles
  el-col(
    :span="10"
  ): el-form-item(
    label="用户角色："
  ): span {{roles}}

//- 邮箱
mixin Email
  el-col(
    :span="10"
  ): el-form-item(
    label="用户邮箱："
  ): el-input(
    placeholder="请输入您的邮箱地址："
    v-model="email"
    clearable
  )

//- 原密码
mixin OriginalPassword
  el-col(
    :span="10"
  ): el-form-item(
    label="用户原密码："
  ): el-input(
    type="password"
    :disabled="!enableUpdatePassword"
    placeholder="请输入您的原密码："
    v-model="password"
    clearable
  )

//- 新密码
mixin NewPassword
  el-col(
    :span="10"
  ): el-form-item(
    label="用户新密码："
  ): el-input(
    type="password"
    :disabled="!enableUpdatePassword"
    placeholder="请输入您的新密码："
    v-model="newPassword"
    clearable
  )

//- 启用更新密码
mixin EnableUpdatePassword
  el-col(
    :span="4"
  ): el-form-item(
    label="更新密码："
  )
    el-checkbox(
      v-model="enableUpdatePassword"
    )

//- 提交与返回
mixin SubmitAndBack
  el-col(
    :span="12"
  ): el-form-item: ex-button(
    :onClick="submit"
  ) 更新 
  el-col(
    :span="12"
  ): el-form-item: el-button.btn(
    @click="goBack"
  ) 返回
el-card.profile
  template(
    #header
  )
    i.el-icon-user-solid
    span 用户信息
  el-form(
    label-width="120px"
    v-loading="processing"
  ): el-row(
    :gutter="15"
  )
    //- 用户角色
    +Roles 

    //- 邮箱地址
    +Email


    //- 原密码
    +OriginalPassword

    //- 新密码
    +NewPassword
    
    //- 是否启用更新密码
    +EnableUpdatePassword

    //- 提交与返回
    +SubmitAndBack

</template>

<script lang="ts">
import { defineComponent } from "vue";

import ExButton from "../components/ExButton.vue";
import { userUpdate, userFetchDetail, userLogout } from "../states/user";
import { ROUTE_LOGIN } from "../router";

export default defineComponent({
  name: "Profile",
  components: {
    ExButton,
  },
  data() {
    return {
      processing: false,
      roles: [],
      email: "",
      newPassword: "",
      password: "",
      enableUpdatePassword: false,
    };
  },
  beforeMount() {
    this.fetch();
  },
  methods: {
    // fetch 拉取信息
    async fetch() {
      this.processing = true;
      try {
        const data = await userFetchDetail();
        this.email = data.email;
        this.roles = data.roles;
      } catch (err) {
        this.$error(err);
      } finally {
        this.processing = false;
      }
    },
    goBack() {
      this.$router.back();
    },
    // submit 提交
    async submit(): Promise<boolean> {
      let isSuccess = false;
      const {
        email,
        newPassword,
        password,
        enableUpdatePassword,
        processing,
      } = this;
      if (processing) {
        return isSuccess;
      }
      const updateData: Record<string, unknown> = {};
      if (enableUpdatePassword) {
        if (!newPassword || !password) {
          this.$error("原密码与新密码不能为空");
          return isSuccess;
        }
        if (newPassword == password) {
          this.$error("新密码不能与原密码相同");
          return isSuccess;
        }
        updateData.newPassword = newPassword;
        updateData.password = password;
      }
      if (email) {
        updateData.email = email;
      }
      if (Object.keys(updateData).length === 0) {
        this.$error("请修改数据再更新");
        return isSuccess;
      }
      try {
        this.processing = true;
        await userUpdate(updateData);
        if (updateData.newPassword) {
          await userLogout();
          this.$message.info("已成功更新，需要重新登录");
          this.$router.replace({
            name: ROUTE_LOGIN,
          });
        } else {
          this.$message.info("成功更新");
        }
        isSuccess = true;
      } catch (err) {
        this.$error(err);
      } finally {
        this.processing = false;
      }
      return isSuccess;
    },
  },
});
</script>

<style lang="stylus" scoped>
@import "../common";
.profile
  margin $mainMargin
.btn
  width 100%
.pagination
  text-align right
  margin-top 15px
</style>
