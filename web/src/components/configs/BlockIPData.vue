<template>
  <el-col :span="8">
    <el-form-item label="拦截IP：" class="blockIPData">
      <el-input placeholder="请输入IP地址或网段" v-model="ip" />
    </el-form-item>
  </el-col>
</template>
<script>
// 黑名单IP拦截数据组件
export default {
  name: "BlockIPData",
  props: {
    data: String
  },
  data() {
    return {
      ip: this.$props.data || ""
    };
  },
  watch: {
    ip: function(value) {
      this.$emit("change", value);
    }
  }
};
</script>
