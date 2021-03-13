<template lang="pug">
//- 配置列表
mixin Table
  config-table(
    :category="category"
    name="请求并发配置"
  )
  .add: el-button.addBtn(
    type="primary"
    @click="add"
  ) 添加
//- 配置编辑
mixin Editor
  config-editor(
    name="添加/更新请求并发配置"
    summary="配置针对各请求实例的并发限制"
    :category="category"
    :defaultValue="defaultValue"
  ): template(
    #data="configProps"
  ): request-concurrency-data(
    :data="configProps.form.data"
    @change.self="configProps.form.data = $event"
  )
  

.requestConcurrency
  //- 配置表格
  div(
    v-if="!editMode"
  )
    +Table
  //- 编辑
  template(
    v-else
  )
    +Editor
</template>

<script lang="ts">
import { defineComponent } from "vue";

import ConfigEditor from "../../components/configs/Editor.vue";
import RequestConcurrencyData from "../../components/configs/RequestConcurrencyData.vue";
import ConfigTable from "../../components/configs/Table.vue";
import { REQUEST_CONCURRENCY, CONFIG_EDIT_MODE } from "../../constants/common";

export default defineComponent({
  name: "RequestConcurrency",
  components: {
    ConfigTable,
    ConfigEditor,
    RequestConcurrencyData,
  },
  data() {
    return {
      defaultValue: {
        category: REQUEST_CONCURRENCY,
      },
      category: REQUEST_CONCURRENCY,
    };
  },
  computed: {
    editMode() {
      return this.$route.query.mode === CONFIG_EDIT_MODE;
    },
  },
  methods: {
    add() {
      this.$router.push({
        query: {
          mode: CONFIG_EDIT_MODE,
        },
      });
    },
  },
});
</script>

<style lang="stylus" scoped>
@import "../../common";

.add
  margin $mainMargin
.addBtn
  width 100%
</style>
