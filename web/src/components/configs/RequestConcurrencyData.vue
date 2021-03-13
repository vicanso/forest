<template lang="pug">
el-col(
  :span="8"
)
//- 实例名称
el-col(
  :span="8"
): el-form-item(
  label="实例："
): el-input(
  placeholder="请输入实例名称"
  v-model="name"
)
//- 最大并发数
el-col(
  :span="8"
): el-form-item(
  label="最大并发："
): el-input(
  type="number"
  placeholder="请输入最大并发限制"
  v-model="max"
)
</template>
<script lang="ts">
import { defineComponent } from "vue";

export default defineComponent({
  name: "RequestConcurrencyData",
  props: {
    data: {
      type: String,
      default: "",
    },
  },
  emits: ["change"],
  data() {
    const data = {
      max: 0,
      name: "",
    };
    if (this.$props.data) {
      Object.assign(data, JSON.parse(this.$props.data));
    }
    return data;
  },
  watch: {
    max() {
      this.handleChange();
    },
    name() {
      this.handleChange();
    },
  },
  methods: {
    handleChange() {
      const { name, max } = this;
      let value = "";
      if (name) {
        value = JSON.stringify({
          name,
          max: Number(max || "0"),
        });
      }
      this.$emit("change", value);
    },
  },
});
</script>
