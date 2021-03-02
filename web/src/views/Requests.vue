<template lang="pug">

mixin AddrColumn
  el-table-column(
    prop="addr"
    key="addr"
    label="Addr"
    width="140"
  )

mixin HostnameColumn
  el-table-column(
    prop="hostname"
    key="hostname"
    label="HostName"
    width="120"
  )

mixin ServiceColumn
  el-table-column(
    prop="service"
    key="service"
    label="服务名称"
    width="120"
  )

mixin URIColumn
  el-table-column(
    label="URI"
    width="80"
  ): template(
    #default="scope"
  ): base-tooltip(
    v-if="scope.row.uri"
    :content="scope.row.method + ' ' + scope.row.uri"
    icon="el-icon-info"
  )

mixin RouteColumn
  el-table-column(
    prop="route"
    key="route"
    label="路由"
    width="200"
  )

mixin StatusColumn
  el-table-column(
    prop="status"
    key="status"
    label="状态码"
    width="100"
  )

mixin UseColumn
  el-table-column(
    label="耗时"
    width="100"
  ): template(
    #default="scope"
  ): el-tooltip(
    placement="bottom" 
    v-if="scope.row.use"
  )
    template(
      #content
    )
      ul
        li dns: {{scope.row.dnsUse}}ms
        li tcp: {{scope.row.tcpUse}}ms
        li tls: {{scope.row.tlsUse}}ms
        li processing: {{scope.row.processingUse}}ms
    span {{scope.row.use}}ms

mixin ResultColumn
  el-table-column(
    label="结果"
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

mixin ErrCategoryColumn
  el-table-column(
    label="出错类型"
    width="120"
    prop="errCategory"
    key="errCategory"
  )

mixin ErrExceptionColumn
  el-table-column(
    label="是否异常"
    width="100"
  ): template(
    #default="scope"
  )
    span(
      v-if="scope.row.exception"
    ) 是
    span(
      v-else
    ) 否

mixin ErrMessageColumn
  el-table-column(
    label="出错信息"
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
  ): template(
    #default="scope"
  ): time-formater(
    :time="scope.row._time"
  )

el-card.requests
  template(
    #header
  )
    i.el-icon-connection
    span 后端HTTP请求
  div(
    v-loading="requests.processing"
  )
    base-filter(
      v-if="inited"
      :fields="filterFields"
      :filter="filter"
    )

    el-table(
      :data="requests.items"
      row-key="_time"
      stripe
      :default-sort="{ prop: '_time', order: 'descending' }"
    )
      //- 服务名称
      +ServiceColumn
      //- 路由名称
      +RouteColumn
      //- 完整的请求地址
      +URIColumn
      //- 接口调用结果
      +ResultColumn
      //- 状态码
      +StatusColumn
      //- 接口调用解析的地址
      +AddrColumn
      //- 请求耗时
      +UseColumn
      //- 机器的host
      +HostnameColumn
      //- 出错类型
      +ErrCategoryColumn
      //- 是否异常
      +ErrExceptionColumn
      //- 出错信息
      +ErrMessageColumn
      //- 时间
      +TimeColumn
      
</template>

<script lang="ts">
import { defineComponent, onUnmounted } from "vue";

import { getDateTimeShortcuts, formatDateWithTZ } from "../helpers/util";
import BaseFilter from "../components/base/Filter.vue";
import BaseTooltip from "../components/Tooltip.vue";
import TimeFormater from "../components/TimeFormater.vue";
import BaseJson from "../components/base/JSON.vue";
import HTTPErrorFormater from "../components/HTTPErrorFormater.vue";
import { PAGE_SIZES } from "../constants/common";
import FilterTable from "../mixins/FilterTable";

import useFluxState, {
  fluxListRequest,
  fluxListRequestClear,
  fluxListRequestService,
  fluxListRequestRoute,
} from "../states/flux";

// 最近一小时
const defaultDateRange = [new Date(Date.now() - 60 * 60 * 1000), new Date()];
// 服务名称列表
const services = [];
// 服务请求路由列表
const routes = [];

const filterFields = [
  {
    label: "结果：",
    key: "result",
    type: "select",
    placeholder: "请选择请求结果",
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
    label: "服务：",
    key: "service",
    type: "select",
    placeholder: "请选择服务名称",
    options: services,
    span: 6,
  },
  {
    label: "路由：",
    key: "route",
    type: "select",
    placeholder: "请选择路由",
    options: routes,
    span: 6,
  },
  {
    label: "异常：",
    key: "exception",
    type: "select",
    placeholder: "请选择是否异常出错",
    options: [
      {
        name: "全部",
        value: "",
      },
      {
        name: "是",
        value: "1",
      },
      {
        name: "否",
        value: "0",
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
    type: "dateTimeRange",
    placeholder: ["开始日期", "结束日期"],
    shortcuts: getDateTimeShortcuts(["1h", "2h", "3h", "12h", "1d"]),
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

export default defineComponent({
  name: "Requests",
  components: {
    BaseFilter,
    BaseTooltip,
    TimeFormater,
    HTTPErrorFormater,
    BaseJson,
  },
  mixins: [FilterTable],
  setup() {
    onUnmounted(() => {
      fluxListRequestClear();
    });
    const fluxState = useFluxState();
    return {
      requests: fluxState.requests,
      requestServices: fluxState.requestServices,
      requestRoutes: fluxState.requestRoutes,
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
        exception: false,
      },
    };
  },
  async beforeMount() {
    try {
      // 获取服务名称列表
      await fluxListRequestService();
      services.length = 0;
      services.push({
        name: "全部",
        value: "",
      });
      this.requestServices.items.forEach((element) => {
        services.push({
          name: element,
          value: element,
        });
      });

      // 获取路由名称列表
      await fluxListRequestRoute();
      routes.length = 0;
      routes.push({
        name: "全部",
        value: "",
      });
      this.requestRoutes.items.forEach((element) => {
        routes.push({
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
    async fetch() {
      const { requests, query } = this;
      if (requests.processing) {
        return;
      }
      const params = Object.assign({}, query);
      const value = params.dateRange;
      if (!value || value.length !== 2) {
        this.$erro("时间区间不能为空");
        return;
      }
      params.begin = formatDateWithTZ(value[0]);
      params.end = formatDateWithTZ(value[1]);
      delete params.dateRange;
      try {
        await fluxListRequest(params);
      } catch (err) {
        this.$error(err);
      }
    },
  },
});
</script>

<style lang="stylus" scoped>
@import "../common";
.requests
  margin $mainMargin
  i
    margin-right 5px
.pagination
  text-align right
  margin-top 15px
</style>
