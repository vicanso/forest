<template>
  <div class="signedKey">
    <div v-if="!editMode">
      <ConfigTable :category="category" name="SignedKey配置" />
      <div class="add">
        <el-button class="addBtn" type="primary" @click="add">添加</el-button>
      </div>
    </div>
    <ConfigEditor
      name="添加/更新SignedKey配置"
      summary="用于配置session中使用的signed key"
      :category="category"
      :defaultValue="defaultValue"
      v-else
    />
  </div>
</template>
<script>
import { SIGNED_KEY } from "@/constants/config";
import { CONFIG_EDITE_MODE } from "@/constants/route";
import ConfigEditor from "@/components/configs/Editor.vue";
import ConfigTable from "@/components/configs/Table.vue";
export default {
  name: "SignedKey",
  components: {
    ConfigTable,
    ConfigEditor
  },
  data() {
    return {
      defaultValue: {
        category: SIGNED_KEY
      },
      category: SIGNED_KEY
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
