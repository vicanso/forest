<template lang="pug">
//- 路由选择
el-col(
  :span="12"
): el-form-item(
  label="路由选择："
): router-selector.selector(
  :router="router"
  @change.self="handleChangeRouter"
) 
//- 最大并发
el-col(
  :span="12"
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

import RouterSelector from "../../components/configs/RouterSelector.vue";

export default defineComponent({
  name: "RouterConcurrencyData",
  components: {
    RouterSelector,
  },
  props: {
    data: String,
  },
  data() {
    const data = {
      router: "",
      method: "",
      route: "",
      max: null,
    };
    if (this.$props.data) {
      Object.assign(data, JSON.parse(this.$props.data));
      data.router = `${data.method} ${data.route}`;
    }
    return data;
  },
  watch: {
    max() {
      this.handleChange();
    },
  },
  methods: {
    handleChangeRouter(value) {
      this.router = value;
      this.handleChange();
    },
    handleChange() {
      const { router, max } = this;
      let value = "";
      if (router) {
        const [method, route] = router.split(" ");
        value = JSON.stringify({
          method,
          route,
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
  width: 100%
</style>
