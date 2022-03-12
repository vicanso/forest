import { defineComponent, onUnmounted, PropType } from "vue";
import { TableColumn } from "naive-ui/lib/data-table/src/interface";
import ExTable, {
  newOPColumn,
  newJSONRenderExpand,
  newStatusValueColumn,
} from "../components/ExTable";
import { formatDate } from "../helpers/util";
import { useConfigsStore } from "../stores/configs";
import { storeToRefs } from "pinia";

function getColumns(): TableColumn[] {
  return [
    {
      title: "名称",
      key: "name",
    },
    newJSONRenderExpand("data"),
    {
      title: "分类",
      key: "category",
    },
    newStatusValueColumn({
      title: "状态",
      key: "status",
    }),
    {
      title: "创建者",
      key: "owner",
    },
    {
      title: "配置生效时间",
      key: "startedAt",
      render(row: Record<string, unknown>) {
        return formatDate(row.startedAt as string);
      },
    },
    {
      title: "配置失效时间",
      key: "endedAt",
      render(row: Record<string, unknown>) {
        return formatDate(row.endedAt as string);
      },
    },
    {
      title: "配置描述",
      key: "description",
      width: 100,
      ellipsis: {
        tooltip: true,
      },
    },
  ];
}

function noop(): void {
  // 无操作
}

export default defineComponent({
  name: "ConfigTable",
  props: {
    title: {
      type: String,
      default: "",
    },
    category: {
      type: String,
      default: () => "",
    },
    onUpdate: {
      type: Function as PropType<(id: number) => void>,
      default: noop,
    },
  },
  setup(props) {
    const configsStore = useConfigsStore();
    const { configs, count } = storeToRefs(configsStore);

    const fetchConfigs = () =>
      configsStore.list({
        category: props.category,
      });

    onUnmounted(() => {
      configsStore.$reset();
    });
    return {
      fetchConfigs,
      count,
      configs,
    };
  },
  render() {
    const { title, onUpdate } = this.$props;
    const { configs, fetchConfigs, $slots, count } = this;
    const columns = getColumns();
    if (onUpdate !== noop) {
      columns.push(
        newOPColumn((row) => {
          onUpdate(row.id as number);
        })
      );
    }
    return (
      <ExTable
        title={title}
        columns={columns}
        data={{
          items: configs,
          count,
        }}
        fetch={fetchConfigs}
      >
        {$slots}
      </ExTable>
    );
  },
});
