<template lang="pug">
mixin AccountColumn
  el-table-column(
    prop="account"
    key="account"
    label="账户"
    width="150"
  )

mixin CategoryColumn
  el-table-column(
    prop="category"
    key="category"
    label="类型"
    width="150"
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

mixin RouteColumn
  el-table-column(
    label="路由"
    prop="route"
    key="route"
  )

mixin TIDColumn
  el-table-column(
    label="Track ID"
  ): template(
    #default="scope"
  ): base-tooltip(
    :content="scope.row.tid"
  )

mixin FullPathColumn
  el-table-column(
    label="完整路径"
    width="80"
  ): template(
    #default="scope"
  ): base-tooltip(
    icon="el-icon-info"
    :content="scope.row.path"
  )

mixin ErrorColumn
  el-table-column(
    label="error"
    width="80"
  ): template(
    #default="scope"
  ): base-tooltip(
    icon="el-icon-info"
    :content="scope.row.message"
  )

mixin TimeColumn
  el-table-column(
    label="时间"
    prop="_time"
    key="_time"
    width="160"
  ): template(
    #default="scope"
  ): time-formater(
    :time="scope.row._time"
  )

el-card.actions
  template(
    #header
  )
    i.el-icon-info
    span 客户端行为查询
  div(
    v-loading="userActions.processing"
  )
    base-filter(
      v-if="inited"
      :fields="filterFields"
      :filter="filter"
    )
    el-table(
      :data="userActions.items"
      row-key="_time"
      stripe
      :default-sort="{ prop: '_time', order: 'descending' }"
    )
      //- 账号
      +AccountColumn

      //- 记录类型
      +CategoryColumn

      //- 记录状态
      +StatusColumn

      //- 触发记录时所有route
      +RouteColumn

      //- tid
      +TIDColumn

      //- full path
      +FullPathColumn

      //- error
      +ErrorColumn

      //- 时间
      +TimeColumn

</template>
<script lang="ts">
import { defineComponent, onUnmounted } from "vue";
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
import useFluxState, {
  fluxListClientActionCategory,
  fluxListClientAction,
  fluxListClientActionClear,
} from "../states/flux";

const defaultDateRange = [today(), today()];
const categoryOptions = [];
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
    key: "category",
    type: "select",
    placeholder: "请选择要筛选的分类",
    options: categoryOptions,
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

export default defineComponent({
  name: "Actions",
  components: {
    BaseFilter,
    BaseTooltip,
    TimeFormater,
    BaseJson,
  },
  mixins: [FilterTable],
  setup() {
    onUnmounted(() => {
      fluxListClientActionClear();
    });
    const fluxState = useFluxState();
    return {
      userActionCategories: fluxState.clientActionCategories,
      userActions: fluxState.clientActions,
    };
  },
  data() {
    return {
      inited: false,
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
  async beforeMount() {
    try {
      await fluxListClientActionCategory();

      categoryOptions.length = 0;
      categoryOptions.push({
        name: "全部",
        value: "",
      });
      this.userActionCategories.items.forEach((element) => {
        categoryOptions.push({
          name: element,
          value: element,
        });
      });
    } catch (err) {
      this.$error(err);
    } finally {
      this.inited = true;
    }
  },
  methods: {
    async fetch() {
      const { userActions, query } = this;
      if (userActions.processing) {
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
        await fluxListClientAction(params);
      } catch (err) {
        this.$error(err);
      }
    },
  },
});
</script>

<style lang="stylus" scoped>
@import "../common";
.actions
  margin $mainMargin
  i
    margin-right 5px
.pagination
  text-align right
  margin-top 15px
</style>
