<template>
  <el-col :span="8">
    <el-form-item label="SignedKey：" class="signedKeyData">
      <el-input placeholder="请输入签名的字符串以,分隔" v-model="key" />
    </el-form-item>
  </el-col>
</template>
<script>
export default {
  name: "SignedKeyData",
  props: {
    data: String
  },
  data() {
    return {
      key: this.$props.data || ""
    };
  },
  watch: {
    key: function(value) {
      this.$emit("change", value);
    }
  }
};
</script>
