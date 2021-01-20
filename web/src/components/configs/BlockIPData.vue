<template lang="pug">
el-col(
  :span="8"
): el-form-item.blockIPData(
  label="拦截IP："
): el-input(
  placeholder="请输入IP地址或网段"
  v-model="ip"
  clearable
)
</template>
<script lang="ts">
import { defineComponent } from "vue";

export default defineComponent({
  name: "BlockIPData",
  emits: ["change"],
  props: {
    data: String,
  },
  data() {
    return {
      ip: this.$props.data || "",
    };
  },
  watch: {
    ip: function (value) {
      this.$emit("change", `${value}`);
    },
  },
});
</script>
