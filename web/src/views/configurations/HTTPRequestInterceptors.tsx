import { useMessage } from "naive-ui";
import { defineComponent, onMounted } from "vue";

import { showError } from "../../helpers/util";
import ExLoading from "../../components/ExLoading";
import { FormItemTypes, FormItem } from "../../components/ExForm";
import ExConfigEditorList from "../../components/ExConfigEditorList";
import { ConfigCategory } from "../../stores/configs";
import {
  getDefaultFormRules,
  newRequireRule,
} from "../../components/ExConfigEditor";
import { useCommonStore } from "../../stores/common";
import { storeToRefs } from "pinia";

export default defineComponent({
  name: "HTTPRequestInterceptorConfigs",
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
      httpInstances,
      fetchingHTTPInstances,
    };
  },
  render() {
    const { fetchingHTTPInstances, httpInstances } = this;
    if (fetchingHTTPInstances) {
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
        key: "data.service",
        type: FormItemTypes.Select,
        placeholder: "请选择拦截配置的实例",
        options: httpInstances.map((item) => {
          return {
            label: item.name,
            value: item.name,
          };
        }),
      },
      {
        name: "方法：",
        key: "data.method",
        type: FormItemTypes.Select,
        options: ["GET", "POST", "PUT", "DELETE"].map((item) => {
          return {
            label: item,
            value: item,
          };
        }),
        placeholder: "请选择请求方法",
      },
      {
        name: "路由：",
        key: "data.route",
        placeholder: "请输入要拦截的路由",
      },
      {
        name: "前置脚本：",
        key: "data.before",
        type: FormItemTypes.TextArea,
        span: 24,
        placeholder: "请输入请求处理前的相关处理脚本",
      },
      {
        name: "后置脚本：",
        key: "data.after",
        type: FormItemTypes.TextArea,
        span: 24,
        placeholder: "请输入请求处理后的相关处理脚本",
      },
    ];

    const rules = getDefaultFormRules({
      "data.name": newRequireRule("服务实例不能为空"),
      "data.method": newRequireRule("请求方法不能为空"),
      "data.router": newRequireRule("请求路由不能为空"),
    });

    return (
      <ExConfigEditorList
        listTitle="HTTP请求拦截配置"
        editorTitle="添加/更新HTTP请求拦截配置"
        editorDescription="设置各HTTP请求实例的拦截处理"
        category={ConfigCategory.HTTPRequestInterceptor}
        extraFormItems={extraFormItems}
        rules={rules}
      />
    );
  },
});
