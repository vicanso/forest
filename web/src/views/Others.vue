<template lang="pug">
.others
  el-card
    template(
      #header
    )
      span 随机字符串生成
    pre.randomKeys(
      v-loading="randomKeys.processing"      
    ) {{randomKeys.items.join("\n")}}

    ex-button(
      :onClick="listRandomKey"
    ) 生成随机字符串

</template>

<script lang="ts">
import { defineComponent } from "vue";

import ExButton from "../components/ExButton.vue";
import useCommonState, { commonListRandomKey } from "../store/common";

export default defineComponent({
  name: "Others",
  components: {
    ExButton,
  },
  setup() {
    const commonState = useCommonState();
    return {
      randomKeys: commonState.randomKeys,
      listRandomKey: commonListRandomKey,
    };
  },
});
</script>
<style lang="stylus" scoped>
@import "../common";

.others
  margin: $mainMargin

.randomKeys
  padding: 10px
  line-height: 2em
  height: 150px
.btn
  width: 100%
</style>
