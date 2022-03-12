import { defineComponent, onMounted } from "vue";
import ExConfigEditorList from "../../components/ExConfigEditorList";
import { ConfigCategory } from "../../stores/configs";
import { FormItemTypes, FormItem } from "../../components/ExForm";
import { useMessage } from "naive-ui";
import { showError } from "../../helpers/util";
import ExLoading from "../../components/ExLoading";
import {
  getDefaultFormRules,
  newRequireRule,
} from "../../components/ExConfigEditor";
import { useCommonStore } from "../../stores/common";
import { storeToRefs } from "pinia";

export default defineComponent({
  name: "RouterConcurrencyConfigs",
  setup() {
    const commonStore = useCommonStore();
    const { routers, fetchingRouters } = storeToRefs(commonStore);
    const message = useMessage();
    onMounted(async () => {
      try {
        await commonStore.listRouter();
      } catch (err) {
        showError(message, err);
      }
    });
    return {
      processing: fetchingRouters,
      routers,
    };
  },
  render() {
    const { routers, processing } = this;
    if (processing) {
      return <ExLoading />;
    }
    const extraFormItems: FormItem[] = [
      {
        name: "路由：",
        key: "data.router",
        type: FormItemTypes.Select,
        placeholder: "请选择路由",
        options: routers.map((item) => {
          const value = `${item.method} ${item.route}`;
          return {
            label: value,
            value,
          };
        }),
      },
      {
        name: "最大并发数：",
        key: "data.max",
        type: FormItemTypes.InputNumber,
        placeholder: "最大并发限制",
      },
      {
        name: "频率限制阈值：",
        key: "data.rate",
        type: FormItemTypes.InputNumber,
        placeholder: "频率限制阈值",
      },
      {
        name: "频率时间区间：",
        key: "data.interval",
        type: FormItemTypes.InputDuration,
        placeholder: "频率时间区间",
      },
    ];
    const rules = getDefaultFormRules({
      "data.router": newRequireRule("路由不能为空"),
    });
    return (
      <ExConfigEditorList
        listTitle="路由并发限制"
        editorTitle="添加/更新路由并发限制"
        editorDescription="设置各路由并发限制"
        category={ConfigCategory.RouterConcurrency}
        extraFormItems={extraFormItems}
        rules={rules}
      />
    );
  },
});
