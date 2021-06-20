import { useMessage } from "naive-ui";
import { defineComponent, onMounted, onUnmounted } from "vue";
import ExLoading from "../components/ExLoading";
import { showError } from "../helpers/util";
import useFluxState, {
  fluxListRequest,
  fluxListRequestClear,
  fluxListRequestRoute,
  fluxListRequestService,
} from "../states/flux";
import ExTable from "../components/ExTable";
import { getHoursAge } from "../helpers/util";

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

function getFilters() {
  return [
    {
      key: "result",
      name: "结果：",
      type: "select",
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
      type: "select",
      placeholder: "请选择要筛选的调用服务",
      options: serviceOptions,
    },
    {
      key: "route",
      name: "路由",
      type: "select",
      placeholder: "请选择要筛选的调用路由",
      options: routeOptions,
    },
    {
      key: "exception",
      name: "异常",
      type: "select",
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
      type: "inputNumber",
      placeholder: "请输入要查询的记录数量",
    },
    {
      key: "begin:end",
      name: "开始结束时间：",
      type: "daterange",
      span: 12,
      defaultValue: [getHoursAge(3).toISOString(), new Date().toISOString()],
    },
  ];
}

function getColumns() {
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
      key: "result",
      width: 80,
      render(row: Record<string, unknown>) {
        if (row.result === "0") {
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
      render(row: Record<string, unknown>) {
        let use = 0;
        if (row.use) {
          use = row.use as number;
        }
        return `${use.toLocaleString()}ms`;
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
      title: "时间",
      key: "createdAt",
      width: 180,
      fixed: "right",
    },
  ];
}
export default defineComponent({
  name: "Requests",
  setup() {
    const message = useMessage();
    const { requestServices, requests, requestRoutes } = useFluxState();

    onMounted(async () => {
      try {
        await fluxListRequestService();
        await fluxListRequestRoute();
      } catch (err) {
        showError(message, err);
      }
    });

    onUnmounted(() => {
      fluxListRequestClear();
    });

    return {
      requests,
      requestServices,
      requestRoutes,
    };
  },
  render() {
    const { requests, requestServices, requestRoutes } = this;
    if (requestServices.processing || requestRoutes.processing) {
      return <ExLoading />;
    }
    if (serviceOptions.length === 1 && requestServices.items.length !== 0) {
      requestServices.items.forEach((item) => {
        serviceOptions.push({
          label: item,
          value: item,
        });
      });
    }

    if (routeOptions.length === 1 && requestRoutes.items.length !== 0) {
      requestRoutes.items.forEach((item) => {
        routeOptions.push({
          label: item,
          value: item,
        });
      });
    }

    return (
      <ExTable
        hidePagination={true}
        title={"HTTP请求统计"}
        data={requests}
        filters={getFilters()}
        columns={getColumns()}
        fetch={fluxListRequest}
      />
    );
  },
});
