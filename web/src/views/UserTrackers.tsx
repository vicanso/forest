import { onMounted, onUnmounted, defineComponent } from "vue";
import { useMessage } from "naive-ui";
import ExTable from "../components/ExTable";
import ExLoading from "../components/ExLoading";
import useFluxState, {
  fluxListUserTracker,
  fluxListUserTrackerClear,
  fluxListUserTrackAction,
} from "../states/flux";
import { showError, today } from "../helpers/util";
import { FormItemTypes } from "../components/ExForm";

function getColumns() {
  return [
    {
      title: "账户",
      key: "account",
      width: 100,
    },
    {
      type: "expand",
      expandable: (data: Record<string, unknown>) =>
        !!data.form || !!data.query || !!data.params,
      renderExpand: (data: Record<string, unknown>) => {
        const arr = [];
        if (data.form) {
          arr.push(<pre>form: {data.form}</pre>);
        }
        if (data.query) {
          arr.push(<pre>query: {data.query}</pre>);
        }
        if (data.params) {
          arr.push(<pre>params: {data.params}</pre>);
        }
        return <div>{arr}</div>;
      },
    },
    {
      title: "类别",
      key: "action",
      width: 160,
    },
    {
      title: "结果",
      key: "resultDesc",
      width: 80,
    },
    {
      title: "TrackID",
      key: "tid",
      width: 220,
    },
    {
      title: "SessionID",
      key: "sid",
      width: 220,
    },
    {
      title: "IP",
      key: "ip",
      width: 140,
    },
    {
      title: "出错信息",
      key: "error",
      ellipsis: {
        tooltip: true,
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

const actionOptions = [
  {
    label: "所有",
    value: "",
  },
];
function getFilters() {
  return [
    {
      key: "account",
      name: "账户：",
      placeholder: "请输入要筛选的账号",
    },
    {
      key: "action",
      name: "类别：",
      placeholder: "请选择要筛选的类别",
      type: FormItemTypes.Select,
      options: actionOptions,
    },
    {
      key: "result",
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
      key: "limit",
      name: "查询数量",
      type: FormItemTypes.InputNumber,
      placeholder: "请输入要查询的记录数量",
    },
    {
      key: "begin:end",
      name: "开始结束时间：",
      type: FormItemTypes.DateRange,
      span: 12,
      defaultValue: [today().toISOString(), new Date().toISOString()],
    },
  ];
}

export default defineComponent({
  name: "UserTrackers",
  setup() {
    const message = useMessage();
    const { userTrackers, userTrackerActions } = useFluxState();
    // 加载用户行为类别
    onMounted(async () => {
      try {
        await fluxListUserTrackAction();
      } catch (err) {
        showError(message, err);
      }
    });
    // 清除数据
    onUnmounted(() => {
      fluxListUserTrackerClear();
    });

    return {
      userTrackers,
      userTrackerActions,
    };
  },
  render() {
    const { userTrackers, userTrackerActions } = this;
    if (userTrackerActions.processing) {
      return <ExLoading />;
    }
    // 添加类别选项
    if (actionOptions.length === 1 && userTrackerActions.items.length !== 0) {
      userTrackerActions.items.forEach((item) =>
        actionOptions.push({
          label: item,
          value: item,
        })
      );
    }
    return (
      <ExTable
        hidePagination={true}
        filters={getFilters()}
        title={"用户行为查询"}
        columns={getColumns()}
        data={userTrackers}
        fetch={fluxListUserTracker}
      />
    );
  },
});
