<template lang="pug">
el-col(
  :span="8"
)
//- 实例名称
el-col(
  :span="8"
): el-form-item(
  label="实例选择："
): HTTPInstanceSelector.selector(
  :instance="name"
  @change.self="handleChangeInstance"
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

import HTTPInstanceSelector from "../../components/configs/HTTPInstanceSelector.vue";

export default defineComponent({
  name: "RequestConcurrencyData",
  components: {
    HTTPInstanceSelector,
  },
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
    handleChangeInstance(value) {
      this.name = value;
      this.handleChange();
    },
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

<style lang="stylus" scoped>
.selector
  width 100%
</style>
