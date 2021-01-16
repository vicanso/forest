<template lang="pug">
el-card.validConfiguration(
  v-loading="currentValid.processing"
)
  template(
    #header
  )
    i.el-icon-s-tools
    span 当前生效配置
  pre {{currentValid.data}}

</template>
<script lang="ts">
import { defineComponent } from "vue";

import { useConfigStore } from "../../store";

export default defineComponent({
  name: "ValidConfiguration",
  methods: {
    async fetch() {
      try {
        await this.getCurrentValid();
      } catch (err) {
        this.$error(err.message);
      }
    },
  },
  beforeMount() {
    this.fetch();
  },
  setup() {
    const configStore = useConfigStore();
    return {
      currentValid: configStore.state.currentValid,
      getCurrentValid: () => configStore.dispatch("getCurrentValid"),
    };
  },
});
</script>
<style lang="stylus" scoped>
@import "../../common";
.validConfiguration
  margin: $mainMargin
</style>
