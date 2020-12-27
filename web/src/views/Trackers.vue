<template>
  <el-card class="trackers">
    <div slot="header">
      <i class="el-icon-user-solid" />
      <span>用户行为查询</span>
    </div>
    <div v-loading="processing">
      <BaseFilter :fields="filterFields" @filter="filter" />
      <el-table :data="trackers" row-key="_time" stripe>
        <el-table-column
          prop="account"
          key="account"
          label="账户"
          width="120"
        />
        <el-table-column prop="action" key="action" label="类型" width="150" />
        <el-table-column label="状态" width="80">
          <template slot-scope="scope">
            <span v-if="scope.row.result == '0'">成功</span>
            <span v-else>失败</span>
          </template>
        </el-table-column>
        <el-table-column prop="form" key="form" label="Form" />
        <el-table-column prop="query" key="query" label="Query" />
        <el-table-column prop="params" key="params" label="Params" />
        <el-table-column label="Session ID">
          <template slot-scope="scope">
            <BaseToolTip :content="scope.row.sid" />
          </template>
        </el-table-column>
        <el-table-column label="Track ID">
          <template slot-scope="scope">
            <BaseToolTip :content="scope.row.tid" />
          </template>
        </el-table-column>
        <el-table-column label="IP">
          <template slot-scope="scope">
            <BaseToolTip icon="el-icon-info" :content="scope.row.ip" />
          </template>
        </el-table-column>
        <el-table-column
          prop="timeDesc"
          key="timeDesc"
          label="时间"
          width="180"
        />
      </el-table>
    </div>
  </el-card>
</template>

<script>
import { mapActions, mapState } from "vuex";
import { today, formatBegin, formatEnd } from "@/helpers/util";
import BaseTable from "@/components/base/Table.vue";
import BaseFilter from "@/components/base/Filter.vue";
import BaseToolTip from "@/components/base/ToolTip.vue";
import { PAGE_SIZES } from "@/constants/common";

const defaultDateRange = [today(), today()];
const filterFields = [
  {
    label: "账号：",
    key: "account",
    placeholder: "请输入要查询的账号",
    clearable: true,
    span: 6
  },
  {
    label: "时间：",
    key: "dateRange",
    type: "dateRange",
    placeholder: ["开始日期", "结束日期"],
    defaultValue: defaultDateRange,
    span: 12
  },
  {
    label: "",
    type: "filter",
    labelWidth: "0px",
    span: 6
  }
];
export default {
  name: "Trackers",
  extends: BaseTable,
  components: {
    BaseFilter,
    BaseToolTip
  },
  data() {
    return {
      disableBeforeMountFetch: true,
      filterFields,
      pageSizes: PAGE_SIZES,
      query: {
        dateRange: defaultDateRange,
        offset: 0,
        limit: PAGE_SIZES[0],
        account: ""
      }
    };
  },
  computed: {
    ...mapState({
      processing: state => state.tracker.listProcessing,
      trackers: state => state.tracker.list.data || []
    })
  },
  methods: {
    ...mapActions(["listTracker"]),
    async fetch() {
      const { query, processing } = this;
      if (processing) {
        return;
      }
      if (!query.account) {
        this.$message.error("账号不能为空");
        return;
      }
      const value = query.dateRange;
      if (value) {
        query.begin = formatBegin(value[0]);
        query.end = formatEnd(value[1]);
      } else {
        query.begin = "";
        query.end = "";
      }
      delete query.dateRange;
      try {
        await this.listTracker(query);
      } catch (err) {
        this.$message.error(err.message);
      }
    }
  }
};
</script>
<style lang="sass" scoped>
@import "@/common.sass"
.trackers
  margin: $mainMargin
  i
    margin-right: 5px
.pagination
  text-align: right
  margin-top: 15px
</style>
