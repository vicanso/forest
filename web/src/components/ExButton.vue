<template lang="pug">
.exButton(
  @click.preventDefault="handleClick"
  :class="{ isProcessing: processing}"
)
  slot
</template>
<script lang="ts">
// 此button的扩展可记录用户行为，防止重复点击等
import { defineComponent } from "vue";

import { getCurrentLocation } from "../router";
import { addUserAction } from "../services/action";

export default defineComponent({
  name: "ExButton",
  props: {
    onClick: {
      type: Function,
      required: true,
    },
    category: {
      type: String,
      required: true,
    },
    extra: {
      type: Object,
      default: function () {
        return {};
      },
    },
  },
  data() {
    return {
      // 是否处理中，避免重复点击
      processing: false,
    };
  },
  methods: {
    async handleClick() {
      if (this.processing) {
        return;
      }
      const currentLocation = getCurrentLocation();
      const data = {
        category: this.$props.category,
        route: currentLocation.name,
        path: currentLocation.path,
        time: Math.floor(Date.now() / 1000),
        extra: this.$props.extra,
      };
      this.processing = true;
      // 由于在onClick会捕获异常处理，因此在此处无法判断是否成功
      try {
        await this.$props.onClick();
      } finally {
        this.processing = false;
      }
      addUserAction(data);
    },
  },
});
</script>
<style lang="stylus" scoped>
.exButton.isProcessing
  opacity: 0.5
</style>
