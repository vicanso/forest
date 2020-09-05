<template>
  <div class="routerConcurrency">
    <div v-if="!editMode">
      <ConfigTable :category="category" name="路由并发配置" />
      <div class="add">
        <el-button class="addBtn" type="primary" @click="add">添加</el-button>
      </div>
    </div>
    <ConfigEditor
      name="添加/更新路由并发配置"
      summary="配置针对各路由并发请求的限制"
      :category="category"
      :defaultValue="defaultValue"
      v-else
    >
      <template v-slot:data="configProps">
        <RouterConcurrencyData
          :data="configProps.form.data"
          @change="configProps.form.data = $event"
        />
      </template>
    </ConfigEditor>
  </div>
</template>
<script>
import { ROUTER_CONCURRENCY } from "@/constants/config";
import { CONFIG_EDITE_MODE } from "@/constants/route";
import ConfigEditor from "@/components/configs/Editor.vue";
import ConfigTable from "@/components/configs/Table.vue";
import RouterConcurrencyData from "@/components/configs/RouterConcurrencyData.vue";

export default {
  name: "RouterConcurrency",
  components: {
    ConfigEditor,
    ConfigTable,
    RouterConcurrencyData
  },
  data() {
    return {
      defaultValue: {
        category: ROUTER_CONCURRENCY
      },
      category: ROUTER_CONCURRENCY
    };
  },
  computed: {
    editMode() {
      return this.$route.query.mode === CONFIG_EDITE_MODE;
    }
  },
  methods: {
    add() {
      this.$router.push({
        query: {
          mode: CONFIG_EDITE_MODE
        }
      });
    }
  }
};
</script>
<style lang="sass" scoped>
@import "@/common.sass"
.add
  margin: $mainMargin
.addBtn
  width: 100%
</style>
