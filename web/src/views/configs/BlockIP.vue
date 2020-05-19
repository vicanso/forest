<template>
  <div class="blockIP">
    <div v-if="!editMode">
      <ConfigTable :category="category" name="黑名单IP配置" />
      <div class="add">
        <el-button class="addBtn" type="primary" @click="add">添加</el-button>
      </div>
    </div>
    <ConfigEditor
      name="添加/更新IP黑名单配置"
      summary="用于拦截访问IP"
      :category="category"
      :defaultValue="defaultValue"
      v-else
    />
  </div>
</template>
<script>
import { BLOCK_IP } from "@/constants/config";
import { CONFIG_EDITE_MODE } from "@/constants/route";
import ConfigEditor from "@/components/configs/Editor.vue";
import ConfigTable from "@/components/configs/Table.vue";

export default {
  name: "BlockIP",
  components: {
    ConfigTable,
    ConfigEditor
  },
  data() {
    return {
      defaultValue: {
        category: BLOCK_IP
      },
      category: BLOCK_IP
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
