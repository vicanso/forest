<template lang="pug">
mixin AccountColumn
  el-table-column(
    prop="account"
    key="account"
    label="账户"
    width="120"
    fixed="left"
  )

mixin CategoryColumn
  el-table-column(
    prop="action"
    key="action"
    label="类型"
    width="200"
  )

mixin StatusColumn
  el-table-column(
    label="状态"
    width="80"
  ): template(
    #default="scope"  
  )
    span(
      v-if="scope.row.result === '0'"
    ) 成功
    span(
      v-else
    ) 失败

mixin FormColumn
  el-table-column(
    label="Form"
    width="100"
  ): template(
    #default="scope"
  ): base-json(
    :content="scope.row.form"
    icon="el-icon-info"
  )

mixin QueryColumn
  el-table-column(
    label="Query"
    width="100"
  ): template(
    #default="scope"
  ): base-json(
    :content="scope.row.query"
    icon="el-icon-info"
  )

mixin ParamsColumn
  el-table-column(
    label="Params"
    width="200"
  ): template(
    #default="scope"
  ): base-json(
    :content="scope.row.params"
  )

mixin SessionIDColumn
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

mixin TrackIDColumn
  el-table-column(
    label="Track ID"
    :filters="trackIDFilters"
    :filter-method="filterTrack"
    width="100"
  ): template(
    #default="scope"
  ): base-tooltip(
    :content="scope.row.tid"
  )

mixin IPColumn
  el-table-column(
    label="IP"
    width="80"
  ): template(
    #default="scope"
  ): base-tooltip(
    icon="el-icon-info"
    :content="scope.row.ip"
  )

mixin ErrorColumn
  el-table-column(
    label="Error"
  ): template(
    #default="scope"
  ): HTTPErrorFormater(
    :message="scope.row.error"
  )

mixin TimeColumn
  el-table-column(
    label="时间"
    prop="_time"
    key="_time"
    width="160"
    fixed="right"
  )
    template(
      #header
    ): StatsSummary(
      v-if="!trackers.processing"
      :data="trackers.items"
      :fields="summaryFields"
      @filter="doFilter"
    )

    template(
      #default="scope"
    ): time-formater(
      :time="scope.row._time"
    )

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
      v-if="inited"
      :fields="filterFields"
      :filter="filter"
    )
    StatsTable(
      v-if="!trackers.processing"
      :data="trackers.items"
      :flux="trackers.flux"
    ): template(
      #default
    )
      //- 账号
      +AccountColumn

      //- 类别
      +CategoryColumn

      //- 状态筛选
      +StatusColumn
      
      //- form参数
      +FormColumn
      
      //- query参数
      +QueryColumn
      
      //- params参数
      +ParamsColumn
      
      //- session id
      +SessionIDColumn

      //- track id
      +TrackIDColumn

      //- ip
      +IPColumn

      //- error
      +ErrorColumn

      //- 时间
      +TimeColumn

</template>

<script lang="ts">
import { defineComponent, onUnmounted, reactive, provide } from "vue";

import {
  today,
  getDateDayShortcuts,
  formatBegin,
  formatEnd,
} from "../helpers/util";
import BaseFilter from "../components/base/Filter.vue";
import BaseTooltip from "../components/Tooltip.vue";
import TimeFormater from "../components/TimeFormater.vue";
import BaseJson from "../components/base/JSON.vue";
import { PAGE_SIZES } from "../constants/common";
import FilterTable from "../mixins/FilterTable";
import HTTPErrorFormater from "../components/HTTPErrorFormater.vue";
import StatsSummary from "../components/StatsSummary.vue";
import StatsTable from "../components/StatsTable.vue";
import userFluxState, {
  fluxListUserTrackAction,
  fluxListUserTracker,
  fluxListUserTrackerClear,
} from "../states/flux";

const defaultDateRange = [today(), today()];
const actionOptions = [];
const filterFields = [
  {
    label: "账号：",
    key: "account",
    placeholder: "请输入要查询的账号",
    clearable: true,
    span: 6,
  },
  {
    label: "类型：",
    key: "action",
    type: "select",
    placeholder: "请选择要筛选的分类",
    options: actionOptions,
    span: 6,
  },
  {
    label: "结果：",
    key: "result",
    type: "select",
    placeholder: "请选择要筛选的结果",
    options: [
      {
        name: "全部",
        value: "",
      },
      {
        name: "成功",
        value: "0",
      },
      {
        name: "失败",
        value: "1",
      },
    ],
    span: 6,
  },
  {
    label: "数量：",
    key: "limit",
    type: "number",
    placeholder: "请输入最大数量",
    clearable: true,
    defaultValue: 100,
    span: 6,
  },
  {
    label: "时间：",
    key: "dateRange",
    type: "dateRange",
    placeholder: ["开始日期", "结束日期"],
    shortcuts: getDateDayShortcuts(["1d", "2d", "3d", "7d"]),
    defaultValue: defaultDateRange,
    span: 18,
  },
  {
    label: "",
    type: "filter",
    labelWidth: "0px",
    span: 6,
  },
];

function getUniqueKey(data: Record<string, unknown>[], key: string) {
  if (!data || !data.length) {
    return [];
  }
  const keys = {};
  data.forEach((item) => {
    if (item[key]) {
      keys[`${item[key]}`] = true;
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
    HTTPErrorFormater,
    BaseJson,
    StatsSummary,
    StatsTable,
  },
  mixins: [FilterTable],
  setup() {
    onUnmounted(() => {
      fluxListUserTrackerClear();
    });
    const statsParams = reactive({
      filters: {},
    });
    provide("statsParams", statsParams);
    const fluxState = userFluxState();
    return {
      statsParams,
      trackers: fluxState.userTrackers,
      trackerActions: fluxState.userTrackerActions,
    };
  },
  data() {
    return {
      inited: false,
      disableBeforeMountFetch: true,
      filterFields,
      pageSizes: PAGE_SIZES,
      summaryFields: ["account", "action", "ip", "sid", "tid", "result"],
      query: {
        dateRange: defaultDateRange,
        offset: 0,
        limit: 100,
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
  async beforeMount() {
    try {
      await fluxListUserTrackAction();

      actionOptions.length = 0;
      actionOptions.push({
        name: "全部",
        value: "",
      });
      this.trackerActions.items.forEach((element) => {
        actionOptions.push({
          name: element,
          value: element,
        });
      });
      this.inited = true;
    } catch (err) {
      this.$error(err);
    }
  },
  methods: {
    doFilter(filters) {
      this.statsParams.filters = filters;
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
      const value = params.dateRange;
      if (!value) {
        this.$erro("时间区间不能为空");
        return;
      }
      params.begin = formatBegin(value[0]);
      params.end = formatEnd(value[1]);
      delete params.dateRange;
      try {
        await fluxListUserTracker(params);
      } catch (err) {
        this.$error(err);
      }
    },
  },
});
</script>

<style lang="stylus" scoped>
@import "../common";
.trackers
  margin $mainMargin
  i
    margin-right 5px
.pagination
  text-align right
  margin-top 15px
</style>
