<template lang="pug">
.configuration
  el-tabs(
    v-model="active"
  )
    el-tab-pane(
      label="所有配置"
      name="all"
    )
      config-table(
        category="*"
        :hiddenHeader="true"
        :hiddenOp="true"
        name="所有配置"
      )
    el-tab-pane(
      label="当前生效配置"
      name="currentValid"
    )
      pre {{currentValid.data}}

</template>
<script lang="ts">
import { defineComponent } from "vue";

import { useConfigStore } from "../../store";
import ConfigTable from "../../components/configs/Table.vue";

export default defineComponent({
  name: "Configuration",
  components: {
    ConfigTable,
  },
  data() {
    return {
      active: "all",
    };
  },
  methods: {
    async fetch() {
      try {
        await this.getCurrentValid();
      } catch (err) {
        this.$error(err);
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
.configuration
  margin: $mainMargin
  padding: 20px 30px
  background-color: $white
pre
  margin: 20px
</style>
