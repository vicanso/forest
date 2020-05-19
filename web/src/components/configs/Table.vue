<template>
  <el-card class="configurationList">
    <div slot="header">
      <i class="el-icon-s-tools" />
      <span>{{ $props.name || "系统配置" }}</span>
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
        width="120"
      />
      <el-table-column
        prop="statusDesc"
        key="statusDesc"
        label="状态"
        width="80"
      />
      <el-table-column
        sortable
        prop="beginDateDesc"
        key="beginDateDesc"
        label="开始时间"
        width="180"
      />
      <el-table-column
        sortable
        prop="endDateDesc"
        key="endDateDesc"
        label="结束时间"
        width="180"
      />
      <el-table-column prop="data" key="data" label="配置数据" />
      <el-table-column prop="owner" key="owner" label="创建者" width="150" />
      <el-table-column
        sortable
        prop="updatedAtDesc"
        key="updatedAtDesc"
        label="更新时间"
        width="180"
      />
      <el-table-column fixed="right" label="操作" width="120">
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
import { CONFIG_EDITE_MODE } from "@/constants/route";

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
      query: {
        category: this.$props.category
      }
    };
  },
  computed: mapState({
    configs: state => state.config.items || [],
    processing: state => state.config.processing,
    userAccount: state => state.user.info.account
  }),
  methods: {
    ...mapActions(["listConfig", "removeConfigByID"]),
    modify(item) {
      this.$router.push({
        query: {
          mode: CONFIG_EDITE_MODE,
          id: item.id
        }
      });
    },
    async remove(item) {
      try {
        await this.removeConfigByID(item.id);
      } catch (err) {
        this.$message.error(err.message);
      }
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
</style>
