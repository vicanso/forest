<template>
  <div class="routerConcurrencyData">
    <el-col :span="12">
      <el-form-item label="路由选择：">
        <RouterSelector
          class="selector"
          :router="router"
          @change="handleChangeRouter"
        />
      </el-form-item>
    </el-col>
    <el-col :span="12">
      <el-form-item label="最大并发：">
        <el-input
          type="nubmer"
          placeholder="请输入最大并发限制"
          v-model="max"
        />
      </el-form-item>
    </el-col>
  </div>
</template>
<script>
import RouterSelector from "@/components/RouterSelector.vue";

export default {
  name: "RouterConcurrencyData",
  components: {
    RouterSelector
  },
  props: {
    data: String
  },
  data() {
    const data = {
      router: "",
      max: null
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
    }
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
          max: Number(max || "0")
        });
      }
      this.$emit("change", value);
    }
  }
};
</script>
<style lang="sass" scoped>
.selector
  width: 100%
</style>
