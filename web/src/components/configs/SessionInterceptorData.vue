<template>
  <el-col :span="8">
    <el-form-item label="提示信息：" class="sessionInterceptorData">
      <el-input placeholder="请输入提示信息" v-model="message" />
    </el-form-item>
  </el-col>
</template>
<script>
// session 拦截的提示信息
export default {
  name: "SessionInterceptorData",
  props: {
    data: String
  },
  data() {
    return {
      message: this.$props.data || ""
    };
  },
  watch: {
    message: function(value) {
      this.$emit("change", value);
    }
  }
};
</script>
