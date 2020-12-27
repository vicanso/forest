<template>
  <el-card class="profile" v-loading="processing">
    <div slot="header">
      <a
        href="#"
        @click.prevent="toggleEnableUpdatePassword"
        class="updatePassword"
        >修改密码</a
      >
      <i class="el-icon-user-solid" />
      <span>我的信息</span>
    </div>
    <el-form label-width="80px" v-if="profile">
      <el-row :gutter="15">
        <el-col :span="6">
          <el-form-item label="账号：">
            {{ profile.account }}
          </el-form-item>
        </el-col>
        <el-col :span="6">
          <el-form-item label="角色：">
            {{ profile.roles.join(",") || "--" }}
          </el-form-item>
        </el-col>
        <el-col :span="6">
          <el-form-item label="部门：">
            {{ profile.groupsDesc.join(",") || "--" }}
          </el-form-item>
        </el-col>
        <el-col :span="6">
          <el-form-item label="状态：">
            {{ profile.statusDesc }}
          </el-form-item>
        </el-col>

        <el-col :span="8">
          <el-form-item label="名称：">
            <el-input placeholder="请输入您的名称" v-model="name" clearable />
          </el-form-item>
        </el-col>
        <el-col :span="8">
          <el-form-item label="手机：">
            <el-input placeholder="请输入手机号码" v-model="mobile" clearable />
          </el-form-item>
        </el-col>
        <el-col :span="8">
          <el-form-item label="邮箱：">
            <el-input placeholder="请输入邮箱地址" v-model="email" clearable />
          </el-form-item>
        </el-col>
        <!-- 修改密码 -->
        <el-col :span="12">
          <el-form-item label="旧密码：">
            <el-input
              :disabled="!enableUpdatePassword"
              placeholder="请输入旧密码(需先点击修改密码)"
              v-model="password"
              clearable
              show-password
            />
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item label="新密码：">
            <el-input
              :disabled="!enableUpdatePassword"
              placeholder="请输入新密码(需先点击修改密码)"
              v-model="newPassword"
              show-password
              clearable
            />
          </el-form-item>
        </el-col>

        <el-col :span="12">
          <el-form-item>
            <el-button class="btn" type="primary" @click="onSubmit"
              >更新</el-button
            >
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item>
            <el-button class="btn" @click="goBack">返回</el-button>
          </el-form-item>
        </el-col>
      </el-row>
    </el-form>
  </el-card>
</template>
<script>
import { mapActions, mapState } from "vuex";
import { LOGIN } from "@/constants/route";
export default {
  name: "Profile",
  data() {
    return {
      enableUpdatePassword: false,
      password: "",
      newPassword: "",
      email: "",
      mobile: "",
      name: ""
    };
  },
  computed: mapState({
    profile: state => state.user.info,
    processing: state => state.user.profileProcessing
  }),
  methods: {
    ...mapActions(["fetchUserInfo", "updateMe", "logout"]),
    async onSubmit() {
      const {
        email,
        mobile,
        profile,
        password,
        newPassword,
        name,
        enableUpdatePassword
      } = this;
      if ((profile.mobile && !mobile) || (profile.email && !email)) {
        this.$message.warning("手机号码与邮箱不能删除");
        return;
      }
      const update = {};
      if (enableUpdatePassword) {
        if (!password) {
          this.$message.warning("请输入旧密码");
          return;
        }
        if (newPassword === password) {
          this.$message.warning("新旧密码不能相同");
          return;
        }
        update.password = password;
        update.newPassword = newPassword;
      }
      if (profile.mobile != mobile) {
        update.mobile = mobile;
      }
      if (profile.email != email) {
        update.email = email;
      }
      if (profile.name != name) {
        update.name = name;
      }
      if (Object.keys(update).length === 0) {
        this.$message.warning("请修改信息后再更新");
        return;
      }
      try {
        await this.updateMe(update);
        if (update.newPassword) {
          await this.logout();
          this.$message.info("信息已成功更新，由于更改了密码需要重新登录");
          this.$router.push({
            name: LOGIN
          });
        } else {
          this.$message.info("信息已成功更新");
        }
      } catch (err) {
        this.$message.error(err.message);
      }
    },
    goBack() {
      this.$router.back();
    },
    toggleEnableUpdatePassword() {
      this.enableUpdatePassword = !this.enableUpdatePassword;
    }
  },
  async beforeMount() {
    try {
      await this.fetchUserInfo();
      const { email, mobile, name } = this.profile;
      this.email = email;
      this.mobile = mobile;
      this.name = name;
    } catch (err) {
      this.$message.error(err.message);
    }
  }
};
</script>
<style lang="sass" scoped>
@import "@/common.sass"
.profile
  margin: $mainMargin
  i
    margin-right: 5px
  .updatePassword
    float: right
    font-size: 13px
.btn
  width: 100%
</style>
