<template lang="pug">
el-select(
  v-model="currentRouter"
  placeholder="请选择路由"
  v-loading="routers.processing"
  @change="handleChange"
): el-option(
  v-for="item in routers.items"
  :key="`${item.method} ${item.route}`"
  :label="`${item.method} ${item.route}`"
  :value="`${item.method} ${item.route}`"
)
</template>
<script lang="ts">
import { defineComponent } from "vue";

import { useCommonStore } from "../../store";

export default defineComponent({
  name: "RouterSelector",
  emits: ["change"],
  props: {
    router: String,
  },
  data() {
    return {
      currentRouter: this.$props.router || "",
    };
  },
  methods: {
    handleChange(value) {
      this.$emit("change", value);
    },
    async fetch() {
      try {
        await this.listRouter();
      } catch (err) {
        this.$error(err);
      }
    },
  },
  beforeMount() {
    this.fetch();
  },
  setup() {
    const commonStore = useCommonStore();
    return {
      routers: commonStore.state.routers,
      listRouter: () => commonStore.dispatch("listRouter"),
    };
  },
});
</script>
