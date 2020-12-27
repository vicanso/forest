<template>
  <el-card class="configurationList">
    <div slot="header">
      <i class="el-icon-s-tools" />
      <span>{{ $props.name || "系统配置" }}</span>
      <span class="filters">
        <el-checkbox title="仅展示有效的" v-model="available"
          >仅展示有效配置</el-checkbox
        >
        <el-checkbox title="展开所有配置" v-model="expanded"
          >展开所有配置</el-checkbox
        >
      </span>
    </div>
    <el-table
      v-loading="processing"
      :data="configs"
      row-key="id"
      stripe
      :default-sort="{ prop: 'updatedAtDesc', order: 'descending' }"
    >
      <el-table-column prop="id" key="id" label="ID" width="80" />
      <el-table-column prop="name" key="name" label="名称" width="120" />
      <el-table-column
        prop="category"
        key="category"
        label="分类"
        width="150"
      />
      <el-table-column
        sortable
        prop="statusDesc"
        key="statusDesc"
        label="状态"
        width="80"
      />
      <el-table-column
        sortable
        prop="startedAtDesc"
        key="startedAtDesc"
        label="开始时间"
        width="180"
      />
      <el-table-column
        sortable
        prop="endedAtDesc"
        key="endedAtDesc"
        label="结束时间"
        width="180"
      />
      <el-table-column
        prop="data"
        key="data"
        label="配置数据"
        :width="configWidth"
      >
        <template slot-scope="scope">
          <pre v-if="expanded">{{
            scope.row.isJSON ? `\n${scope.row.data}` : scope.row.data
          }}</pre>
          <el-tooltip placement="bottom" v-else-if="scope.row.isJSON">
            <pre slot="content">{{ scope.row.data }}</pre>
            <i class="el-icon-info" />
          </el-tooltip>
          <span v-else>{{ scope.row.data }}</span>
        </template>
      </el-table-column>
      <el-table-column prop="owner" key="owner" label="创建者" width="150" />
      <el-table-column
        sortable
        prop="updatedAtDesc"
        key="updatedAtDesc"
        label="更新时间"
        width="180"
      />
      <el-table-column fixed="right" label="操作">
        <template slot-scope="scope">
          <div v-if="scope.row.owner === userAccount">
            <el-popconfirm
              title="确定要删除该配置吗？"
              @onConfirm="remove(scope.row)"
            >
              <el-button class="op" slot="reference" type="text" size="small"
                >删除</el-button
              >
            </el-popconfirm>
            <el-button
              class="op"
              type="text"
              size="small"
              @click="modify(scope.row)"
              >编辑</el-button
            >
          </div>
          <span v-else>--</span>
        </template>
      </el-table-column>
    </el-table>
  </el-card>
</template>
<script>
import { mapActions, mapState } from "vuex";
import { CONFIG_EDIT_MODE } from "@/constants/route";
import { CONFIG_ENABLED } from "@/constants/config";

export default {
  name: "ConfigurationList",
  props: {
    name: String,
    category: {
      type: String,
      required: true
    }
  },
  data() {
    return {
      expanded: false,
      available: false,
      query: {
        category: this.$props.category
      }
    };
  },
  computed: {
    configWidth() {
      if (this.expanded) {
        return 200;
      }
      return 80;
    },
    ...mapState({
      configs: function(state) {
        const { available } = this;
        const arr = (state.config.items || []).filter(item => {
          // 如果非选择仅展示有效的
          if (!available) {
            return true;
          }
          if (item.status !== CONFIG_ENABLED) {
            return false;
          }
          const now = Date.now();
          const beginDate = new Date(item.startedAt).getTime();
          const endDate = new Date(item.endedAt).getTime();
          // 如果未到开始时间或者已结束
          if (beginDate > now || endDate < now) {
            return false;
          }
          return true;
        });
        return arr;
      },
      processing: state => state.config.processing,
      userAccount: state => state.user.info.account
    })
  },
  methods: {
    ...mapActions(["listConfig", "removeConfigByID"]),
    modify(item) {
      this.$router.push({
        query: {
          mode: CONFIG_EDIT_MODE,
          id: item.id
        }
      });
    },
    async remove() {
      this.$message.warning("不支持删除配置，请将配置禁用即可");
    }
  },
  async beforeMount() {
    const { query } = this;
    try {
      await this.listConfig(query);
    } catch (err) {
      this.$message.error(err.message);
    }
  }
};
</script>
<style lang="sass" scoped>
@import "@/common.sass"
.configurationList
  margin: $mainMargin
  i
    margin-right: 3px
  .op
    margin: 0 10px
  .filters
    margin-left: 20px
</style>
