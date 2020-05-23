<template>
  <el-card class="profile" v-loading="processing">
    <div slot="header">
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
            {{ profile.rolesDesc.join(",") || "--" }}
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
        <el-col :span="12">
          <el-form-item label="手机：">
            <el-input placeholder="请输入手机号码" v-model="mobile" clearable />
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item label="邮箱：">
            <el-input placeholder="请输入邮箱地址" v-model="email" clearable />
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
export default {
  name: "Profile",
  data() {
    return {
      email: "",
      mobile: ""
    };
  },
  computed: mapState({
    profile: state => state.user.profile,
    processing: state => state.user.profileProcessing
  }),
  methods: {
    ...mapActions(["getUserProfile", "updateMe"]),
    async onSubmit() {
      const { email, mobile, profile } = this;
      if ((profile.mobile && !mobile) || (profile.email && !email)) {
        this.$message.warning("手机号码与邮箱不能删除");
        return;
      }
      const update = {};
      if (profile.mobile != mobile) {
        update.mobile = mobile;
      }
      if (profile.email != email) {
        update.email = email;
      }
      if (Object.keys(update) === 0) {
        this.$message.warning("请修改信息后再更新");
        return;
      }
      try {
        await this.updateMe(update);
        this.$message.info("信息已成功更新");
      } catch (err) {
        this.$message.error(err.message);
      }
    },
    goBack() {
      this.$router.back();
    }
  },
  async beforeMount() {
    try {
      await this.getUserProfile();
      const { email, mobile } = this.profile;
      this.email = email;
      this.mobile = mobile;
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
.btn
  width: 100%
</style>
