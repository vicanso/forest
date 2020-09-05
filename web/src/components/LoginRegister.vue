<template>
  <div class="loginRegister">
    <el-card>
      <div slot="header" class="clearfix">
        <i class="el-icon-user-solid" />
        <span>{{ title }}</span>
      </div>
      <el-form
        v-loading="processing"
        ref="form"
        :model="form"
        label-width="80px"
      >
        <el-form-item label="账号：">
          <el-input
            placeholder="请输入账号"
            v-model="form.account"
            autofocus="true"
            clearable
          />
        </el-form-item>
        <el-form-item label="密码：">
          <el-input
            v-model="form.password"
            show-password
            placeholder="请输入密码"
          />
        </el-form-item>
        <el-form-item label="验证码：">
          <el-row>
            <el-col :span="18">
              <el-input
                class="code"
                v-model="form.captcha"
                maxlength="4"
                clearable
                @keyup.enter.native="onSubmit"
                placeholder="请输入验证码"
              />
            </el-col>
            <el-col :span="6">
              <div class="captcha" @click="refreshCaptcha">
                <img
                  v-if="captchaData"
                  :src="`data:image/jpeg;base64,${captchaData.data}`"
                />
              </div>
            </el-col>
          </el-row>
        </el-form-item>
        <el-form-item>
          <el-button class="submit" type="primary" @click="onSubmit"
            >{{ submitText }}
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>
<script>
import { mapState, mapActions } from "vuex";
import { LOGIN } from "@/constants/route";

const registerType = "register";

export default {
  name: "LoginRegister",
  props: {
    type: String
  },
  data() {
    const { type } = this.$props;
    let title = "用户登录";
    let submitText = "立即登录";
    if (type === registerType) {
      title = "用户注册";
      submitText = "立即注册";
    }
    return {
      title,
      submitText,
      captchaData: null,
      form: {
        account: "",
        password: "",
        captcha: ""
      }
    };
  },
  computed: mapState({
    processing: state => state.user.processing
  }),
  methods: {
    ...mapActions(["login", "register", "getCaptcha"]),
    async refreshCaptcha() {
      try {
        const data = await this.getCaptcha();
        this.captchaData = data;
      } catch (err) {
        this.$message.error(err.message);
      }
    },
    async onSubmit() {
      const { account, password, captcha } = this.form;
      if (!account || !password || !captcha) {
        this.$message.warning("账号、密码以及验证码不能为空");
        return;
      }
      if (this.submitting) {
        return;
      }
      this.submitting = true;
      const params = {
        account,
        password,
        captcha: `${this.captchaData.id}:${captcha}`
      };
      try {
        if (this.$props.type === registerType) {
          await this.register(params);
          this.$router.replace({
            name: LOGIN
          });
        } else {
          await this.login(params);
          this.$router.back();
        }
      } catch (err) {
        // 图形验证码只可校验一次，因此出错则刷新
        this.refreshCaptcha();
        this.$message.error(err.message);
      } finally {
        this.submitting = false;
      }
    }
  },
  mounted() {
    this.refreshCaptcha();
  }
};
</script>
<style lang="sass" scoped>
@import "@/common.sass"
.loginRegister
  margin: 100px auto
  max-width: 600px
  i
    margin-right: 5px
.captcha
  cursor: pointer
  overflow: hidden
  height: 40px
  text-align: center
.code
  width: 100%
.submit
  width: 100%
</style>
