<template>
  <el-card class="trackers">
    <div slot="header">
      <i class="el-icon-user-solid" />
      <span>用户行为查询</span>
    </div>
    <div v-loading="processing">
      <BaseFilter :fields="filterFields" @filter="filter" />
      <el-table
        :data="trackers"
        row-key="_time"
        stripe
        :default-sort="{ prop: 'timeDesc', order: 'descending' }"
      >
        <el-table-column
          prop="account"
          key="account"
          label="账户"
          width="120"
        />
        <el-table-column prop="action" key="action" label="类型" width="150" />
        <el-table-column
          label="状态"
          width="80"
          :filters="resultFilters"
          :filter-method="filterResult"
        >
          <template slot-scope="scope">
            <span v-if="scope.row.result == '0'">成功</span>
            <span v-else>失败</span>
          </template>
        </el-table-column>
        <el-table-column label="Form">
          <template slot-scope="scope">
            <BaseJSON :content="scope.row.form" />
          </template>
        </el-table-column>
        <el-table-column label="Query">
          <template slot-scope="scope">
            <BaseJSON :content="scope.row.query" />
          </template>
        </el-table-column>
        <el-table-column label="Params">
          <template slot-scope="scope">
            <BaseJSON :content="scope.row.params" />
          </template>
        </el-table-column>
        <el-table-column
          label="Session ID"
          :filters="sessionIDFilters"
          :filter-method="filterSession"
          width="110"
        >
          <template slot-scope="scope">
            <BaseToolTip :content="scope.row.sid" />
          </template>
        </el-table-column>
        <el-table-column
          label="Track ID"
          :filters="trackIDFilters"
          :filter-method="filterTrack"
          width="90"
        >
          <template slot-scope="scope">
            <BaseToolTip :content="scope.row.tid" />
          </template>
        </el-table-column>
        <el-table-column label="IP" width="80">
          <template slot-scope="scope">
            <BaseToolTip icon="el-icon-info" :content="scope.row.ip" />
          </template>
        </el-table-column>
        <el-table-column
          prop="timeDesc"
          key="timeDesc"
          label="时间"
          width="180"
          sortable
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
import BaseJSON from "@/components/base/JSON.vue";
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

function getUniqueKey(data, key) {
  if (!data || !data.length) {
    return [];
  }
  const keys = {};
  data.forEach(item => {
    if (item[key]) {
      keys[item[key]] = true;
    }
  });
  return Object.keys(keys).map(item => {
    return {
      text: item,
      value: item
    };
  });
}

export default {
  name: "Trackers",
  extends: BaseTable,
  components: {
    BaseFilter,
    BaseToolTip,
    BaseJSON
  },
  data() {
    return {
      resultFilters: [
        {
          text: "成功",
          value: "0"
        },
        {
          text: "失败",
          value: "1"
        }
      ],
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
      trackIDFilters: state => getUniqueKey(state.tracker.list.data, "tid"),
      sessionIDFilters: state => getUniqueKey(state.tracker.list.data, "sid"),
      processing: state => state.tracker.listProcessing,
      trackers: state => state.tracker.list.data || []
    })
  },
  methods: {
    ...mapActions(["listTracker"]),
    filterResult(value, row) {
      return row.result == value;
    },
    filterTrack(value, row) {
      return row.tid == value;
    },
    filterSession(value, row) {
      return row.sid == value;
    },
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
