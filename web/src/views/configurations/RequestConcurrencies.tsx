import { useMessage } from "naive-ui";
import { defineComponent, onMounted } from "vue";
import ExConfigEditorList from "../../components/ExConfigEditorList";
import { FormItemTypes, FormItem } from "../../components/ExForm";
import ExLoading from "../../components/ExLoading";
import { showError } from "../../helpers/util";
import { ConfigCategory } from "../../stores/configs";
import {
  getDefaultFormRules,
  newRequireRule,
} from "../../components/ExConfigEditor";
import { useCommonStore } from "../../stores/common";
import { storeToRefs } from "pinia";

export default defineComponent({
  name: "RequestConcurrencyConfigs",
  setup() {
    const message = useMessage();
    const commonStore = useCommonStore();
    const { httpInstances, fetchingHTTPInstances } = storeToRefs(commonStore);

    onMounted(async () => {
      try {
        await commonStore.listHTTPInstance();
      } catch (err) {
        showError(message, err);
      }
    });

    return {
      processing: fetchingHTTPInstances,
      httpInstances,
    };
  },

  render() {
    const { httpInstances, processing } = this;
    if (processing) {
      return <ExLoading />;
    }
    const extraFormItems: FormItem[] = [
      {
        type: FormItemTypes.Blank,
        name: "",
        key: "",
      },
      {
        name: "实例：",
        key: "data.name",
        type: FormItemTypes.Select,
        placeholder: "请选择限制并发数的实例",
        options: httpInstances.map((item) => {
          return {
            label: item.name,
            value: item.name,
          };
        }),
      },
      {
        name: "并发数：",
        key: "data.max",
        type: FormItemTypes.InputNumber,
        placeholder: "请输入并发限制",
      },
    ];
    const rules = getDefaultFormRules({
      "data.name": newRequireRule("服务实例不能为空"),
      "data.max": {
        required: true,
        message: "并发限制不能为空",
        trigger: "blur",
        validator(rule, value) {
          if (!value) {
            return false;
          }
          return true;
        },
      },
    });

    return (
      <ExConfigEditorList
        listTitle="HTTP请求实例并发配置"
        editorTitle="添加/更新HTTP请求实例并发限制"
        editorDescription="设置各HTTP请求实例的并发请求数"
        category={ConfigCategory.RequestConcurrency}
        extraFormItems={extraFormItems}
        rules={rules}
      />
    );
  },
});
