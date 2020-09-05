<template>
  <el-select
    v-model="currentRouter"
    placeholder="请选择路由"
    v-loading="processing"
    @change="handleChange"
  >
    <el-option
      v-for="item in routers"
      :key="item.key"
      :label="item.key"
      :value="item.key"
    >
    </el-option>
  </el-select>
</template>
<script>
import { mapState, mapActions } from "vuex";

export default {
  name: "RouterSelector",
  props: {
    router: String
  },
  data() {
    return {
      currentRouter: this.$props.router || ""
    };
  },
  computed: mapState({
    processing: state => state.common.processing,
    routers: state => state.common.routers || []
  }),
  methods: {
    ...mapActions(["listRouter"]),
    handleChange(value) {
      this.$emit("change", value);
    }
  },
  async beforeMount() {
    try {
      await this.listRouter();
    } catch (err) {
      this.$message.error(err.message);
    }
  }
};
</script>
