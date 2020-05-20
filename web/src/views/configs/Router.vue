<template>
  <div class="router">
    <div v-if="!editMode">
      <ConfigTable :category="category" name="路由配置" />
      <div class="add">
        <el-button class="addBtn" type="primary" @click="add">添加</el-button>
      </div>
    </div>
    <ConfigEditor
      name="添加/更新路由配置"
      summary="配置针对各路由响应的Mock"
      :category="category"
      :defaultValue="defaultValue"
      v-else
    >
      <template v-slot:data="configProps">
        <RouterData
          :data="configProps.form.data"
          @change="configProps.form.data = $event"
        />
      </template>
    </ConfigEditor>
  </div>
</template>
<script>
import { ROUTER } from "@/constants/config";
import { CONFIG_EDITE_MODE } from "@/constants/route";
import ConfigEditor from "@/components/configs/Editor.vue";
import ConfigTable from "@/components/configs/Table.vue";
import RouterData from "@/components/configs/RouterData.vue";

export default {
  name: "Router",
  components: {
    RouterData,
    ConfigTable,
    ConfigEditor
  },
  data() {
    return {
      defaultValue: {
        category: ROUTER
      },
      category: ROUTER
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
