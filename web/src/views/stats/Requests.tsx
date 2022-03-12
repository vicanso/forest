import { useMessage, NPopover } from "naive-ui";
import { defineComponent, onMounted, onUnmounted, VNode } from "vue";
import { css } from "@linaria/core";
import { TableColumn } from "naive-ui/lib/data-table/src/interface";
import ExLoading from "../../components/ExLoading";
import { showError } from "../../helpers/util";
import ExTable from "../../components/ExTable";
import { getHoursAge } from "../../helpers/util";
import { FormItemTypes } from "../../components/ExForm";
import ExFluxDetail from "../../components/ExFluxDetail";
import { useRequestsStore } from "../../stores/requests";
import { storeToRefs } from "pinia";

const serviceOptions = [
  {
    label: "所有",
    value: "",
  },
];

const routeOptions = [
  {
    label: "所有",
    value: "",
  },
];

const useTimeListClass = css`
  margin: 0;
  padding: 0;
  list-style-position: inside;
`;

function getFilters() {
  return [
    {
      key: "rslt",
      name: "结果：",
      type: FormItemTypes.Select,
      placeholder: "请选择要筛选的结果",
      options: [
        {
          label: "所有",
          value: "",
        },
        {
          label: "成功",
          value: "0",
        },
        {
          label: "失败",
          value: "1",
        },
      ],
    },
    {
      key: "service",
      name: "服务：",
      type: FormItemTypes.Select,
      placeholder: "请选择要筛选的调用服务",
      options: serviceOptions,
    },
    {
      key: "route",
      name: "路由",
      type: FormItemTypes.Select,
      placeholder: "请选择要筛选的调用路由",
      options: routeOptions,
    },
    {
      key: "exception",
      name: "异常",
      type: FormItemTypes.Select,
      placeholder: "请选择是否筛选异常出错",
      options: [
        {
          label: "所有",
          value: "",
        },
        {
          label: "是",
          value: "true",
        },
      ],
    },
    {
      key: "limit",
      name: "查询数量：",
      type: FormItemTypes.InputNumber,
      placeholder: "请输入要查询的记录数量",
    },
    {
      key: "useGt",
      name: "请求耗时大于：",
      type: FormItemTypes.InputNumber,
      placeholder: "请输入要查询的耗时记录",
    },
    {
      key: "begin:end",
      name: "开始结束时间：",
      type: FormItemTypes.DateRange,
      span: 8,
      defaultValue: [getHoursAge(3).toISOString(), new Date().toISOString()],
    },
  ];
}

function getColumns(): TableColumn[] {
  return [
    {
      title: "服务名称",
      key: "service",
      width: 100,
    },
    {
      title: "请求路由",
      key: "route",
      width: 200,
      ellipsis: {
        tooltip: true,
      },
    },
    {
      title: "完整地址",
      key: "uri",
      width: 100,
      ellipsis: {
        tooltip: true,
      },
    },
    {
      title: "结果",
      key: "rslt",
      width: 80,
      render(row: Record<string, unknown>) {
        if (row.rslt === "0") {
          return "成功";
        }
        return "失败";
      },
    },
    {
      title: "状态码",
      key: "status",
      width: 70,
    },
    {
      title: "请求地址",
      key: "addr",
      width: 150,
    },
    {
      title: "耗时",
      key: "use",
      width: 100,
      render(row: Record<string, unknown>) {
        let use = 0;
        if (row.use) {
          use = row.use as number;
        }
        if (!use) {
          return "--";
        }
        const slots = {
          trigger: () => <span>{use.toLocaleString()}ms</span>,
        };
        const list: VNode[] = [];
        const append = (key: string, name: string) => {
          if (!row[key]) {
            return;
          }
          const use = row[key] as number;
          list.push(
            <li>
              {name}: {use.toLocaleString()}ms
            </li>
          );
        };
        append("dnsUse", "DNS");
        append("tcpUse", "TCP");
        append("tlsUse", "TLS");
        append("processingUse", "PROCESSING");
        append("transferUse", "TRANSFER");

        return (
          <NPopover v-slots={slots}>
            <ul class={useTimeListClass}>{list}</ul>
          </NPopover>
        );
      },
    },
    {
      title: "出错类型",
      key: "errCategory",
      width: 100,
    },
    {
      title: "异常",
      key: "exception",
      width: 80,
      render(row: Record<string, unknown>) {
        if (row.exception) {
          return "是";
        }
        return "否";
      },
    },
    {
      title: "出错信息",
      key: "error",
    },
    {
      title: "完整记录",
      key: "requestDetail",
      width: 90,
      render(row: Record<string, unknown>) {
        return <ExFluxDetail data={row} />;
      },
    },
    {
      title: "时间",
      key: "createdAt",
      width: 180,
      fixed: "right",
    },
  ];
}
export default defineComponent({
  name: "RequestStats",
  setup() {
    const message = useMessage();
    // const { requestServices, requestRoutes } = useFluxState();
    const requestsStore = useRequestsStore();
    const {
      requests,
      count,
      routes,
      fetchingRoutes,
      services,
      fetchingServices,
    } = storeToRefs(requestsStore);

    onMounted(async () => {
      try {
        await requestsStore.listRoute();
        await requestsStore.listService();
      } catch (err) {
        showError(message, err);
      }
    });

    onUnmounted(() => {
      requestsStore.$reset();
    });

    return {
      fetch: requestsStore.list,
      requests,
      count,
      routes,
      fetchingRoutes,
      services,
      fetchingServices,
    };
  },
  render() {
    const {
      requests,
      count,
      routes,
      fetchingRoutes,
      services,
      fetchingServices,
      fetch,
    } = this;
    if (fetchingServices || fetchingRoutes) {
      return <ExLoading />;
    }
    if (serviceOptions.length === 1 && services.length !== 0) {
      services.forEach((item) => {
        serviceOptions.push({
          label: item,
          value: item,
        });
      });
    }

    if (routeOptions.length === 1 && routes.length !== 0) {
      routes.forEach((item) => {
        routeOptions.push({
          label: item,
          value: item,
        });
      });
    }

    return (
      <ExTable
        disableAutoFetch={true}
        hidePagination={true}
        title={"HTTP请求统计"}
        data={{
          items: requests,
          count,
        }}
        filters={getFilters()}
        columns={getColumns()}
        fetch={fetch}
      />
    );
  },
});
