<template lang="pug">
el-card.trackers
  template(
    #header
  )
    i.el-icon-user-solid
    span 用户行为查询
  div(
    v-loading="trackers.processing"
  )
    base-filter(
      :fields="filterFields"
      @filter="filter"
    )
    el-table(
      :data="trackers.items"
      row-key="_time"
      stripe
      :default-sort="{ prop: '_time', order: 'descending' }"
    )
      el-table-column(
        prop="account"
        key="account"
        label="账户"
        width="120"
      )
      el-table-column(
        prop="action"
        key="action"
        label="类型"
        width="150"
      )
      //- 状态筛选
      el-table-column(
        label="状态"
        width="80"
        :filters="resultFilters"
        :filter-method="filterResult"
      ): template(
        #default="scope"  
      )
        span(
          v-if="scope.row.result === '0'"
        ) 成功
        span(
          v-else
        ) 失败
      //- form参数
      el-table-column(
        label="Form"
      ): template(
        #default="scope"
      ): base-json(
        :content="scope.row.form"
      )
      //- query参数
      el-table-column(
        label="Query"
      ): template(
        #default="scope"
      ): base-json(
        :content="scope.row.query"
      )
      //- params参数
      el-table-column(
        label="Params"
      ): template(
        #default="scope"
      ): base-json(
        :content="scope.row.params"
      )
      //- session id
      el-table-column(
        label="Session ID"
        :filters="sessionIDFilters"
        :filter-method="filterSession"
        width="110"
      ): template(
        #default="scope"
      ): base-tooltip(
        :content="scope.row.sid"
      )
      //- track id
      el-table-column(
        label="Track ID"
        :filters="trackIDFilters"
        :filter-method="filterTrack"
        width="90"
      ): template(
        #default="scope"
      ): base-tooltip(
        :content="scope.row.tid"
      )
      //- ip
      el-table-column(
        label="IP"
        width="80"
      ): template(
        #default="scope"
      ): base-tooltip(
        icon="el-icon-info"
        :content="scope.row.ip"
      )
      //- 时间
      el-table-column(
        label="时间"
        prop="_time"
        key="_time"
        width="120"
      ): template(
        #default="scope"
      ): time-formater(
        :time="scope.row._time"
      )
</template>

<script lang="ts">
import { defineComponent } from "vue";

import { today, formatBegin, formatEnd } from "../helpers/util";
import BaseFilter from "../components/base/Filter.vue";
import BaseTooltip from "../components/Tooltip.vue";
import TimeFormater from "../components/TimeFormater.vue";
import BaseJson from "../components/base/JSON.vue";
import { PAGE_SIZES } from "../constants/common";
import FilterTable from "../mixins/FilterTable";
import { useTrackerStore } from "../store";

const defaultDateRange = [today(), today()];
const filterFields = [
  {
    label: "账号：",
    key: "account",
    placeholder: "请输入要查询的账号",
    clearable: true,
    span: 6,
  },
  {
    label: "时间：",
    key: "dateRange",
    type: "dateRange",
    placeholder: ["开始日期", "结束日期"],
    defaultValue: defaultDateRange,
    span: 12,
  },
  {
    label: "",
    type: "filter",
    labelWidth: "0px",
    span: 6,
  },
];

function getUniqueKey(data: any[], key: string) {
  if (!data || !data.length) {
    return [];
  }
  const keys = {};
  data.forEach((item) => {
    if (item[key]) {
      keys[item[key]] = true;
    }
  });
  return Object.keys(keys).map((item) => {
    return {
      text: item,
      value: item,
    };
  });
}

export default defineComponent({
  name: "Trackers",
  components: {
    BaseFilter,
    BaseTooltip,
    TimeFormater,
    BaseJson,
  },
  mixins: [FilterTable],
  data() {
    return {
      resultFilters: [
        {
          text: "成功",
          value: "0",
        },
        {
          text: "失败",
          value: "1",
        },
      ],
      disableBeforeMountFetch: true,
      filterFields,
      pageSizes: PAGE_SIZES,
      query: {
        dateRange: defaultDateRange,
        offset: 0,
        limit: PAGE_SIZES[0],
        account: "",
      },
    };
  },
  computed: {
    trackIDFilters() {
      return getUniqueKey(this.trackers.items, "tid");
    },
    sessionIDFilters() {
      return getUniqueKey(this.trackers.items, "sid");
    },
  },
  methods: {
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
      const { trackers, query } = this;
      if (trackers.processing) {
        return;
      }
      const params = Object.assign({}, query);
      if (!params.account) {
        this.$error("账号不能为空");
        return;
      }
      const value = params.dateRange;
      if (!value) {
        this.$erro("时间区间不能为空");
        return;
      }
      params.begin = formatBegin(value[0]);
      params.end = formatEnd(value[1]);
      delete params.dateRange;
      try {
        await this.list(params);
      } catch (err) {
        this.$error(err);
      }
    },
  },
  setup() {
    const trackerStore = useTrackerStore();
    return {
      trackers: trackerStore.state.trackers,
      list: (params) => trackerStore.dispatch("list", params),
    };
  },
});
</script>

<style lang="stylus" scoped>
@import "../common";
.trackers
  margin: $mainMargin
  i
    margin-right: 5px
.pagination
  text-align: right
  margin-top: 15px
</style>
