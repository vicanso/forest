<template lang="pug">
el-card.httpErrors
  template(
    #header
  )
    i.el-icon-user-solid
    span HTTP出错查询 
  div(
    v-loading="httpErrors.processing"
  )
    base-filter(
      :fields="filterFields"
      @filter="filter"
    )
    el-table(
      :data="httpErrors.items"
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
        prop="method"
        key="method"
        label="Method"
        width="80"
      )
      el-table-column(
        prop="route"
        key="route"
        label="路由"
        width="150"
      )
      el-table-column(
        prop="category"
        key="category"
        label="类型"
        width="100"
      )
      el-table-column(
        prop="statusCode"
        key="statusCode"
        label="状态码"
        width="80"
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
        prop="ip"
        key="ip"
        label="IP"
        width="100"
      )
      //- uri 
      el-table-column(
        label="URI"
        width="80"
      ): template(
        #default="scope"
      ): base-tooltip(
        icon="el-icon-info"
        :content="scope.row.uri"
      )
      //- error 
      el-table-column(
        prop="error"
        key="error"
        label="Error"
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
import { useFluxStore } from "../store";

// 最近一小时
const defaultDateRange = [new Date(Date.now() - 60 * 60 * 1000), new Date()];
const filterFields = [
  {
    label: "账号：",
    key: "account",
    placeholder: "请输入要查询的账号",
    clearable: true,
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
    type: "dateTimeRange",
    placeholder: ["开始日期", "结束日期"],
    defaultValue: defaultDateRange,
    span: 9,
  },
  {
    label: "",
    type: "filter",
    labelWidth: "0px",
    span: 3,
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
      disableBeforeMountFetch: true,
      filterFields,
      pageSizes: PAGE_SIZES,
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
      return getUniqueKey(this.httpErrors.items, "tid");
    },
    sessionIDFilters() {
      return getUniqueKey(this.httpErrors.items, "sid");
    },
  },
  methods: {
    // filterResult(value, row) {
    //   return row.result == value;
    // },
    filterTrack(value, row) {
      return row.tid == value;
    },
    filterSession(value, row) {
      return row.sid == value;
    },
    async fetch() {
      const { httpErrors, query } = this;
      if (httpErrors.processing) {
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
        await this.listHTTPError(params);
      } catch (err) {
        this.$error(err);
      }
    },
  },
  setup() {
    const fluxStore = useFluxStore();
    return {
      httpErrors: fluxStore.state.httpErrors,
      listHTTPError: (params) => fluxStore.dispatch("listHTTPError", params),
    };
  },
});
</script>

<style lang="stylus" scoped>
@import "../common";
.httpErrors
  margin: $mainMargin
  i
    margin-right: 5px
.pagination
  text-align: right
  margin-top: 15px
</style>
