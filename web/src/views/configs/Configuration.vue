<template lang="pug">
//- 所有配置
mixin AllSetting
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

//- 生效配置
mixin ActiveSetting
  el-tab-pane(
    label="当前生效配置"
    name="currentValid"
  )
    pre {{currentValid.data}}
.configuration
  el-tabs(
    v-model="active"
  )
    //- 所有配置
    +AllSetting

    //- 生效配置
    +ActiveSetting


</template>
<script lang="ts">
import { defineComponent } from "vue";

import useConfigState, { configListValid } from "../../states/config";
import ConfigTable from "../../components/configs/Table.vue";

export default defineComponent({
  name: "Configuration",
  components: {
    ConfigTable,
  },
  setup() {
    const configState = useConfigState();
    return {
      currentValid: configState.currentValidConfig,
    };
  },
  data() {
    return {
      active: "all",
    };
  },
  beforeMount() {
    this.fetch();
  },
  methods: {
    async fetch() {
      try {
        await configListValid();
      } catch (err) {
        this.$error(err);
      }
    },
  },
});
</script>
<style lang="stylus" scoped>
@import "../../common";
.configuration
  margin $mainMargin
  padding 20px 30px
  background-color $white
pre
  margin 20px
</style>
