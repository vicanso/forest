<template lang="pug">
el-select(
  v-model="currentInstance"
  placeholder="请选择HTTP实例"
  v-loading="instances.processing"
  @change="handleChange"
): el-option(
  v-for="item in instances.items"
  :key="item.name"
  :label="getLabel(item)"
  :value="item.name"
)
</template>

<script lang="ts">
import { defineComponent } from "vue";

import useCommonState, { commonListHTTPInstance } from "../../states/common";

export default defineComponent({
  name: "HTTPInstanceSelector",
  props: {
    instance: {
      type: String,
      default: "",
    },
  },
  emits: ["change"],
  setup() {
    const commonState = useCommonState();
    return {
      instances: commonState.httpInstances,
    };
  },
  data() {
    return {
      currentInstance: this.$props.instance || "",
    };
  },
  beforeMount() {
    this.fetch();
  },
  methods: {
    getLabel(item) {
      const maxConcurrency = item.maxConcurrency || 0;
      const concurrency = item.concurrency || 0;
      return `${item.name} (max:${maxConcurrency}, current:${concurrency})`;
    },
    handleChange(value) {
      this.$emit("change", value);
    },
    async fetch() {
      try {
        await commonListHTTPInstance();
      } catch (err) {
        this.$error(err);
      }
    },
  },
});
</script>
